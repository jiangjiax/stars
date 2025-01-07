package server

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jiangjiax/stars/internal/asset"
	"github.com/jiangjiax/stars/internal/config"
	"github.com/jiangjiax/stars/internal/post"
	"github.com/jiangjiax/stars/internal/rss"
	"github.com/jiangjiax/stars/internal/template"
)

type Server struct {
	config     *config.Config
	projectDir string
	port       int
	engine     *template.Engine
	posts      *post.Store
	assets     *asset.Pipeline
	router     chi.Router
}

func New(projectDir string, cfg *config.Config, port int) (*Server, error) {
	// 加载项目配置
	cfg, err := config.LoadConfig(filepath.Join(projectDir, "config.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 初始化资源管理器
	assets := asset.New(projectDir, cfg.Theme)

	// 构建资源
	if err := assets.BuildAssets(); err != nil {
		return nil, fmt.Errorf("failed to build assets: %w", err)
	}

	// 初始化模板引擎，设置为服务器模式
	engine, err := template.New(projectDir, cfg, false, assets)
	if err != nil {
		return nil, fmt.Errorf("failed to create template engine: %w", err)
	}

	// 创建文章存储
	posts := post.New()

	// 递归加载所有文章
	postsDir := filepath.Join(projectDir, "content/posts")
	err = filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理 .md 文件
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}

		// 解析文章
		parsePost, err := post.ParsePost(path)
		if err != nil {
			return fmt.Errorf("failed to parse post %s: %w", path, err)
		}

		// 如果没有验证信息或内容有变化，更新 contentHash
		if parsePost.Verification == nil || parsePost.ContentChanged() {
			if err := parsePost.UpdateContentHash(); err != nil {
				return fmt.Errorf("failed to update content hash: %w", err)
			}
		}

		// 如果没有设置 slug，使用相对路径作为 URL
		if parsePost.Slug == "" {
			relPath, err := filepath.Rel(postsDir, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}
			// 移除 .md 后缀
			relPath = strings.TrimSuffix(relPath, ".md")
			// 将路径分隔符转换为 URL 分隔符
			parsePost.Slug = strings.ReplaceAll(relPath, string(filepath.Separator), "/")
		}

		if err := posts.Add(parsePost); err != nil {
			return fmt.Errorf("failed to add post %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load posts: %w", err)
	}

	// 创建服务器实例
	srv := &Server{
		config:     cfg,
		projectDir: projectDir,
		port:       port,
		engine:     engine,
		posts:      posts,
		assets:     assets,
	}

	return srv, nil
}

func (s *Server) Start() error {
	// 创建路由器
	router := chi.NewRouter()
	s.router = router

	// 创建文件系统
	staticFS := http.Dir(filepath.Join(s.projectDir, "static"))
	themeFS := http.Dir(filepath.Join(s.projectDir, "themes", s.config.Theme, "static"))

	// 静态文件处理器
	fileServer := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 先尝试从项目静态目录提供文件
		if f, err := staticFS.Open(r.URL.Path); err == nil {
			f.Close()
			http.FileServer(staticFS).ServeHTTP(w, r)
			return
		}

		// 如果项目静态目录没有，则从主题静态目录提供
		if f, err := themeFS.Open(r.URL.Path); err == nil {
			f.Close()
			http.FileServer(themeFS).ServeHTTP(w, r)
			return
		}

		// 如果都没有找到，返回 404
		http.NotFound(w, r)
	})

	// 注册路由
	router.Handle("/static/*", http.StripPrefix("/static/", s.addCorrectMIMETypes(fileServer)))

	// RSS feed 路由
	router.Get("/feed.xml", func(w http.ResponseWriter, r *http.Request) {
		// 获取所有文章并按日期排序
		posts := s.posts.GetAll()
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Date.After(posts[j].Date)
		})

		// 生成 RSS feed
		feed, err := rss.GenerateFeed(posts, s.config.Title, s.config.BaseURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 输出 RSS
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "  ")
		encoder.Encode(feed)
	})

	// 其他路由
	router.Get("/*", s.handleContent)

	// 启动服务器
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server at http://localhost%s", addr)
	return http.ListenAndServe(addr, router)
}

func (s *Server) handleContent(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling request for path: %s", r.URL.Path)

	// 处理首页
	if r.URL.Path == "/" {
		log.Printf("Rendering home page")
		data := map[string]interface{}{
			"Title": s.config.Title,
			"Site":  s.config,
		}
		html, err := s.engine.RenderHome(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
		return
	}

	// 处理标签页面
	if strings.HasPrefix(r.URL.Path, "/tags/") {
		s.handleTaxonomy(w, r, "tags")
		return
	}

	// 处理系列页面
	if strings.HasPrefix(r.URL.Path, "/series/") {
		s.handleTaxonomy(w, r, "series")
		return
	}

	// 处理文章列表页和分页
	if r.URL.Path == "/posts" || strings.HasPrefix(r.URL.Path, "/posts/page/") {
		// 获取页码
		page := 1
		if strings.HasPrefix(r.URL.Path, "/posts/page/") {
			pageStr := strings.TrimPrefix(r.URL.Path, "/posts/page/")
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		// 更新请求参数
		q := r.URL.Query()
		q.Set("page", strconv.Itoa(page))
		r.URL.RawQuery = q.Encode()

		s.handleList(w, r)
		return
	}

	// 处理单篇文章页面
	if strings.HasPrefix(r.URL.Path, "/posts/") && !strings.Contains(r.URL.Path, "/page/") {
		s.handlePost(w, r)
		return
	}

	// 处理标签云页面
	if r.URL.Path == "/tags" {
		data := map[string]interface{}{
			"Title":     "标签云 - " + s.config.Title,
			"AllTags":   s.posts.GetAllTags(),
			"TagsStats": s.posts.GetTagsStats(),
			"Site":      s.config,
		}
		html, err := s.engine.RenderTags(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, html)
		return
	}

	http.NotFound(w, r)
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	// 获取文章 slug
	slug := strings.TrimPrefix(r.URL.Path, "/posts/")

	// 查找对应的文章
	p, err := s.posts.GetBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 获取所有文章用系列导航
	allPosts := s.posts.GetAll()

	// 创建模板数据
	data := map[string]interface{}{
		"Title": p.Title + " - " + s.config.Title,
		"Post":  p,
		"Posts": allPosts, // 传递所有文章数据
		"Site":  s.config,
	}

	// 渲染模板
	html, err := s.engine.RenderSingle(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func (s *Server) loadPosts() error {
	contentDir := filepath.Join(s.projectDir, "content")
	contentFS := os.DirFS(contentDir)
	if err := s.posts.LoadFromFS(contentFS, "posts"); err != nil {
		return fmt.Errorf("failed to load posts: %w", err)
	}
	return nil
}

func (s *Server) handleGetPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := s.posts.GetBySlug(slug)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Title": post.Title + " - " + s.config.Title,
		"Post":  post,
	}

	html, err := s.engine.RenderSingle(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// 添加新的辅助方法来设置正确的 MIME 类型
func (s *Server) addCorrectMIMETypes(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 根据文件扩展名设置正确的 MIME 类型
		switch {
		case strings.HasSuffix(r.URL.Path, ".css"):
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case strings.HasSuffix(r.URL.Path, ".js"):
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case strings.HasSuffix(r.URL.Path, ".svg"):
			w.Header().Set("Content-Type", "image/svg+xml")
		case strings.HasSuffix(r.URL.Path, ".png"):
			w.Header().Set("Content-Type", "image/png")
		case strings.HasSuffix(r.URL.Path, ".jpg"), strings.HasSuffix(r.URL.Path, ".jpeg"):
			w.Header().Set("Content-Type", "image/jpeg")
		}
		next.ServeHTTP(w, r)
	})
}

// handleList 处理文章列表页面
func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	// 获取分页参数
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 6

	// 获取并解码筛选参数
	tag, _ := url.QueryUnescape(r.URL.Query().Get("tag"))
	series, _ := url.QueryUnescape(r.URL.Query().Get("series"))

	// 获取文章列表
	var posts []*post.Post
	var total int

	// 根据筛选条件获取文章
	if tag != "" && series != "" {
		posts, total = s.posts.ListByTagAndSeries(tag, series, page, pageSize)
	} else if tag != "" {
		posts, total = s.posts.ListByTag(tag, page, pageSize)
	} else if series != "" {
		posts, total = s.posts.ListBySeries(series, page, pageSize)
	} else {
		posts, total = s.posts.ListPaged(page, pageSize)
	}

	// 获取预计算的统计数据
	seriesStats := s.posts.GetSeriesStats()
	tagsStats := s.posts.GetTagsStats()
	allTags := s.posts.GetAllTags()

	// 准备模板数据
	data := map[string]interface{}{
		"Title":       "文章列表 - " + s.config.Title,
		"Posts":       posts,
		"CurrentPage": page,
		"PageSize":    pageSize,
		"TotalPosts":  total,
		"Tag":         tag,
		"Series":      series,
		"Site":        s.config,
		"BaseURL":     "/posts",
		"SeriesStats": seriesStats, // 添加系列统计
		"TagsStats":   tagsStats,   // 添加标签统计
		"AllTags":     allTags,     // 添加所有标签
		"Pagination":  template.NewPagination(page, pageSize, total, "/posts"),
	}

	// 在传入数据前对 Series 进行排序
	sort.Slice(s.config.Series, func(i, j int) bool {
		return s.config.Series[i].Order < s.config.Series[j].Order
	})

	html, err := s.engine.RenderList(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// handleTaxonomy 处理分类页面（标签和系列）
func (s *Server) handleTaxonomy(w http.ResponseWriter, r *http.Request, taxonomy string) {
	// 解析路径获取分类项和页码
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	// URL 解码 term
	term, err := url.QueryUnescape(parts[1])
	if err != nil {
		http.Error(w, "Invalid term", http.StatusBadRequest)
		return
	}

	page := 1

	// 检查是否有分页
	if len(parts) > 3 && parts[2] == "page" {
		if p, err := strconv.Atoi(parts[3]); err == nil && p > 0 {
			page = p
		}
	}

	// 获取文章列表
	pageSize := 6
	var posts []*post.Post
	var totalPosts int

	if taxonomy == "tags" {
		posts = s.posts.GetPostsByTag(term)
	} else {
		posts = s.posts.GetPostsBySeries(term)
	}
	totalPosts = len(posts)

	// 按日期排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// 计算当前页的文章
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > totalPosts {
		end = totalPosts
	}

	// 创建分页数据
	baseURL := fmt.Sprintf("/%s/%s", taxonomy, term)
	pagination := template.NewPagination(page, pageSize, totalPosts, baseURL)

	data := map[string]interface{}{
		"Title":       fmt.Sprintf("%s: %s - 第%d页", taxonomy, term, page),
		"Posts":       posts[start:end],
		"Site":        s.config,
		"Taxonomy":    taxonomy,
		"Term":        term,
		"SeriesStats": s.posts.GetSeriesStats(),
		"TagsStats":   s.posts.GetTagsStats(),
		"AllTags":     s.posts.GetAllTags(),
		"Pagination":  pagination,
		"TotalPosts":  totalPosts,
	}

	html, err := s.engine.RenderList(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

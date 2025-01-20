package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jiangjiax/stars/internal/asset"
	"github.com/jiangjiax/stars/internal/post"
	"github.com/jiangjiax/stars/internal/rss"
	"github.com/jiangjiax/stars/internal/sitemap"
	"github.com/jiangjiax/stars/internal/template"
)

// Build generates the static website from the project
func (p *Project) Build() error {
	// 初始化资源管理器
	assets := asset.New(p.Path, p.Site.Theme)

	// 构建资源
	if err := assets.BuildAssets(); err != nil {
		return fmt.Errorf("failed to build assets: %w", err)
	}

	engine, err := template.New(p.Path, p.Site, true, assets)
	if err != nil {
		return fmt.Errorf("failed to initialize template engine: %w", err)
	}

	builder := &Builder{
		project:   p,
		publicDir: filepath.Join(p.Path, "public"),
		engine:    engine,
	}
	return builder.Build()
}

// Builder handles the static site generation process
type Builder struct {
	project   *Project
	publicDir string
	engine    *template.Engine
}

// Build executes the full build process
func (b *Builder) Build() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"Initialize templates", b.project.initBuildTemplates},
		{"Clean public directory", b.cleanPublicDir},
		{"Parse posts", b.parsePosts},
		{"Generate posts", b.generatePosts},
		{"Generate index page", b.generateIndex},
		{"Generate list page", b.generateList},
		{"Copy static files", b.copyStaticFiles},
		{"Generate paginated lists", b.generatePaginatedLists},
		{"Generate taxonomy pages", b.generateTaxonomyPages},
		{"Generate tags page", b.generateTagsPage},
	}

	for _, step := range steps {
		fmt.Printf("Building: %s\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	// 生成 RSS feed
	feedPath := filepath.Join(b.publicDir, "feed.xml")
	if err := rss.GenerateAndSaveFeed(
		b.project.Posts,
		b.project.Site.Title,
		b.project.Site.BaseURL,
		feedPath,
	); err != nil {
		return fmt.Errorf("failed to generate RSS feed: %w", err)
	}

	// 生成 sitemap
	sitemapGen := sitemap.New(b.project.Site, b.project.Posts)
	if err := sitemapGen.Generate(b.publicDir); err != nil {
		return fmt.Errorf("failed to generate sitemap: %w", err)
	}

	return nil
}

// cleanPublicDir cleans and recreates the public directory
func (b *Builder) cleanPublicDir() error {
	if err := os.RemoveAll(b.publicDir); err != nil {
		return fmt.Errorf("failed to clean public directory: %w", err)
	}
	return os.MkdirAll(b.publicDir, 0755)
}

// parsePosts reads and parses all markdown posts
func (b *Builder) parsePosts() error {
	postsDir := filepath.Join(b.project.Path, "content", "posts")
	err := filepath.Walk(postsDir, func(path string, info os.FileInfo, err error) error {
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

		// 添加与 server 命令相同的 contentHash 检查
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
			relPath = strings.TrimSuffix(relPath, ".md")
			parsePost.Slug = strings.ReplaceAll(relPath, string(filepath.Separator), "/")
		}

		b.project.Posts = append(b.project.Posts, parsePost)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to parse posts: %w", err)
	}

	// Sort posts by series and order
	sort.Slice(b.project.Posts, func(i, j int) bool {
		if b.project.Posts[i].Series == b.project.Posts[j].Series {
			return b.project.Posts[i].SeriesOrder < b.project.Posts[j].SeriesOrder
		}
		return b.project.Posts[i].Series < b.project.Posts[j].Series
	})

	return nil
}

// generatePosts generates HTML pages for all posts
func (b *Builder) generatePosts() error {
	for _, post := range b.project.Posts {
		if err := b.generatePost(post); err != nil {
			return fmt.Errorf("failed to generate post %s: %w", post.Slug, err)
		}
	}
	return nil
}

// generatePost generates a single post's HTML page
func (b *Builder) generatePost(post *post.Post) error {
	// Create post directory
	postDir := filepath.Join(b.publicDir, "posts", post.Slug)
	if err := os.MkdirAll(postDir, 0755); err != nil {
		return fmt.Errorf("failed to create post directory: %w", err)
	}

	// Set current post for template rendering
	b.project.Post = post

	// Prepare template data
	templateData := map[string]interface{}{
		"Title": post.Title + " - " + b.project.Site.Title,
		"Post":  post,
		"Posts": b.project.Posts,
		"Site":  b.project.Site,
	}

	// 使用模板引擎渲染
	html, err := b.engine.RenderSingle(templateData)
	if err != nil {
		return fmt.Errorf("failed to render post: %w", err)
	}

	// Write post HTML file
	indexPath := filepath.Join(postDir, "index.html")
	if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write post file: %w", err)
	}

	return nil
}

// generateIndex generates the site's index page
func (b *Builder) generateIndex() error {
	data := map[string]interface{}{
		"Title": b.project.Site.Title,
		"Posts": b.project.Posts,
		"Site":  b.project.Site,
	}

	// 用模板引擎渲染
	html, err := b.engine.RenderHome(data)
	if err != nil {
		return fmt.Errorf("failed to render index page: %w", err)
	}

	indexPath := filepath.Join(b.publicDir, "index.html")
	if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write index.html: %w", err)
	}

	return nil
}

// copyStaticFiles copies theme static files to public directory
func (b *Builder) copyStaticFiles() error {
	// 复制主题静态文件
	themeStaticDir := filepath.Join(b.project.Path, "themes", b.project.Site.Theme, "static")
	publicStaticDir := filepath.Join(b.publicDir, "static")

	// 复制主题静态文件
	err := filepath.Walk(themeStaticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(themeStaticDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		destPath := filepath.Join(publicStaticDir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy theme static files: %w", err)
	}

	// 复制项目静态文件（如果存在）
	projectStaticDir := filepath.Join(b.project.Path, "static")
	if _, err := os.Stat(projectStaticDir); err == nil {
		err := filepath.Walk(projectStaticDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			relPath, err := filepath.Rel(projectStaticDir, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			destPath := filepath.Join(publicStaticDir, relPath)
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			if err := os.WriteFile(destPath, content, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}

			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to copy project static files: %w", err)
		}
	}

	return nil
}

// generateList generates the list page for posts
func (b *Builder) generateList() error {
	// 创建 post.Store 实例来获取统计信息
	store := post.New()
	for _, p := range b.project.Posts {
		if err := store.Add(p); err != nil {
			return fmt.Errorf("failed to add post to store: %w", err)
		}
	}

	// 对系列按照 Order 排序
	sort.Slice(b.project.Site.Series, func(i, j int) bool {
		return b.project.Site.Series[i].Order < b.project.Site.Series[j].Order
	})

	data := map[string]interface{}{
		"Title":       "Posts - " + b.project.Site.Title,
		"Posts":       b.project.Posts,
		"Site":        b.project.Site,
		"BuildMode":   true,
		"SeriesStats": store.GetSeriesStats(),
		"TagsStats":   store.GetTagsStats(),
		"AllTags":     store.GetAllTags(),
	}

	html, err := b.engine.RenderList(data)
	if err != nil {
		return fmt.Errorf("failed to render list page: %w", err)
	}

	// 创建 posts 目录
	postsDir := filepath.Join(b.publicDir, "posts")
	if err := os.MkdirAll(postsDir, 0755); err != nil {
		return fmt.Errorf("failed to create posts directory: %w", err)
	}

	// 写入 index.html
	indexPath := filepath.Join(postsDir, "index.html")
	if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write posts/index.html: %w", err)
	}

	return nil
}

// generatePaginatedLists 生成所有分页列表
func (b *Builder) generatePaginatedLists() error {
	pageSize := 6
	posts := b.project.Posts

	// 按日期排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	totalPosts := len(posts)
	baseURL := "/posts"

	// 创建 post.Store 实例来获取统计信息
	store := post.New()
	for _, p := range posts {
		if err := store.Add(p); err != nil {
			return fmt.Errorf("failed to add post to store: %w", err)
		}
	}

	// 计算总页数
	totalPages := (totalPosts + pageSize - 1) / pageSize

	for page := 1; page <= totalPages; page++ {
		start := (page - 1) * pageSize
		end := start + pageSize
		if end > totalPosts {
			end = totalPosts
		}

		// 创建分页数据
		pagination := template.NewPagination(page, pageSize, totalPosts, baseURL)

		data := map[string]interface{}{
			"Title":       fmt.Sprintf("文章列表 - 第%d页 - %s", page, b.project.Site.Title),
			"Posts":       posts[start:end],
			"Site":        b.project.Site,
			"BuildMode":   true,
			"SeriesStats": store.GetSeriesStats(),
			"TagsStats":   store.GetTagsStats(),
			"AllTags":     store.GetAllTags(),
			"Pagination":  pagination,
		}

		// 生成页面
		html, err := b.engine.RenderList(data)
		if err != nil {
			return fmt.Errorf("failed to render page %d: %w", page, err)
		}

		// 创建目录并写入文件
		var pageDir string
		if page == 1 {
			pageDir = filepath.Join(b.publicDir, "posts")
		} else {
			pageDir = filepath.Join(b.publicDir, "posts", "page", fmt.Sprintf("%d", page))
		}

		if err := os.MkdirAll(pageDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		indexPath := filepath.Join(pageDir, "index.html")
		if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
			return fmt.Errorf("failed to write page %d: %w", page, err)
		}
	}

	return nil
}

// generateTaxonomyPages 生成分类页面
func (b *Builder) generateTaxonomyPages() error {
	// 创建 post.Store 实例来获取标签和系列统计
	store := post.New()
	for _, p := range b.project.Posts {
		if err := store.Add(p); err != nil {
			return fmt.Errorf("failed to add post to store: %w", err)
		}
	}

	// 处理标签页面
	tagPosts := make(map[string][]*post.Post)
	for _, p := range b.project.Posts {
		for _, tag := range p.Tags {
			tagPosts[tag] = append(tagPosts[tag], p)
		}
	}

	// 处理系列页面
	seriesPosts := make(map[string][]*post.Post)
	for _, p := range b.project.Posts {
		if p.Series != "" {
			seriesPosts[p.Series] = append(seriesPosts[p.Series], p)
		}
	}

	// 为配置文件中定义的所有系列创建页面，即使没有文章
	for _, series := range b.project.Site.Series {
		if _, exists := seriesPosts[series.Name]; !exists {
			// 如果这个系列还没有文章，也要生成页面
			if err := b.generatePaginatedTaxonomyPages("series", series.Name, []*post.Post{}); err != nil {
				return fmt.Errorf("failed to generate empty series page for %s: %w", series.Name, err)
			}
		}
	}

	// 生成标签页面
	for tag, posts := range tagPosts {
		if err := b.generatePaginatedTaxonomyPages("tags", tag, posts); err != nil {
			return fmt.Errorf("failed to generate tag pages: %w", err)
		}
	}

	// 生成系列页面
	for series, posts := range seriesPosts {
		if err := b.generatePaginatedTaxonomyPages("series", series, posts); err != nil {
			return fmt.Errorf("failed to generate series pages: %w", err)
		}
	}

	return nil
}

func (b *Builder) generatePaginatedTaxonomyPages(taxonomy, term string, posts []*post.Post) error {
	pageSize := 6
	totalPosts := len(posts)
	totalPages := (totalPosts + pageSize - 1) / pageSize
	// 显示空状态页面
	// emptyState := false
	if totalPosts == 0 {
		// emptyState = true
		totalPages = 1
	}

	// 按日期排序
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// 创建 post.Store 实例来获取统计信息
	store := post.New()
	for _, p := range b.project.Posts {
		if err := store.Add(p); err != nil {
			return fmt.Errorf("failed to add post to store: %w", err)
		}
	}

	// 生成分页
	for page := 1; page <= totalPages; page++ {
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
			"Site":        b.project.Site,
			"Taxonomy":    taxonomy,
			"Term":        term,
			"SeriesStats": store.GetSeriesStats(),
			"TagsStats":   store.GetTagsStats(),
			"AllTags":     store.GetAllTags(),
			"Pagination":  pagination,
			"TotalPosts":  totalPosts,
			"BuildMode":   true,
		}

		// 生成页面
		html, err := b.engine.RenderList(data)
		if err != nil {
			return fmt.Errorf("failed to render page %d: %w", page, err)
		}

		// 处理文件夹名称：将空格替换为连字符，移除特殊字符
		dirName := strings.ReplaceAll(term, " ", "-")

		// 创建目录并写入文件
		var pageDir string
		if page == 1 {
			pageDir = filepath.Join(b.publicDir, taxonomy, dirName)
		} else {
			pageDir = filepath.Join(b.publicDir, taxonomy, dirName, "page", fmt.Sprintf("%d", page))
		}

		if err := os.MkdirAll(pageDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		indexPath := filepath.Join(pageDir, "index.html")
		if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
			return fmt.Errorf("failed to write page %d: %w", page, err)
		}
	}

	return nil
}

// generateTagsPage 生成标签云页面
func (b *Builder) generateTagsPage() error {
	// 创建 post.Store 实例来获取标签统计
	store := post.New()
	for _, p := range b.project.Posts {
		if err := store.Add(p); err != nil {
			return fmt.Errorf("failed to add post to store: %w", err)
		}
	}

	// 获取所有标签和统计信息
	allTags := store.GetAllTags()
	tagsStats := store.GetTagsStats()

	// 生成标签云页面
	data := map[string]interface{}{
		"Title":     "标签云 - " + b.project.Site.Title,
		"AllTags":   allTags,
		"TagsStats": tagsStats,
		"Site":      b.project.Site,
		"BuildMode": true,
	}

	html, err := b.engine.RenderTags(data)
	if err != nil {
		return fmt.Errorf("failed to render tags page: %w", err)
	}

	// 创建标签页面目录
	tagsDir := filepath.Join(b.publicDir, "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return fmt.Errorf("failed to create tags directory: %w", err)
	}

	// 写入 index.html
	indexPath := filepath.Join(tagsDir, "index.html")
	if err := os.WriteFile(indexPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write tags index: %w", err)
	}

	return nil
}

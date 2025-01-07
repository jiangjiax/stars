package generator

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	stdtmpl "text/template"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/jiangjiax/stars/internal/config"
	"github.com/jiangjiax/stars/internal/post"
	"github.com/jiangjiax/stars/internal/template"
	"github.com/jiangjiax/stars/internal/template/funcs"
)

//go:embed templates/config.yaml templates/example-posts/*
var templates embed.FS

// Project represents a new Stars blog project
type Project struct {
	Path        string
	Name        string
	Description string
	Date        string
	Site        *config.Config
	Posts       []*post.Post
	Post        *post.Post        // 当前文章（用于单文章页面）
	template    *stdtmpl.Template // 使用标准库的模板类型
}

// Now 返回当前时间
func (p *Project) Now() time.Time {
	return time.Now()
}

// New creates a new Stars blog project
func New(path string, useTemplate bool) (*Project, error) {
	name := filepath.Base(path)
	p := &Project{
		Path:        path,
		Name:        name,
		Description: "A Stars Web3 blog",
		Date:        time.Now().Format("2006-01-02"),
	}

	var cfg *config.Config
	var err error

	if useTemplate {
		// 用于 new 命令：使用模板配置
		cfg, err = config.LoadFromTemplate(templates, "templates/config.yaml", p)
	} else {
		// 用于 build 命令：使用项目配置
		cfg, err = config.LoadConfig(filepath.Join(path, "config.yaml"))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	p.Site = cfg

	// 初始化模板引擎
	tmpl := stdtmpl.New("").Funcs(funcs.DefaultFuncs)
	p.template = tmpl

	return p, nil
}

// Generate creates the project directory structure and files
func (p *Project) Generate() error {
	// 1. 创建基本目录结构
	dirs := []string{
		"content/posts",
		"public",
		"themes",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(p.Path, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 使用 defer 处理错误清理
	var success bool
	defer func() {
		if !success {
			// 如果生成失败，清理创建的目录
			if err := os.RemoveAll(p.Path); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to clean up directory %s: %v\n", p.Path, err)
			}
		}
	}()

	// 2. 安装默认主题
	if err := p.installDefaultTheme(); err != nil {
		return fmt.Errorf("failed to install default theme: %w", err)
	}

	// 3. 初始化模板
	if err := p.initTemplates(); err != nil {
		return fmt.Errorf("failed to init templates: %w", err)
	}

	// 4. 在主题目录下安装依赖
	if err := p.installDependencies(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	// 5. 复制示例文章
	if err := p.copyExamplePosts(); err != nil {
		return fmt.Errorf("failed to copy example posts: %w", err)
	}

	// 6. 生成配置文件
	configPath := filepath.Join(p.Path, "config.yaml")
	if err := p.generateFromTemplate(configPath, "templates/config.yaml", p); err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// 标记生成成功
	success = true
	return nil
}

// 复制示例文章
func (p *Project) copyExamplePosts() error {
	entries, err := templates.ReadDir("templates/example-posts")
	if err != nil {
		return fmt.Errorf("failed to read example posts directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		content, err := templates.ReadFile(filepath.Join("templates/example-posts", entry.Name()))
		if err != nil {
			return fmt.Errorf("failed to read example post %s: %w", entry.Name(), err)
		}

		destPath := filepath.Join(p.Path, "content/posts", entry.Name())
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write example post %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// initBuildTemplates 专门用于构建时的模板初始化
func (p *Project) initBuildTemplates() error {
	// 创建基础模板实例
	tmpl := stdtmpl.New("").Funcs(funcs.DefaultFuncs)
	p.template = tmpl

	// 加载主题目录下的所有模板
	themePath := filepath.Join(p.Path, "themes", p.Site.Theme, "layouts")
	err := filepath.Walk(themePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 HTML 文件
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		// 读取模板文件
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", path, err)
		}

		// 获取相对路径作为模板名称
		relPath, err := filepath.Rel(themePath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// 使用不带目录的文件名作为模板名称
		templateName := filepath.Base(relPath)

		// 解析模板
		if _, err := tmpl.New(templateName).Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		fmt.Printf("Loaded template: %s as %s\n", relPath, templateName)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// 打印所有已加载的模板
	fmt.Printf("Available templates: %v\n", tmpl.DefinedTemplates())

	return nil
}

// 初始化模板
func (p *Project) initTemplates() error {
	// 创建新的模板实例
	tmpl := stdtmpl.New("").Funcs(funcs.DefaultFuncs)

	// 解析所有模板文件
	themePath := filepath.Join(p.Path, "themes", p.Site.Theme, "layouts")
	err := filepath.Walk(themePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 HTML 文件
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			// 读取模板文件
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read template file %s: %w", path, err)
			}

			// 获取模板名称
			relPath, err := filepath.Rel(themePath, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}

			// 使用文件名作为模板名称
			templateName := filepath.Base(relPath)

			// 解析模板
			if _, err := tmpl.New(templateName).Parse(string(content)); err != nil {
				return fmt.Errorf("failed to parse template %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 保存模板实例
	p.template = tmpl
	return nil
}

// 生成单篇文章页面
func (p *Project) generatePost(currentPost *post.Post) error {
	// 创建文章目录
	postDir := filepath.Join(p.Path, "public", "posts", currentPost.Slug)
	if err := os.MkdirAll(postDir, 0755); err != nil {
		return fmt.Errorf("failed to create post directory: %w", err)
	}

	// 设置当前文章，用于模板渲染
	p.Post = currentPost

	// 确保 Posts 不为空
	if p.Posts == nil {
		p.Posts = make([]*post.Post, 0)
	}

	// 创建模板数据
	templateData := map[string]interface{}{
		"Title": currentPost.Title + " - " + p.Site.Title,
		"Post":  currentPost,
		"Posts": p.Posts,
		"Site":  p.Site,
	}

	// 渲染文章内容
	var buf bytes.Buffer
	if err := p.template.ExecuteTemplate(&buf, "single.html", templateData); err != nil {
		return fmt.Errorf("failed to render post %s: %w", currentPost.Slug, err)
	}

	// 写入该文章的 index.html，此处不是在复制模板文件，而是在生成每篇文章的 HTML 页面
	indexPath := filepath.Join(postDir, "index.html")
	if err := os.WriteFile(indexPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write post index.html: %w", err)
	}

	return nil
}

// 生成所有文章页面
func (p *Project) generatePosts() error {
	// 确保 Posts 不为空
	if p.Posts == nil {
		p.Posts = make([]*post.Post, 0)
	}

	if len(p.Posts) == 0 {
		return fmt.Errorf("no posts found")
	}

	// 先按系列和序号排序所有文章
	sort.Slice(p.Posts, func(i, j int) bool {
		if p.Posts[i].Series == p.Posts[j].Series {
			return p.Posts[i].SeriesOrder < p.Posts[j].SeriesOrder
		}
		return p.Posts[i].Series < p.Posts[j].Series
	})

	// 生成每篇文章
	for _, onePost := range p.Posts {
		if err := p.generatePost(onePost); err != nil {
			return fmt.Errorf("failed to generate post %s: %w", onePost.Slug, err)
		}
	}
	return nil
}

// 复制主题文件
func (p *Project) copyThemeFiles() error {
	return fs.WalkDir(templates, "templates/theme", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel("templates/theme", path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// 目标路径
		destPath := filepath.Join(p.Path, "themes/default", relPath)

		// 如果是目录，创建它
		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// 复制文件
		content, err := templates.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		if err := os.WriteFile(destPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", destPath, err)
		}

		return nil
	})
}

// generateFromTemplate creates a file from a template
func (p *Project) generateFromTemplate(destPath, templatePath string, data interface{}) error {
	// 创建模板引擎并加函数映射
	tmpl := stdtmpl.New(filepath.Base(templatePath)).Funcs(funcs.DefaultFuncs)

	// 解析模板
	content, err := templates.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templatePath, err)
	}

	if _, err := tmpl.Parse(string(content)); err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// 确保目标目录存在
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Create destination file
	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", destPath, err)
	}
	defer dest.Close()

	// 执行模板
	if err := tmpl.Execute(dest, data); err != nil {
		return fmt.Errorf("generateFromTemplate failed to execute template %s: %w", templatePath, err)
	}

	return nil
}

// installDependencies 安装项目依赖
func (p *Project) installDependencies() error {
	themePath := filepath.Join(p.Path, "themes", "default")
	pkgPath := filepath.Join(themePath, "package.json")

	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return fmt.Errorf("package.json not found")
	}

	// 安装依赖
	cmd := exec.Command("npm", "install")
	cmd.Dir = themePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *Project) generateTaxonomyPage(taxonomy, term string, posts []*post.Post) error {
	pageSize := 6
	totalPosts := len(posts)
	baseURL := fmt.Sprintf("/%s/%s", taxonomy, term)

	// 使用我们自己的 template 包创建分页数据
	pagination := &template.Pagination{
		CurrentPage: 1,
		PageSize:    pageSize,
		TotalPosts:  totalPosts,
		BaseURL:     baseURL,
		HasPrev:     false,
		HasNext:     totalPosts > pageSize,
		NextPage:    2,
	}

	// 准备模板数据
	data := map[string]interface{}{
		"Title":      fmt.Sprintf("%s: %s - %s", taxonomy, term, p.Site.Title),
		"Posts":      posts[:min(pageSize, totalPosts)],
		"Site":       p.Site,
		"Pagination": pagination,
		"BuildMode":  true,
	}

	// 渲染模板
	var buf bytes.Buffer
	if err := p.template.ExecuteTemplate(&buf, "list.html", data); err != nil {
		return fmt.Errorf("failed to render taxonomy page: %w", err)
	}

	// 创建目录
	pageDir := filepath.Join(p.Path, "public", taxonomy, term)
	if err := os.MkdirAll(pageDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 写入 index.html
	indexPath := filepath.Join(pageDir, "index.html")
	if err := os.WriteFile(indexPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// installDefaultTheme 安装默认主题
func (p *Project) installDefaultTheme() error {
	themesDir := filepath.Join(p.Path, "themes")
	defaultTheme := "default"

	// 使用 go-git 从 GitHub 克隆主题
	_, err := git.PlainClone(
		filepath.Join(themesDir, defaultTheme),
		false,
		&git.CloneOptions{
			URL:      "https://github.com/jiangjiax/stars-theme-default",
			Progress: os.Stdout,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to clone default theme: %w", err)
	}

	return nil
}

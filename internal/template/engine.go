package template

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/jiangjiax/stars/internal/asset"
	"github.com/jiangjiax/stars/internal/config"
	"github.com/jiangjiax/stars/internal/template/funcs"
)

type Engine struct {
	config     *config.Config
	layoutDir  string
	projectDir string
	buildMode  bool
	assets     *asset.Pipeline
}

// New 创建新的模板引擎
func New(projectDir string, cfg *config.Config, buildMode bool, assets *asset.Pipeline) (*Engine, error) {
	// 设置布局目录
	layoutDir := filepath.Join(projectDir, "themes", cfg.Theme, "layouts")
	if _, err := os.Stat(layoutDir); err != nil {
		return nil, fmt.Errorf("theme layouts directory not found: %w", err)
	}

	engine := &Engine{
		config:     cfg,
		layoutDir:  layoutDir,
		projectDir: projectDir,
		buildMode:  buildMode,
		assets:     assets,
	}

	// 设置资源管理器到模板函数
	funcs.SetAssetPipeline(engine.assets)

	return engine, nil
}

// render 执行模板渲染 - 使用模板方法模式重构
func (e *Engine) render(kind, section string, data interface{}) (string, error) {
	var buf strings.Builder

	// 1. 创建基础模板
	baseTemplate, err := e.createBaseTemplate()
	if err != nil {
		return "", err
	}

	// 2. 加载基础布局
	if err := e.loadBaseLayout(baseTemplate); err != nil {
		return "", err
	}

	// 3. 加载组件
	if err := e.loadComponents(baseTemplate); err != nil {
		return "", err
	}

	// 4. 加载页面模板
	if err := e.loadPageTemplate(baseTemplate, kind, section); err != nil {
		return "", err
	}

	// 5. 处理数据
	processedData, err := e.processTemplateData(data, kind, section)
	if err != nil {
		return "", err
	}

	// 6. 执行模板
	if err := baseTemplate.ExecuteTemplate(&buf, "_default/baseof.html", processedData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// 创建基础模板
func (e *Engine) createBaseTemplate() (*template.Template, error) {
	return template.New("").Funcs(funcs.DefaultFuncs), nil
}

// 加载基础布局
func (e *Engine) loadBaseLayout(t *template.Template) error {
	baseofPath := filepath.Join(e.layoutDir, "_default/baseof.html")
	baseContent, err := os.ReadFile(baseofPath)
	if err != nil {
		return fmt.Errorf("failed to read baseof.html: %w", err)
	}

	if _, err := t.New("_default/baseof.html").Parse(string(baseContent)); err != nil {
		return fmt.Errorf("failed to parse baseof.html: %w", err)
	}
	return nil
}

// 加载组件
func (e *Engine) loadComponents(t *template.Template) error {
	componentsDir := filepath.Join(e.layoutDir, "components")
	return filepath.Walk(componentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".html") {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read component: %w", err)
		}

		if _, err := t.Parse(string(content)); err != nil {
			return fmt.Errorf("failed to parse component: %w", err)
		}
		return nil
	})
}

// 加载页面模板
func (e *Engine) loadPageTemplate(t *template.Template, kind, section string) error {
	templateName := e.lookupTemplate(kind, section)
	content, err := os.ReadFile(filepath.Join(e.layoutDir, templateName))
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	if _, err := t.New(templateName).Parse(string(content)); err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}
	return nil
}

// 处理模板数据
func (e *Engine) processTemplateData(data interface{}, kind, section string) (interface{}, error) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return data, nil
	}

	if _, exists := m["Site"]; !exists {
		m["Site"] = e.config
	}

	if e.buildMode {
		depth := e.calculatePathDepth(kind, section)
		prefix := e.buildPathPrefix(depth)
		e.setStaticPaths(m, prefix)
		e.processImagePaths(m, prefix)
	} else {
		m["StaticPath"] = "/static"
		m["ImagesPath"] = "/static/images"
	}

	m["BuildMode"] = e.buildMode
	return m, nil
}

// 计算路径深度
func (e *Engine) calculatePathDepth(kind, section string) int {
	depth := 0
	if section != "" {
		depth++
	}
	if kind != "index" {
		depth++
	}
	return depth
}

// 构建路径前缀
func (e *Engine) buildPathPrefix(depth int) string {
	if depth == 0 {
		return "./"
	}
	return strings.Repeat("../", depth)
}

// 设置静态资源路径
func (e *Engine) setStaticPaths(m map[string]interface{}, prefix string) {
	m["StaticPath"] = prefix + "static"
	m["ImagesPath"] = prefix + "static/images"
}

// 处理图片路径
func (e *Engine) processImagePaths(m map[string]interface{}, prefix string) {
	if site, ok := m["Site"].(*config.Config); ok && e.buildMode {
		if strings.HasPrefix(site.Author.Avatar, "/") {
			site.Author.Avatar = prefix + strings.TrimPrefix(site.Author.Avatar, "/")
		}
		for i := range site.Author.Projects {
			if strings.HasPrefix(site.Author.Projects[i].Image, "/") {
				site.Author.Projects[i].Image = prefix + strings.TrimPrefix(site.Author.Projects[i].Image, "/")
			}
		}
	}
}

// lookupTemplate 按照优先级查找模板
func (e *Engine) lookupTemplate(kind, section string) string {
	// 模板查找顺序
	lookupOrder := []string{
		filepath.Join(section, kind+".html"),        // posts/single.html
		filepath.Join("_default", kind+".html"),     // _default/single.html
		filepath.Join(section, "_"+kind+".html"),    // posts/_single.html
		filepath.Join("_default", "_"+kind+".html"), // _default/_single.html
	}

	// 检查文件是否存在
	for _, name := range lookupOrder {
		path := filepath.Join(e.layoutDir, name)
		if _, err := os.Stat(path); err == nil {
			return name
		}
	}

	// 默认返回
	if kind == "index" {
		return "index.html"
	}
	return filepath.Join("_default", kind+".html")
}

// RenderHome 渲染首页
func (e *Engine) RenderHome(data map[string]interface{}) (string, error) {
	return e.render("index", "", data)
}

// RenderList 渲染列表页
func (e *Engine) RenderList(data map[string]interface{}) (string, error) {
	return e.render("list", "posts", data)
}

// RenderSingle 渲染单页
func (e *Engine) RenderSingle(data map[string]interface{}) (string, error) {
	return e.render("single", "posts", data)
}

// RenderTags 渲染标签云页面
func (e *Engine) RenderTags(data map[string]interface{}) (string, error) {
	return e.render("tags", "", data)
}

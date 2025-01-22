package sitemap

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jiangjiax/stars/internal/config"
	"github.com/jiangjiax/stars/internal/post"
)

type Sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq"`
	Priority   float64 `xml:"priority"`
}

// Generator sitemap 生成器
type Generator struct {
	config *config.Config
	posts  []*post.Post
}

// New 创建新的 sitemap 生成器
func New(cfg *config.Config, posts []*post.Post) *Generator {
	return &Generator{
		config: cfg,
		posts:  posts,
	}
}

// Generate 生成 sitemap.xml
func (g *Generator) Generate(outputDir string) error {
	sitemap := &Sitemap{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  g.collectURLs(),
	}

	// 确保输出目录存在
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 创建文件
	outputPath := filepath.Join(outputDir, "sitemap.xml")
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create sitemap file: %w", err)
	}
	defer file.Close()

	// 写入 XML 头
	if _, err := file.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}

	// 创建编码器并设置缩进
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	// 编码并写入文件
	if err := encoder.Encode(sitemap); err != nil {
		return fmt.Errorf("failed to encode sitemap: %w", err)
	}

	return nil
}

// collectURLs 收集所有 URL
func (g *Generator) collectURLs() []URL {
	var urls []URL
	baseURL := g.config.BaseURL

	// 添加首页
	urls = append(urls, URL{
		Loc:        baseURL,
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	// 添加文章列表页
	urls = append(urls, URL{
		Loc:        baseURL + "posts",
		ChangeFreq: "daily",
		Priority:   0.9,
	})

	// 添加标签页
	urls = append(urls, URL{
		Loc:        baseURL + "tags",
		ChangeFreq: "weekly",
		Priority:   0.8,
	})

	// 添加所有文章页面
	for _, p := range g.posts {
		if !p.Draft {
			urls = append(urls, URL{
				Loc:        baseURL + "posts/" + p.Slug,
				LastMod:    p.Date.Format("2006-01-02"),
				ChangeFreq: "monthly",
				Priority:   0.7,
			})
		}
	}

	return urls
} 
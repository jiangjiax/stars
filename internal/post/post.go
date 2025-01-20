package post

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-emoji"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	ghtml "github.com/yuin/goldmark/renderer/html"
	"github.com/jiangjiax/stars/internal/config"
)

var cfg *config.Config

func init() {
	// 从环境变量或默认路径加载配置
	var err error
	cfg, err = config.LoadConfig("config.yaml")
	if err != nil {
		// 如果加载失败，使用空配置
		cfg = &config.Config{}
	}
}

// PostMeta 文章元数据
type PostMeta struct {
	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	Description string    `yaml:"description"`
	Tags        []string  `yaml:"tags"`
	Slug        string    `yaml:"slug"`
	Series      string    `yaml:"series"`
	SeriesOrder int       `yaml:"seriesOrder"`
	ReadingTime int       `yaml:"readingTime"`
}

// Post 完整文章
type Post struct {
	PostMeta                      // 嵌入元数据
	Content         template.HTML // 文章内容
	RawContent      string        // 原始 Markdown
	Title           string        `yaml:"title"`
	Slug            string        `yaml:"slug"`
	Date            time.Time     `yaml:"date"`
	Description     string        `yaml:"description"`
	Tags            []string      `yaml:"tags"`
	Series          string        `yaml:"series"`
	SeriesOrder     int           `yaml:"seriesOrder"`
	Draft           bool          `yaml:"draft"`
	TableOfContents []*TableOfContentsItem
	ReadingTime     int                  `yaml:"readingTime"`
	Verification    *config.Verification `yaml:"verification"`
	FilePath        string               `yaml:"-"`
}

type TableOfContentsItem struct {
	Title    string
	ID       string
	Level    int
	Children []*TableOfContentsItem
}

// 创建一个全局的 markdown 解析器
var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,            // GitHub Flavored Markdown
		extension.Footnote,       // 脚注支持
		extension.DefinitionList, // 定义列表
		extension.Table,          // 表格支持
		extension.Strikethrough,  // 删除线
		extension.Linkify,        // 自链接
		extension.TaskList,       // 任务列表
		meta.Meta,                // Front Matter 持
		emoji.Emoji,              // Emoji 支持
		highlighting.NewHighlighting( // 代码高亮
			highlighting.WithStyle("monokai"),
			highlighting.WithGuessLanguage(true),
		),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(), // 自动生成标题 ID
		parser.WithAttribute(),     // 属性支持
	),
	goldmark.WithRendererOptions(
		ghtml.WithHardWraps(), // 保留换行
		ghtml.WithXHTML(),     // 使用 XHTML
		ghtml.WithUnsafe(),    // 允许原始 HTML
	),
)

// ParsePosts 解析指定目录的所有文章
func ParsePosts(contentDir string) ([]*Post, error) {
	var posts []*Post

	// 直接使用传入的目录路径
	err := filepath.Walk(contentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理 .md 文件
		if info.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		// 解析文章
		post, err := ParsePost(path)
		if err != nil {
			return fmt.Errorf("failed to parse post %s: %w", path, err)
		}

		// 只添加非草稿文章
		if !post.Draft {
			posts = append(posts, post)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// ParsePost 解析单个文章文件
func ParsePost(filePath string) (*Post, error) {
	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 创建新的文章实例
	post := &Post{
		FilePath:   filePath,
		RawContent: string(content),
	}

	// 分离 Front Matter 和内容
	parts := bytes.Split(content, []byte("---\n"))
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid post format: missing front matter")
	}

	// 解析 Front Matter 到 post 结构体
	if err := yaml.Unmarshal(parts[1], &post); err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %w", err)
	}

	// 设置文章内容
	post.RawContent = string(bytes.Join(parts[2:], []byte("---\n")))

	// 渲染 Markdown 内容
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert([]byte(post.RawContent), &buf, parser.WithContext(context)); err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}

	renderedContent := buf.String()
	post.Content = template.HTML(renderedContent)

	// 从渲染后的 HTML 中提取标题和 ID 来生成目录
	toc := generateTOCFromHTML(renderedContent)
	post.TableOfContents = toc

	// 如果没有设置 slug，使用标题生成
	if post.Slug == "" {
		post.Slug = slugify(post.Title)
	}

	// 计算阅读时间
	post.ReadingTime = calculateReadingTime(post.RawContent)

	// 如果没有设置作者地址，从配置文件读取
	if post.Verification != nil && post.Verification.Author == "" {
		post.Verification.Author = cfg.Author.WalletAddress
	}

	if post.Verification == nil {
		post.Verification = &config.Verification{}
	}

	return post, nil
}

// ParseContent 解析文章内容
func ParseContent(content string) (*Post, error) {
	// 移除可能的 BOM 头
	content = removeBOM(content)

	// 规范化换行符
	content = normalizeNewlines(content)

	// 分离 Front Matter 和文章内容
	parts := bytes.Split([]byte(content), []byte("---\n"))
	if len(parts) < 3 {
		// 尝试使用 CRLF
		parts = bytes.Split([]byte(content), []byte("---\r\n"))
		if len(parts) < 3 {
			return nil, fmt.Errorf("invalid post format: missing front matter")
		}
	}

	// 跳过第一个空部分
	frontMatter := parts[1]
	content = string(bytes.Join(parts[2:], []byte("---\n")))

	var post Post
	if err := yaml.Unmarshal(frontMatter, &post); err != nil {
		return nil, fmt.Errorf("failed to parse front matter: %w", err)
	}

	// 先渲染 Markdown 内容
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := md.Convert([]byte(content), &buf, parser.WithContext(context)); err != nil {
		return nil, fmt.Errorf("failed to render markdown: %w", err)
	}

	renderedContent := buf.String()
	post.Content = template.HTML(renderedContent)

	// 从渲染后的 HTML 中提取标题和 ID 来生成目录
	toc := generateTOCFromHTML(renderedContent)
	post.TableOfContents = toc

	// 如果没有设置 slug，使用标题生成
	if post.Slug == "" {
		post.Slug = slugify(post.Title)
	}

	// 计算阅读时间
	post.ReadingTime = calculateReadingTime(content)

	// 保存原始内容
	post.RawContent = content

	// 如果没有设置作者地址，从配置文件读取
	if post.Verification != nil && post.Verification.Author == "" {
		post.Verification.Author = cfg.Author.WalletAddress
	}

	return &post, nil
}

// slugify 将标题转换为 URL 友好的格式
func slugify(title string) string {
	// 如果标题是中文，使用 pinyin 或其他方式转换
	// 这我们先用一个简单的方案：为中文标题生成一个唯一的 ID
	if regexp.MustCompile(`[\p{Han}]`).MatchString(title) {
		// 移除空格并转小写
		title = strings.ToLower(strings.TrimSpace(title))
		// 将标题转换为 bytes 并进行 base64 编的前8位作为 ID
		hash := fmt.Sprintf("%x", md5.Sum([]byte(title)))[:8]
		return hash
	}

	// 文标题使用原来的逻辑
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(slug, "")
	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

// calculateReadingTime 计算文章阅读时间
func calculateReadingTime(content string) int {
	// 移除代码块和 HTML 标签
	content = regexp.MustCompile("```[\\s\\S]*?```").ReplaceAllString(content, "")
	content = regexp.MustCompile("<[^>]+>").ReplaceAllString(content, "")
	
	// 分别计算中文和英文字数
	chineseCount := len(regexp.MustCompile(`[\p{Han}]`).FindAllString(content, -1))
	
	// 移除中文字符后计算英文单词数
	contentWithoutChinese := regexp.MustCompile(`[\p{Han}]`).ReplaceAllString(content, "")
	words := strings.Fields(contentWithoutChinese)
	englishCount := len(words)
	
	// 中文阅读速度：约每分钟 300 字
	// 英文阅读速度：约每分钟 200 词
	chineseTime := (chineseCount + 299) / 300
	englishTime := (englishCount + 199) / 200
	
	// 取较大值，并确保至少返回 1 分钟
	readingTime := max(chineseTime, englishTime)
	if readingTime < 1 {
		readingTime = 1
	}
	
	return readingTime
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// generateTOCFromHTML 从渲染后的 HTML 中提取标题和 ID 来生成目录
func generateTOCFromHTML(html string) []*TableOfContentsItem {
	// 匹配 HTML 标题标签及其 ID 属性，使用命名捕获组
	re := regexp.MustCompile(`<h([1-6])[^>]*id="([^"]+)"[^>]*>(.*?)</h[1-6]>`)
	matches := re.FindAllStringSubmatch(html, -1)

	if len(matches) == 0 {
		return nil
	}

	var items []*TableOfContentsItem
	var stack []*TableOfContentsItem
	currentLevel := 0

	for _, match := range matches {
		level, _ := strconv.Atoi(match[1])
		id := match[2]
		// 移除标题中可能的 HTML 标签
		title := regexp.MustCompile(`<[^>]+>`).ReplaceAllString(match[3], "")

		item := &TableOfContentsItem{
			Title: title,
			ID:    id,
			Level: level,
		}

		// 处理层级关系
		if level > currentLevel {
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, item)
			} else {
				items = append(items, item)
			}
			stack = append(stack, item)
		} else if level == currentLevel {
			if len(stack) > 1 {
				parent := stack[len(stack)-2]
				parent.Children = append(parent.Children, item)
			} else {
				items = append(items, item)
			}
			stack[len(stack)-1] = item
		} else {
			// 回溯到正确的层级
			for len(stack) > 0 && stack[len(stack)-1].Level >= level {
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				parent.Children = append(parent.Children, item)
			} else {
				items = append(items, item)
			}
			stack = append(stack, item)
		}
		currentLevel = level
	}

	return items
}

// 添加 UpdateFrontMatter 函数
func UpdateFrontMatter(content []byte, updates map[string]string) error {
	// 分离 Front Matter 和文章内容
	parts := bytes.Split(content, []byte("---\n"))
	if len(parts) < 3 {
		return fmt.Errorf("invalid post format: missing front matter")
	}

	// 析 Front Matter
	var frontMatter map[string]interface{}
	if err := yaml.Unmarshal(parts[1], &frontMatter); err != nil {
		return fmt.Errorf("failed to parse front matter: %w", err)
	}

	// 更新字段
	for key, value := range updates {
		// 处理布尔值
		if key == "draft" {
			frontMatter[key] = value == "true"
		} else {
			frontMatter[key] = value
		}
	}

	// 重新序列化 Front Matter
	updatedFrontMatter, err := yaml.Marshal(frontMatter)
	if err != nil {
		return fmt.Errorf("failed to marshal updated front matter: %w", err)
	}

	// 新组合文件内容
	var result bytes.Buffer
	result.WriteString("---\n")
	result.Write(updatedFrontMatter)
	result.WriteString("---\n")
	result.Write(bytes.Join(parts[2:], []byte("---\n")))

	// 写回文件
	if err := os.WriteFile(string(content), result.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write updated content: %w", err)
	}

	return nil
}

// SaveMetadata 保存更新后的元数据到文件
func (p *Post) SaveMetadata() error {
	// 检查文件路径是否为空
	if p.FilePath == "" {
		return fmt.Errorf("file path is empty")
	}

	// 在保存元数据时自动计算内容哈希
	if p.Verification == nil {
		p.Verification = &config.Verification{}
	}
	p.Verification.ContentHash = p.calculateContentHash()

	// 读取原始文件
	content, err := os.ReadFile(p.FilePath)
	if err != nil {
		return err
	}

	// 分离前置元数据和内容
	parts := bytes.Split(content, []byte("---\n"))
	if len(parts) < 3 {
		return fmt.Errorf("invalid post format")
	}

	// 更新元数据
	metadata := map[string]interface{}{
		"title":        p.Title,
		"date":         p.Date,
		"description":  p.Description,
		"tags":         p.Tags,
		"slug":         p.Slug,
		"series":       p.Series,
		"seriesOrder":  p.SeriesOrder,
		"verification": p.Verification, // 添加验证信息到元数据
	}

	// 序列化元数据
	metadataBytes, err := yaml.Marshal(metadata)
	if err != nil {
		return err
	}

	// 组合新文件内容
	newContent := []byte("---\n")
	newContent = append(newContent, metadataBytes...)
	newContent = append(newContent, []byte("---\n")...)
	newContent = append(newContent, parts[2]...)

	// 写回文件
	return os.WriteFile(p.FilePath, newContent, 0644)
}

// 移除 BOM 头
func removeBOM(content string) string {
	if strings.HasPrefix(content, "\xEF\xBB\xBF") {
		return content[3:]
	}
	return content
}

// 规范化换行符
func normalizeNewlines(content string) string {
	// 将 Windows 格的换行符转换为 Unix 风格
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func (p *Post) calculateContentHash() string {
	hasher := sha3.NewLegacyKeccak256()

	// 组合所有需要计算哈希的内容
	content := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		Title:   p.Title,
		Content: p.RawContent,
	}

	// 序列化为 JSON 以确保一致性
	data, _ := json.Marshal(content)
	hasher.Write(data)
	hash := hasher.Sum(nil)
	return "0x" + hex.EncodeToString(hash) + p.Verification.NFT.Version
}

func (p *Post) parseMetadata(meta map[string]interface{}) error {
	if verification, ok := meta["verification"].(map[string]interface{}); ok {
		p.Verification = &config.Verification{}

		// 处理 NFT 配置
		if nftConfig, ok := verification["nft"].(map[string]interface{}); ok {
			p.Verification.NFT = &config.NFTConfig{}

			// 解析配置
			if price, ok := nftConfig["price"].(string); ok {
				p.Verification.NFT.Price = price
			}
			if maxSupply, ok := nftConfig["maxSupply"].(int); ok {
				p.Verification.NFT.MaxSupply = maxSupply
			}
			if royaltyFee, ok := nftConfig["royaltyFee"].(int); ok {
				p.Verification.NFT.RoyaltyFee = royaltyFee
			}
			if onePerAddress, ok := nftConfig["onePerAddress"].(bool); ok {
				p.Verification.NFT.OnePerAddress = onePerAddress
			}
			if version, ok := nftConfig["version"].(string); ok {
				p.Verification.NFT.Version = version
			}

			// 验证参数
			if err := p.Verification.NFT.ValidateNFTConfig(); err != nil {
				// 记录警告信息
				log.Printf("Warning: Invalid NFT config in post %s: %v, using default values",
					p.Slug, err)
				// 使用默认配置
				p.Verification.NFT = config.GetDefaultNFTConfig()
			}
		}

		// 如果没配置 NFT，使用默认配置
		if p.Verification.NFT == nil {
			p.Verification.NFT = config.GetDefaultNFTConfig()
		}
	}

	return nil
}

// ContentChanged 检查内容是否变化
func (p *Post) ContentChanged() bool {
	if p.Verification == nil || p.Verification.ContentHash == "" {
		return true
	}
	currentHash := p.calculateContentHash()
	return currentHash != p.Verification.ContentHash
}

// UpdateContentHash 更新内容哈希并保存
func (p *Post) UpdateContentHash() error {
	if p.Verification == nil {
		p.Verification = &config.Verification{}
	}
	p.Verification.ContentHash = p.calculateContentHash()
	return p.SaveMetadata()
}
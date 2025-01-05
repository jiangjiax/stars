package rss

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"stars/internal/post"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel
}

type Channel struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	Language      string `xml:"language"`
	LastBuildDate string `xml:"lastBuildDate"`
	Items         []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
}

func GenerateFeed(posts []*post.Post, siteTitle string, baseURL string) (*RSS, error) {
	feed := &RSS{
		Version: "2.0",
		Channel: Channel{
			Title:         siteTitle,
			Link:          baseURL,
			Description:   "最新文章更新",
			Language:      "zh-CN",
			LastBuildDate: time.Now().Format(time.RFC1123Z),
		},
	}

	// 默认显示最近10篇
	limit := 10

	// 添加文章
	for i, p := range posts {
		if i >= limit {
			break
		}

		item := Item{
			Title:       p.Title,
			Link:        baseURL + "/" + p.Slug,
			Description: p.Description,
			PubDate:     p.Date.Format(time.RFC1123Z),
			GUID:        baseURL + "/" + p.Slug,
		}
		feed.Channel.Items = append(feed.Channel.Items, item)
	}

	return feed, nil
}

// GenerateAndSaveFeed 生成 RSS feed 并保存到文件
func GenerateAndSaveFeed(posts []*post.Post, siteTitle string, baseURL string, outputPath string) error {
	// 生成 feed
	feed, err := GenerateFeed(posts, siteTitle, baseURL)
	if err != nil {
		return fmt.Errorf("failed to generate feed: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create feed file: %w", err)
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
	if err := encoder.Encode(feed); err != nil {
		return fmt.Errorf("failed to encode feed: %w", err)
	}

	return nil
}

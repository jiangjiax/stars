package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/jiangjiax/stars/internal/config"
	"github.com/jiangjiax/stars/internal/post"
)

var (
	postSeries string
	postOrder  int
	postTags   []string
	postSlug   string
	postDraft  bool
	postDesc   string
)

var postCmd = &cobra.Command{
	Use:   "post [title]",
	Short: "Create a new post",
	Long: `Create a new post with the given title.
This command will create a new markdown file in the content/posts directory.`,
	Example: `  stars post "My First Post"
  stars post "Hello World" --series "Getting Started" --order 1
  stars post "New Feature" --tags "feature,update"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]

		// 获取当前工作目录
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// 检查是否是有效的 Stars 项目
		configFile := filepath.Join(projectDir, "config.yaml")
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			return fmt.Errorf("not a valid Stars project (config.yaml not found)")
		}

		// 生成文章文件名
		fileName := generateFileName(title)
		if postSlug != "" {
			fileName = postSlug + ".md"
		}

		// 创建文章目录
		postsDir := filepath.Join(projectDir, "content", "posts")
		if err := os.MkdirAll(postsDir, 0755); err != nil {
			return fmt.Errorf("failed to create posts directory: %w", err)
		}

		// 检查文件是否已存在
		filePath := filepath.Join(postsDir, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			return fmt.Errorf("post file already exists: %s", fileName)
		}

		// 创建文章文件
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create post file: %w", err)
		}
		defer file.Close()

		// 生成文章内容
		content, err := generatePostContent(title)
		if err != nil {
			return fmt.Errorf("failed to generate post content: %w", err)
		}

		// 写入文件
		if _, err := file.WriteString(content); err != nil {
			return fmt.Errorf("failed to write post content: %w", err)
		}

		fmt.Printf("Created new post: %s\n", filePath)
		return nil
	},
}

var postTemplate = `---
title: "{{ .Title }}"
date: {{ .Date.Format "2006-01-02" }}
{{- if .Description }}
description: "{{ .Description }}"
{{- end }}
{{- if .Series }}
series: "{{ .Series }}"
{{- end }}
{{- if .SeriesOrder }}
seriesOrder: {{ .SeriesOrder }}
{{- end }}
{{- if .Tags }}
tags: [{{ range $i, $tag := .Tags }}{{if $i}}, {{end}}"{{ $tag }}"{{ end }}]
{{- end }}
{{- if .Slug }}
slug: "{{ .Slug }}"
{{- end }}
{{- if .Draft }}
draft: true
{{- end }}
{{- if .Verification }}
verification:
  arweaveId: "{{ .Verification.ArweaveId }}"
  nftContract: "{{ .Verification.NftContract }}"
  author: "{{ .Verification.Author }}"
  contentHash: "{{ .Verification.ContentHash }}"
  nft:
    price: "{{ .Verification.NFT.Price }}"
    maxSupply: {{ .Verification.NFT.MaxSupply }}
    royaltyFee: {{ .Verification.NFT.RoyaltyFee }}
    onePerAddress: {{ .Verification.NFT.OnePerAddress }}
    version: "{{ .Verification.NFT.Version }}"
    chainId: {{ .Verification.NFT.ChainId }}
    tokenSymbol: "{{ .Verification.NFT.TokenSymbol }}"
{{- end }}
---

在这里写下你的文章内容...
`

func generateFileName(title string) string {
	// 移除特殊字符，将空格替换为连字符
	fileName := strings.ToLower(title)
	fileName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, fileName)
	fileName = strings.ReplaceAll(fileName, " ", "-")
	fileName = strings.Trim(fileName, "-")

	// 添加日期前缀和扩展名
	return fmt.Sprintf("%s-%s.md", time.Now().Format("2006-01-02"), fileName)
}

func generatePostContent(title string) (string, error) {
	// 获取当前工作目录
	projectDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	// 加载配置文件
	cfg, err := config.LoadConfig(filepath.Join(projectDir, "config.yaml"))
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// 创建新的 Post 实例
	p := &post.Post{
		Title:       title,
		Date:        time.Now(),
		Description: postDesc,
		Series:      postSeries,
		SeriesOrder: postOrder,
		Tags:        postTags,
		Draft:       postDraft,
		Slug:        postSlug,
		Verification: cfg.Verification,  // 从配置文件获取默认验证信息
	}

	// 解析模板
	tmpl, err := template.New("post").Parse(postTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// 渲染模板
	var buf strings.Builder
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("generatePostContent failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func init() {
	rootCmd.AddCommand(postCmd)

	// 添加命令行标志
	postCmd.Flags().StringVar(&postSeries, "series", "", "series name for the post")
	postCmd.Flags().IntVar(&postOrder, "order", 0, "order in the series")
	postCmd.Flags().StringSliceVar(&postTags, "tags", []string{}, "comma-separated list of tags")
	postCmd.Flags().StringVar(&postSlug, "slug", "", "custom URL slug")
	postCmd.Flags().BoolVar(&postDraft, "draft", false, "mark post as draft")
	postCmd.Flags().StringVar(&postDesc, "desc", "", "post description")
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"github.com/jiangjiax/stars/internal/post"
)

var draftCmd = &cobra.Command{
	Use:   "draft",
	Short: "Manage draft posts",
	Long:  `Create, list, and publish draft posts.`,
}

var draftListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all draft posts",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		posts, err := post.ParsePosts(filepath.Join(projectDir, "content"))
		if err != nil {
			return err
		}

		fmt.Println("Draft posts:")
		for _, p := range posts {
			if p.Draft {
				fmt.Printf("- %s (%s)\n", p.Title, p.Date.Format("2006-01-02"))
			}
		}
		return nil
	},
}

var draftPublishCmd = &cobra.Command{
	Use:   "publish [slug]",
	Short: "Publish a draft post",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slug := args[0]
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// 读取文章文件
		filePath := filepath.Join(projectDir, "content", "posts", slug+".md")
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read post file: %w", err)
		}

		// 更新 front matter，移除 draft: true
		post.UpdateFrontMatter(content, map[string]string{"draft": "false"})

		fmt.Printf("Published post: %s\n", slug)
		return nil
	},
}

func init() {
	draftCmd.AddCommand(draftListCmd)
	draftCmd.AddCommand(draftPublishCmd)
	rootCmd.AddCommand(draftCmd)
}

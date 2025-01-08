package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jiangjiax/stars/internal/generator"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the static website",
	Long: `Build command generates a static website from your Stars project.
It processes all content, applies themes, and creates a deployable website.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()

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

		// 创建生成器实例
		gen, err := generator.New(projectDir, false)
		if err != nil {
			return fmt.Errorf("failed to create generator: %w", err)
		}

		// 生成静态网站
		fmt.Println("Generating static website...")
		if err := gen.Build(); err != nil {
			return fmt.Errorf("failed to build website: %w", err)
		}

		elapsed := time.Since(start)
		fmt.Printf("\nBuild completed in %s\n", elapsed)
		fmt.Println("Your website is ready in the 'public' directory")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

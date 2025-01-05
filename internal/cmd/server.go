package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"stars/internal/config"
	"stars/internal/server"
)

var (
	port int
)

type ServerCommand struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start development server",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// 加载配置
		cfg, err := config.LoadConfig(filepath.Join(projectDir, "config.yaml"))
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// 创建带取消的上下文
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		// 创建并启动开发服务器
		srv, err := server.New(projectDir, cfg, port)
		if err != nil {
			return fmt.Errorf("failed to create server: %w", err)
		}

		// 处理优雅退出
		go func() {
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			<-sigCh
			fmt.Println("\nShutting down server...")
			cancel()
		}()

		fmt.Printf("Starting development server on http://localhost:%d\n", port)

		// 启动服务器
		errCh := make(chan error, 1)
		go func() {
			errCh <- srv.Start()
		}()

		// 等待服务器退出或上下文取消
		select {
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&port, "port", "p", 1313, "port to run the server on")
}

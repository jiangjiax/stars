package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stars",
	Short: "Stars is a Web3-enabled static blog generator",
	Long: `Stars (繁星) is a modern static blog generator with Web3 features,
built with Go. It allows you to create and deploy decentralized blogs.`,
}

func Execute() error {
	return rootCmd.Execute()
}

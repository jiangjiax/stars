package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"stars/internal/theme"
)

var themeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage themes",
	Long:  `Create, list, and manage themes for your Stars blog.`,
}

var themeNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		manager := theme.New(projectDir)
		if err := manager.Create(name); err != nil {
			return fmt.Errorf("failed to create theme: %w", err)
		}

		fmt.Printf("Created new theme: %s\n", name)
		return nil
	},
}

var themeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		manager := theme.New(projectDir)
		themes, err := manager.List()
		if err != nil {
			return fmt.Errorf("failed to list themes: %w", err)
		}

		// 使用 tabwriter 格式化输出
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tVERSION\tAUTHOR\tDESCRIPTION")
		for _, t := range themes {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.Name, t.Version, t.Author, t.Description)
		}
		w.Flush()

		return nil
	},
}

var themeUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Switch to a theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		manager := theme.New(projectDir)
		if err := manager.Use(name); err != nil {
			return fmt.Errorf("failed to switch theme: %w", err)
		}

		fmt.Printf("Switched to theme: %s\n", name)
		return nil
	},
}

var themeInstallCmd = &cobra.Command{
	Use:   "install [repo]",
	Short: "Install theme from git repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := args[0]

		projectDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		manager := theme.New(projectDir)
		if err := manager.Install(repo); err != nil {
			return fmt.Errorf("failed to install theme: %w", err)
		}

		fmt.Printf("Successfully installed theme from %s\n", repo)
		return nil
	},
}

func init() {
	themeCmd.AddCommand(themeNewCmd)
	themeCmd.AddCommand(themeListCmd)
	themeCmd.AddCommand(themeUseCmd)
	themeCmd.AddCommand(themeInstallCmd)
	rootCmd.AddCommand(themeCmd)
}

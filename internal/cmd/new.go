package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"stars/internal/generator"
)

var newCmd = &cobra.Command{
	Use:   "new [path]",
	Short: "Create a new Stars blog",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Convert to absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path: %w", err)
		}

		// Check if directory already exists
		if _, err := os.Stat(absPath); !os.IsNotExist(err) {
			return fmt.Errorf("directory %s already exists", absPath)
		}

		// Create new project
		project, err := generator.New(absPath)
		if err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}

		// Generate project files
		if err := project.Generate(); err != nil {
			return fmt.Errorf("failed to generate project: %w", err)
		}

		fmt.Printf("Successfully created new Stars blog at %s\n", absPath)
		fmt.Println("\nNext steps:")
		fmt.Printf("  cd %s\n", path)
		fmt.Println("  stars server    # Start development server")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

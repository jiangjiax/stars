package asset

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Pipeline struct {
	projectDir string
	assetMap   map[string]string // 用于存储资源映射
	theme      string
}

func New(projectDir string, theme string) *Pipeline {
	return &Pipeline{
		projectDir: projectDir,
		assetMap:   make(map[string]string),
		theme:      theme,
	}
}

// BuildAssets 构建所有资源
func (p *Pipeline) BuildAssets() error {
	// 先构建 CSS
	if err := p.BuildCSS(); err != nil {
		return fmt.Errorf("failed to build CSS: %w", err)
	}

	// 构建 JS
	fmt.Println("Building JavaScript...")
	jsCmd := exec.Command("npm", "run", "build:js")
	jsCmd.Dir = filepath.Join(p.projectDir, "themes", p.theme)
	jsCmd.Stdout = os.Stdout
	jsCmd.Stderr = os.Stderr
	if err := jsCmd.Run(); err != nil {
		return fmt.Errorf("failed to build JavaScript: %w", err)
	}

	// 读取生成的资源映射
	assetManifest := filepath.Join(p.projectDir, "themes", p.theme, "static/dist/manifest.json")
	data, err := os.ReadFile(assetManifest)
	if err != nil {
		return fmt.Errorf("failed to read asset manifest: %w", err)
	}

	if err := json.Unmarshal(data, &p.assetMap); err != nil {
		return fmt.Errorf("failed to parse asset manifest: %w", err)
	}

	return nil
}

// BuildCSS 构建 CSS 文件
func (p *Pipeline) BuildCSS() error {
	// 检查必要文件是否存在
	files := []string{
		filepath.Join("themes", p.theme, "package.json"),
		filepath.Join("themes", p.theme, "tailwind.config.js"),
		filepath.Join("themes", p.theme, "static", "css", "styles.css"),
	}

	for _, file := range files {
		if _, err := os.Stat(filepath.Join(p.projectDir, file)); os.IsNotExist(err) {
			return fmt.Errorf("required file not found: %s, dir: %s", file, p.projectDir)
		}
	}

	// 检查 npm 是否安装
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("node.js and npm are required. Please install them first: %w", err)
	}

	// 检查 node_modules 是否存在且完整
	nodeModulesPath := filepath.Join(p.projectDir, "themes", p.theme, "node_modules")
	tailwindPath := filepath.Join(nodeModulesPath, "tailwindcss")
	if _, err := os.Stat(tailwindPath); os.IsNotExist(err) {
		// 安装依赖
		fmt.Println("Installing dependencies...")
		installCmd := exec.Command("npm", "install")
		installCmd.Dir = filepath.Join(p.projectDir, "themes", p.theme)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install dependencies: %w", err)
		}
	} else {
		fmt.Println("Dependencies already installed, skipping npm install")
	}

	// 构建 CSS
	fmt.Println("Building CSS...")
	buildCmd := exec.Command("npm", "run", "build:css")
	buildCmd.Dir = filepath.Join(p.projectDir, "themes", p.theme)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build CSS: %w", err)
	}

	return nil
}

// 添加获取资源路径的方法
func (p *Pipeline) GetAssetPath(name string) string {
	if hash, ok := p.assetMap[name]; ok {
		return hash
	}
	return name
}

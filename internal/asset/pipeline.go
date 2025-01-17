package asset

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	// 在构建新资源之前清理旧文件
	if err := p.CleanOldAssets(); err != nil {
		return err
	}

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

	// 处理 ABI 文件
	if err := p.processABIFiles(); err != nil {
		return fmt.Errorf("failed to process ABI files: %w", err)
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

// CleanOldAssets 清理旧的资源文件
func (p *Pipeline) CleanOldAssets() error {
	// 获取 dist 目录
	distDir := filepath.Join(p.projectDir, "themes", p.theme, "static/dist")

	// 读取目录下所有文件
	entries, err := os.ReadDir(distDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // dist 目录不存在时直接返回
		}
		return fmt.Errorf("failed to read dist directory: %w", err)
	}

	// 删除所有 .js 文件
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".js") {
			filePath := filepath.Join(distDir, entry.Name())
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove old asset file %s: %w", filePath, err)
			}
		}
	}

	return nil
}

// 添加新方法处理 ABI 文件
func (p *Pipeline) processABIFiles() error {
	// ABI 文件源目录
	abiDir := filepath.Join(p.projectDir, "themes", p.theme, "static", "abi")
	
	// 确保目标目录存在
	distDir := filepath.Join(p.projectDir, "themes", p.theme, "static", "dist", "abi")
	if err := os.MkdirAll(distDir, 0755); err != nil {
		return err
	}

	// 遍历处理所有 ABI 文件
	entries, err := os.ReadDir(abiDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			// 读取文件内容
			content, err := os.ReadFile(filepath.Join(abiDir, entry.Name()))
			if err != nil {
				return err
			}

			// 计算内容哈希
			hashBytes := sha256.Sum256(content)
			hash := fmt.Sprintf("%x", hashBytes[:16])
			
			// 生成新文件名
			baseName := strings.TrimSuffix(entry.Name(), ".json")
			newName := fmt.Sprintf("%s.%s.json", baseName, hash)
			
			// 写入新文件
			if err := os.WriteFile(filepath.Join(distDir, newName), content, 0644); err != nil {
				return err
			}

			// 更新资源映射
			p.assetMap[entry.Name()] = fmt.Sprintf("abi/%s", newName)
		}
	}

	return nil
}
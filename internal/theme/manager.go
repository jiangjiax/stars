package theme

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Manager 主题管理器
type Manager struct {
	projectDir string
	themesDir  string
}

// Theme 主题配置
type Theme struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Author      string   `yaml:"author"`
	Description string   `yaml:"description"`
	Homepage    string   `yaml:"homepage"`
	License     string   `yaml:"license"`
	Tags        []string `yaml:"tags"`
}

// New 创建主题管理器
func New(projectDir string) *Manager {
	return &Manager{
		projectDir: projectDir,
		themesDir:  filepath.Join(projectDir, "themes"),
	}
}

// List 列出所有已安装的主题
func (m *Manager) List() ([]Theme, error) {
	var themes []Theme

	entries, err := os.ReadDir(m.themesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read themes directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		theme, err := m.loadTheme(entry.Name())
		if err != nil {
			return nil, err
		}
		themes = append(themes, *theme)
	}

	return themes, nil
}

// Create 创建新主题
func (m *Manager) Create(name string) error {
	themeDir := filepath.Join(m.themesDir, name)

	// 检查主题是否已存在
	if _, err := os.Stat(themeDir); !os.IsNotExist(err) {
		return fmt.Errorf("theme %s already exists", name)
	}

	// 创建主题目录结构
	dirs := []string{
		"layouts",
		"layouts/_default",
		"layouts/partials",
		"layouts/components",
		"static/css",
		"static/js",
		"static/images",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(themeDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 从默认主题复制基础模板和资源
	defaultThemeDir := filepath.Join(m.projectDir, "themes", "default")
	if err := copyThemeFiles(defaultThemeDir, themeDir); err != nil {
		return fmt.Errorf("failed to copy theme files: %w", err)
	}

	// 创建主题配置文件
	if err := m.createThemeConfig(name, themeDir); err != nil {
		return err
	}

	return nil
}

// Use 切换使用的主题
func (m *Manager) Use(name string) error {
	// 检查主题是否存在
	themeDir := filepath.Join(m.themesDir, name)
	if _, err := os.Stat(themeDir); os.IsNotExist(err) {
		return fmt.Errorf("theme %s not found", name)
	}

	// 读取配置文件
	configPath := filepath.Join(m.projectDir, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析当前配置
	var config map[string]interface{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 更新主题设置
	config["theme"] = name

	// 重新序列化配置，保持原有格式和注释
	updatedConfig, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 写回配置文件
	if err := os.WriteFile(configPath, updatedConfig, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// loadTheme 加载主题配置
func (m *Manager) loadTheme(name string) (*Theme, error) {
	configPath := filepath.Join(m.themesDir, name, "theme.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme config: %w", err)
	}

	var theme Theme
	if err := yaml.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme config: %w", err)
	}

	return &theme, nil
}

// createThemeConfig 创建主题配置文件
func (m *Manager) createThemeConfig(name string, themeDir string) error {
	theme := Theme{
		Name:        name,
		Version:     "0.1.0",
		Author:      "Your Name",
		Description: "A new Stars theme",
		License:     "MIT",
	}

	data, err := yaml.Marshal(theme)
	if err != nil {
		return fmt.Errorf("failed to marshal theme config: %w", err)
	}

	configPath := filepath.Join(themeDir, "theme.yaml")
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme config: %w", err)
	}

	return nil
}

// copyThemeFiles 从默认主题复制文件
func copyThemeFiles(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算目标路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 复制文件
		return copyFile(path, dstPath)
	})
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// Install 从 Git 仓库安装主题
func (m *Manager) Install(repo string) error {
	// 提取主题名称
	name := filepath.Base(repo)
	name = strings.TrimSuffix(name, ".git")

	// 检查主题是否已存在
	themeDir := filepath.Join(m.themesDir, name)
	if _, err := os.Stat(themeDir); !os.IsNotExist(err) {
		return fmt.Errorf("theme %s already exists", name)
	}

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "stars-theme-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// 克隆仓库
	cmd := exec.Command("git", "clone", repo, tmpDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// 验证主题结构
	if err := m.validateTheme(tmpDir); err != nil {
		return fmt.Errorf("invalid theme structure: %w", err)
	}

	// 移动到主题目录
	if err := os.Rename(tmpDir, themeDir); err != nil {
		return fmt.Errorf("failed to install theme: %w", err)
	}

	return nil
}

// validateTheme 验证主题结构
func (m *Manager) validateTheme(dir string) error {
	required := []string{
		"theme.yaml",
		"layouts",
		"layouts/_default",
		"layouts/_default/single.html",
		"layouts/_default/list.html",
		"static",
	}

	for _, path := range required {
		fullPath := filepath.Join(dir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("missing required file/directory: %s", path)
		}
	}

	return nil
}

// Update 更新主题
func (m *Manager) Update(name string) error {
	// TODO: 实现主题更新
	return nil
}

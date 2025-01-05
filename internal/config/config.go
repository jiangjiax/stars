package config

import (
	"embed"
	"fmt"
	"gopkg.in/yaml.v3"
	"math/big"
	"os"
	"regexp"
	"strings"
	"text/template"
)

// NFTConfig 表示 NFT 相关配置
type NFTConfig struct {
	Price         string `yaml:"price"`         // NFT 铸造价格(ETH)
	MaxSupply     int    `yaml:"maxSupply"`     // 最大铸造数量
	RoyaltyFee    int    `yaml:"royaltyFee"`    // 版税比例(0-5000)
	OnePerAddress bool   `yaml:"onePerAddress"` // 每个地址限铸一个
	Version       string `yaml:"version"`       // NFT 版本号
	ChainId       int    `yaml:"chainId"`       // 链 ID
}

// Verification 表示内容验证信息
type Verification struct {
	ArweaveId   string     `yaml:"arweaveId"`   // Arweave 交易 ID
	NftContract string     `yaml:"nftContract"` // NFT 合约地址
	Author      string     `yaml:"author"`      // 作者钱包地址
	ContentHash string     `yaml:"contentHash"` // 内容哈希值
	NFT         *NFTConfig `yaml:"nft"`         // NFT 配置
}

type Series struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Order       int    `yaml:"order"`
}

// NFT 参数限制
const (
	MinPrice      = "0.001" // 最小价格 0.001 ETH
	MaxPrice      = "10"    // 最大价格 10 ETH
	MinMaxSupply  = 1       // 最小供应量
	MaxMaxSupply  = 10000   // 最大供应量
	MinRoyaltyFee = 0       // 最小版税比例 0%
	MaxRoyaltyFee = 5000    // 最大版税比例 50%
)

// ValidateNFTConfig 验证 NFT 配置参数
func (c *NFTConfig) ValidateNFTConfig() error {
	// 验证价格
	price, ok := new(big.Float).SetString(c.Price)
	if !ok {
		return fmt.Errorf("invalid price format: %s", c.Price)
	}
	minPrice, _ := new(big.Float).SetString(MinPrice)
	maxPrice, _ := new(big.Float).SetString(MaxPrice)
	if price.Cmp(minPrice) < 0 || price.Cmp(maxPrice) > 0 {
		return fmt.Errorf("price must be between %s and %s ETH", MinPrice, MaxPrice)
	}

	// 验证最大供应量
	if c.MaxSupply < MinMaxSupply || c.MaxSupply > MaxMaxSupply {
		return fmt.Errorf("maxSupply must be between %d and %d", MinMaxSupply, MaxMaxSupply)
	}

	// 验证版税比例
	if c.RoyaltyFee < MinRoyaltyFee || c.RoyaltyFee > MaxRoyaltyFee {
		return fmt.Errorf("royaltyFee must be between %d and %d (0%% - 50%%)",
			MinRoyaltyFee, MaxRoyaltyFee)
	}

	// 验证版本号格式
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(c.Version) {
		return fmt.Errorf("invalid version format, must be semver (e.g. 1.0.0)")
	}

	return nil
}

// GetDefaultNFTConfig 获取默认的 NFT 配置
func GetDefaultNFTConfig() *NFTConfig {
	return &NFTConfig{
		Price:         "0.01",
		MaxSupply:     100,
		RoyaltyFee:    1000, // 10%
		OnePerAddress: true,
		Version:       "1.0.0",
	}
}

type Config struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	BaseURL     string `yaml:"baseURL"`
	Theme       string `yaml:"theme"`

	// 个人信息
	Author struct {
		Name          string `yaml:"name"`
		Title         string `yaml:"title"`
		Avatar        string `yaml:"avatar"`
		Bio           string `yaml:"bio"`
		Location      string `yaml:"location"`
		GitHub        string `yaml:"github"`
		Status        string `yaml:"status"`
		WalletAddress string `yaml:"walletAddress"`

		// 技术栈
		Skills []struct {
			Name       string `yaml:"name"`
			Level      string `yaml:"level"`
			Percentage int    `yaml:"percentage"`
		} `yaml:"skills"`

		// 项目展示
		Projects []struct {
			Name        string   `yaml:"name"`
			Description string   `yaml:"description"`
			URL         string   `yaml:"url"`     // GitHub 仓库地址
			Website     string   `yaml:"website"` // 项目网址
			Image       string   `yaml:"image"`
			Highlights  []string `yaml:"highlights"`
			Tags        []string `yaml:"tags"`
			Stars       int      `yaml:"stars"`
			Forks       int      `yaml:"forks"`
		} `yaml:"projects"`

		// 开源贡献
		Contributions []struct {
			Name          string   `yaml:"name"`
			Description   string   `yaml:"description"`
			URL           string   `yaml:"url"`
			Contributions []string `yaml:"contributions"`
			Impact        string   `yaml:"impact"`
		} `yaml:"contributions"`

		// 社交链接
		SocialLinks []struct {
			Name string `yaml:"name"`
			URL  string `yaml:"url"`
			Icon string `yaml:"icon"`
		} `yaml:"socialLinks"`

		// 教育背景
		Education []struct {
			School       string   `yaml:"school"`
			Degree       string   `yaml:"degree"`
			Major        string   `yaml:"major"`
			Year         string   `yaml:"year"`
			Achievements []string `yaml:"achievements"`
		} `yaml:"education"`

		// 工作经历
		Experience []struct {
			Company          string   `yaml:"company"`
			Position         string   `yaml:"position"`
			Period           string   `yaml:"period"`
			Responsibilities []string `yaml:"responsibilities"`
			Technologies     []string `yaml:"technologies"`
		} `yaml:"experience"`

		// 证书和资质
		Certifications []struct {
			Name   string `yaml:"name"`
			Issuer string `yaml:"issuer"`
			Date   string `yaml:"date"`
			Icon   string `yaml:"icon"`
		} `yaml:"certifications"`

		// 联系方式
		Contact struct {
			Email    string `yaml:"email"`
			WeChat   string `yaml:"wechat"` // 微信号
			Phone    string `yaml:"phone"`
			Telegram string `yaml:"telegram"` // Telegram 用户名
			Twitter  string `yaml:"twitter"`  // Twitter 用户名
		} `yaml:"contact"`

		// 推荐的文章系列
		RecommendedSeries []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Icon        string `yaml:"icon"`
		} `yaml:"recommendedSeries"`
	} `yaml:"author"`

	// 博客设置
	Blog struct {
		// 分页设置
		Pagination struct {
			PostsPerPage int `yaml:"postsPerPage"`
			ShowPages    int `yaml:"showPages"`
		} `yaml:"pagination"`
	} `yaml:"blog"`

	Series []Series `yaml:"series"`

	Newsletter Newsletter `yaml:"newsletter"`

	SEO SEO `yaml:"seo"`
}

// Newsletter 表示邮件订阅配置
type Newsletter struct {
	Enabled     bool   `yaml:"enabled"`     // 是否启用邮件订阅
	Provider    string `yaml:"provider"`    // 邮件服务提供商
	Description string `yaml:"description"` // 订阅描述文本
	Buttondown  struct {
		Username string `yaml:"username"` // Buttondown 用户名
		// APIKey   string `yaml:"apiKey"`  // Buttondown API key
	} `yaml:"buttondown"`
}

// SEO 配置
type SEO struct {
	Keywords []string `yaml:"keywords"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// LoadFromTemplate 从模板创建配置
func LoadFromTemplate(fs embed.FS, templatePath string, data interface{}) (*Config, error) {
	tmpl, err := template.ParseFS(fs, templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute config template: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal([]byte(buf.String()), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

默认主题包含了 Stars 博客引擎的基础主题功能和 Web3 集成。

## 特性

### 🎨 设计

- 响应式设计,完美适配移动端
- 暗色主题
- 星空动态背景
- 现代化动效
- 清晰的排版

### 🔧 功能

- 文章系列支持
- 标签云
- 目录导航
- RSS 订阅
- 邮件订阅
- 社交分享
- 二维码生成

### ⛓ Web3 集成

- NFT 铸造
- 钱包连接
- 区块链验证

## 安装

```bash
stars theme install github.com/jiangjiax/stars-theme-default
```

## 开发

### 开发环境

要求:
- Node.js 18+
- npm 9+

安装依赖:
```bash
npm install
```

开发:
```bash
# 监听 CSS 变更
npm run watch:css

# 监听 JS 变更
npm run watch:js
```

构建:
```bash
npm run build:css
npm run build:js
```

## 自定义

### 颜色变量

在 `static/css/theme.css` 中修改:

```css
:root {
    --stars-primary: #0a192f;
    --stars-secondary: #112240;
    --stars-accent: #64ffda;
    --stars-gold: #ffd700;
    --stars-text: #ccd6f6;
    --stars-muted: #8892b0;
}
```

### 布局

修改 `layouts/_default/baseof.html` 调整整体布局。

### 组件

在 `layouts/components/` 中修改或添加组件。

## 贡献

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交改动 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 作者

[@jiangjiax](https://github.com/jiangjiax)

## 致谢

感谢所有贡献者的支持! 
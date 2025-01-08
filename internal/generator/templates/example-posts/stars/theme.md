---
date: 2025-01-07
description: 详细介绍 Stars 主题系统的使用方法和自定义主题开发指南
series: Stars 教程
seriesOrder: 4
slug: theme-guide
tags:
  - 教程
  - 主题
  - 定制化
title: Stars 主题
verification:
    arweaveId: a_YLv8-4tzPWY0jn7ZcYpoAtp2CI5Biux0p2vIgV9pg
    nftContract: 0x760410d585110e149233919357E7C866bb51A841
    author: 0x16572b97410200e79AB6c9423F8d9778F0Fb9C54
    contentHash: 
    nft:
        price: "0"
        maxSupply: 9999
        royaltyFee: 0
        onePerAddress: true
        version: "1.0.0"
        chainId: 11155111
---

Stars 提供了灵活的主题系统，让你可以轻松定制网站的外观。本文将介绍如何使用和开发 Stars 主题。

## 主题基础

### 目录结构

一个典型的 Stars 主题目录结构如下：

```
themes/default/
├── layouts/         # 布局模板
│   ├── _default/    # 默认模板
│   │   ├── baseof.html
│   │   ├── list.html
│   │   └── single.html
│   │   └── tags.html
│   └── components/  # 可复用组件
│   └── index.html   # 首页模板
├── static/          # 静态资源
│   ├── css/        # 样式文件
│   ├── js/         # 脚本文件
│   └── images/     # 图片资源
│   └── abi/        # abi 文件
└── theme.yaml       # 主题配置
├── package.json     # 依赖管理
├── postcss.config.js # PostCSS 配置
├── tailwind.config.js # Tailwind 配置
└── webpack.config.js  # Webpack 配置
```

### 切换主题

有两种方式可以切换主题：

1. **通过命令行**：

```bash
# 列出已安装的主题
stars theme list

# 切换主题
stars theme use <theme-name>

# 安装新主题
stars theme install <theme-repo>

# 创建新主题
stars theme new <theme-name>
```

2. **直接修改配置文件**：

在项目根目录的 `config.yaml` 中修改 theme 字段：

```yaml
# config.yaml
title: "My Blog"
description: "A blog powered by Stars"
baseURL: "http://localhost:1313"
theme: "default"  # 将 "default" 改为你想使用的主题名称
```

> 注意：theme 的值必须与 themes 目录下的主题文件夹名称相对应。

## 模板系统

Stars 使用 Go 的模板引擎，支持强大的模板功能。

### 基础模板

`layouts/_default/baseof.html` 是最基础的模板：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }} - {{ .Site.Title }}</title>
    {{ partial "head.html" . }}
</head>
<body>
    {{ partial "header.html" . }}
    <main>
        {{ block "main" . }}{{ end }}
    </main>
    {{ partial "footer.html" . }}
</body>
</html>
```

### 页面类型

Stars 支持以下页面类型：

1. **首页** (`index.html`)
   - 网站首页模板
   - 可以展示最新文章、特色内容等

2. **列表页** (`_default/list.html`)
   - 文章列表页面
   - 支持分类、标签筛选

3. **文章页** (`_default/single.html`)
   - 单篇文章页面
   - 显示文章内容和元数据

4. **标签页** (`tags/list.html`)
   - 标签云页面
   - 展示所有标签和统计

## 资源处理

### CSS 样式

Stars 使用 Tailwind CSS 作为默认样式框架。

### JavaScript

支持添加自定义脚本。

## 开发自定义主题

### 1. 创建主题

```bash
stars theme new my-theme
```

这将创建一个基本的主题结构。

### 2. 主题配置

在 `theme.yaml` 中定义主题信息：

```yaml
name: "My Theme"
version: "1.0.0"
author: "Your Name"
description: "A beautiful Stars theme"

# 主题选项
options:
  colorScheme: "light"
  showAuthor: true
  enableComments: false
```

### 3. 模板变量

可用的主要变量：

- `.Site`: 网站配置信息
- `.Title`: 页面标题
- `.Content`: 页面内容
- `.Posts`: 文章列表
- `.Series`: 系列信息
- `.Tags`: 标签信息

### 4. Web3 集成

主题可以集成 Web3 功能：

```html
<!-- NFT 铸造按钮 -->
<button class="mint-nft" 
        data-contract="{{ .Verification.NftContract }}"
        data-price="{{ .Verification.Nft.Price }}">
    铸造 NFT
</button>

<!-- 内容验证信息 -->
<div class="verification-info">
    <p>Content Hash: {{ .Verification.ContentHash }}</p>
    <p>Arweave ID: {{ .Verification.ArweaveId }}</p>
</div>
```

## 最佳实践

1. **性能优化**
   - 优化图片资源
   - 合理使用缓存

2. **可维护性**
   - 添加必要的注释
   - 使用版本控制

3. **用户体验**
   - 支持响应式设计
   - 优化加载性能
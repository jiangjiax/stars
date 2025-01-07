---
date: 2025-01-07
description: 介绍 Stars 的基本使用方法
series: Stars 教程
seriesOrder: 1
slug: welcome-to-stars
tags:
  - 入门
  - 教程
title: Stars 入门
verification:
    arweaveId: xyz789abc123def456
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

Stars 是一个基于 Go 语言开发的去中心化个人网站生成器，支持 Web3 功能。本文将指导你从零开始使用 Stars 创建你的个人网站。

## 环境准备

### 1. 安装 Go

Stars 需要 Go 1.21 或更高版本。

**macOS**:
```bash
# 使用 Homebrew 安装
brew install go

# 验证安装
go version
```

**Linux**:
```bash
# 下载并解压 Go
wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

**Windows**:
1. 访问 [Go 下载页面](https://go.dev/dl/)
2. 下载 Windows 安装包
3. 运行安装程序
4. 打开命令提示符验证: `go version`

### 2. 配置 Go 环境

```bash
# 设置 GOPROXY（国内用户推荐）
go env -w GOPROXY=https://goproxy.cn,direct

# 启用 Go modules
go env -w GO111MODULE=on
```

### 3. 安装 Node.js

Stars 的主题系统需要 Node.js 18+ 和 npm。

**使用 nvm 安装（推荐）**:
```bash
# 安装 nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 安装 Node.js
nvm install 18
nvm use 18

# 验证安装
node --version
npm --version
```

## 安装 Stars

```bash
# 安装最新版本
go install github.com/jiangjiax/stars/cmd/stars@latest

# 验证安装
stars -v
```

## 创建新网站

### 1. 创建项目

```bash
# 创建新的博客项目
stars new my-blog

# 进入项目目录
cd my-blog
```

### 2. 目录结构

```
my-blog/
├── content/          # 内容目录
│   └── posts/        # 文章存放处，可以嵌套文件夹
├── themes/           # 主题目录
│   └── default/      # 默认主题
├── public/           # 生成的静态文件
└── config.yaml       # 配置文件
```

### 3. 创建文章

Stars 提供了便捷的命令行工具来创建文章：

基本用法：

```bash
stars post "文章标题"
```

指定系列、标签、描述和自定义url等：

```bash
stars post "文章标题" --series "系列名称" --tags "功能,更新" --slug "about" -desc "这是一篇关于技术的分享文章"
```

也可以直接在 `content/posts` 目录下创建 Markdown 文件。

### 4. 本地预览

```bash
# 启动开发服务器
stars server -p 1313

# 访问 http://localhost:1313
```

### 5. 构建网站

```bash
# 生成静态文件
stars build

# 生成的文件在 public 目录
```

## Web3 功能

具体如何使用 web3 功能，请参考[这里](./content-creation)。

1. **NFT 铸造**
- 将文章铸造为 NFT
- 支持多链部署
- 版税分成机制

2. **内容验证**
- 链上内容验证
- Arweave 永久存储
- 数字签名支持

## 获取帮助

- GitHub Issues: [报告问题](https://github.com/jiangjiax/stars/issues)
- 社区：[Discussions](https://github.com/jiangjiax/stars/discussions)
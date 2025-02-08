# Stars(繁星) - Web3 个人网站生成器

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Release](https://img.shields.io/github/v/release/jiangjiax/stars)](https://github.com/jiangjiax/stars/releases)

[English](./README_EN.md) | 简体中文

## 项目概述

Stars(繁星)是一个基于 Go 语言开发的去中心化个人网站生成器，支持 Web3 功能。它能帮助你快速创建一个现代化的个人网站，并提供 Web3 集成功能。本项目旨在为创作者提供一个去中心化的内容发布平台。

### 智能合约

Stars 项目支持多链部署的 NFT 合约系统，用户可以在文章的元数据中填写 Stars 官方部署的 NFT 合约地址，也可以自行部署。

- **Ethereum Sepolia**
  - Address: `0x5c83f2287833F567b1D80D7B981084eb5CaeF445`
  - [Contract Code](https://github.com/jiangjiax/stars-contracts)
  - [Etherscan](https://sepolia.etherscan.io/address/0x5c83f2287833F567b1D80D7B981084eb5CaeF445)
  - Chain ID: 11155111

- **Telos Testnet**
  - Address: `0x903e48Ca585dBF4dFeb74f2864501feB6f0dF369`
  - [Contract Code](https://github.com/jiangjiax/stars-contracts)
  - [TelosScan](https://testnet.teloscan.io/address/0x903e48Ca585dBF4dFeb74f2864501feB6f0dF369)
  - Chain ID: 41

- **EDU Chain Testnet**
  - Address: `0xcA3Dbe8eF976e606B8c96052aaC22763aDeAEE0A`
  - [Contract Code](https://github.com/jiangjiax/stars-contracts)
  - [TelosScan](https://edu-chain-testnet.blockscout.com/address/0xcA3Dbe8eF976e606B8c96052aaC22763aDeAEE0A)
  - Chain ID: 656476

### 主要特性

🚀 **高性能静态站点生成器**

🎨 **现代化主题系统**

📱 **移动适配**

🔗 **Web3 功能集成**

🛠 **开发者友好**

### 贡献指南

我们欢迎所有形式的贡献，无论是新功能、文档改进还是问题报告。请查看我们的[贡献指南](./CONTRIBUTING.md)了解更多信息。

### 贡献者

感谢以下贡献者的支持：

<a href="https://github.com/jiangjiax/stars/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=jiangjiax/stars" />
</a>

## 快速开始

[Stars 入门](./internal/generator/templates/example-posts/stars/welcome.md)

## 常见问题

### Q: Stars 是什么？
A: Stars(繁星)是一个去中心化的个人网站生成器，让创作者能够以去中心化的方式发布和管理内容，同时支持 Web3 功能集成。

### Q: 文章内容如何存储在区块链上？
A: 与 NFT 艺术品类似，Stars 将文章内容存储在去中心化存储系统中，链上只存储内容的去中心化地址。无论是图片、文章还是其他形式的数字作品，都采用相同的存储机制。

### Q: Stars 是否会对内容进行审查？
A: Stars 是完全去中心化的平台，内容的发布决定权完全在创作者手中。平台本身不会也无法对内容进行审查。

### Q: 已发布的内容可以修改吗？
A: 由于区块链的不可篡改特性，已上链的内容无法直接修改。但 Stars 支持版本管理机制，创作者可以发布新版本的内容（如 v1.0.0 升级到 v1.0.1）。创作者可以自由选择展示或隐藏特定版本的内容。

### Q: Stars 支持哪些区块链网络？
A: Stars 目前支持多个区块链网络，包括 Ethereum Sepolia、Telos Testnet 和 EDU Chain Testnet。用户可以选择使用官方部署的智能合约，也可以自行部署到支持的网络上。

### Q: 如何保证内容的永久性？
A: Stars 推荐用户使用去中心化存储系统来存储内容，并在文章的元数据上填写去中心化存储文件的地址，确保内容的持久性和可访问性。即使用户的 Stars 个人网站不再运营，内容仍然可以通过去中心化网络访问。

### Q: Stars 如何保护创作者的权益？
A: 通过区块链技术，每篇文章都会生成唯一的链上记录，为创作者提供了内容所有权的证明。创作者可以通过 NFT 的形式确权和变现自己的作品。

### Q: 使用 Stars 需要付费吗？
A: Stars 本身是开源免费的。但在发布内容到区块链时，需要支付相应网络的 Gas 费用。这些费用由所选择的区块链网络决定，与 Stars 无关。
---
date: 2025-01-07
description: 详细介绍 Stars 文章的元数据配置和使用方法
series: Stars 教程
seriesOrder: 2
slug: metadata-guide
tags:
  - 教程
  - Markdown
  - Front Matter
title: Stars 文章的元数据介绍
verification:
    arweaveId: xyz789abc123def456
    nftContract: 0x5c83f2287833F567b1D80D7B981084eb5CaeF445
    author: 0x16572b97410200e79AB6c9423F8d9778F0Fb9C54
    contentHash: bGP4K5KQZJvKe_JOqrD9u99WS8YDHlIuDjjfCuaGtUk
    nft:
        price: "0"
        maxSupply: 9999
        royaltyFee: 0
        onePerAddress: true
        version: "1.0.0"
        chainId: 11155111
---

Stars 在头部定义文章的元数据。本文将详细介绍每个元数据字段的含义和使用方法。

## 基础元数据

### 必填字段

必填的元数据字段包括：

- `title`: 文章的标题
- `date`: 发布日期，格式为 YYYY-MM-DD
- `description`: 文章简短描述，用于预览和SEO

### 可选字段

- `slug`: 自定义 URL，默认使用文件名
- `draft`: 是否为草稿，草稿不会被发布

## 分类和组织

### 标签系统

使用 `tags` 字段为文章添加标签，推荐使用列表格式：

```yaml
tags:
  - Web3
  - 教程
  - 入门
```

### 文章系列

使用 `series` 和 `seriesOrder` 组织相关文章：

- `series`: 系列名称，必须和config.yaml中的series.name一致
- `seriesOrder`: 在系列中的顺序，从1开始，如果seriesOrder为0，则不会显示在系列导航中

## Web3 相关配置

### 验证信息

使用 `verification` 字段配置区块链相关信息：

- `arweaveId`: Arweave 存储 ID，你也可以放置别的存储方ID，但我建议使用去中心化的基础设施，比如[ArDrive](https://ardrive.io/)
- `nftContract`: NFT 合约地址，你可以在[这里](https://github.com/jiangjiax/stars/blob/main/CONTRACTS.md)查看当前我们部署的智能合约地址，我们部署了多个区块链网路的地址，你也可以修改合约代码并部署得到你自己的合约地址
- `author`: 作者钱包地址
- `contentHash`: 文章内容哈希

#### 文章内容哈希的生成逻辑

1. 使用 Keccak-256 哈希算法（以太坊使用的哈希算法）
2. 将文章的标题和内容组合成一个 JSON 对象，将 JSON 对象序列化为字符串，对序列化后的字符串进行哈希计算
3. 将哈希结果转换为十六进制字符串，在结果前加上 "0x" 前缀，在结果后加上版本号，文章的内容哈希就生成了
4. 在 ContentChanged 方法中用于检查内容是否发生变化
5. 在 UpdateContentHash 方法中用于更新内容哈希
6. 在文章保存时会自动计算并更新哈希值
7. 文章的内容哈希和作者的地址将作为NFT的唯一标识

### NFT 配置

在 `verification.nft` 中配置 NFT 相关参数：

- `price`: NFT 铸造价格（ETH），0就是免费铸造
- `maxSupply`: NFT 最大供应量，超出后不可再铸造
- `royaltyFee`: 版税比例，范围 0-5000（0%-50%）
- `onePerAddress`: 是否每个地址限购一个，true 表示每个地址限购一个，false 表示不限制
- `version`: NFT 版本号，如 "1.0.0"
- `chainId`: 链 ID（如：1=以太坊主网，11155111=Sepolia测试网）你可以在[这里](https://github.com/jiangjiax/stars/blob/main/CONTRACTS.md)查看当前支持的智能合约的chainId

## 示例

一个完整的文章配置示例：

```yaml
---
title: 如何使用 Stars 创建 Web3 博客
date: 2025-01-07
description: 详细介绍如何使用 Stars 搭建支持 Web3 功能的个人博客
slug: how-to-create-web3-blog
tags:
  - Web3
  - 教程
  - Stars
series: Stars 教程
seriesOrder: 3
draft: false
verification:
  arweaveId: xyz789abc123def456
  nftContract: 0x5c83f2287833F567b1D80D7B981084eb5CaeF445
  author: 0x16572b97410200e79AB6c9423F8d9778F0Fb9C54
  contentHash: 0x987654321fedcba
  nft:
    price: "0.001"
    maxSupply: 100
    royaltyFee: 1000
    onePerAddress: true
    version: "1.0.0"
    chainId: 11155111
---
```

## 常见问题

1. **为什么我的文章没有显示？**
   - 检查 `draft` 是否设置为 `true`
   - 确认日期是否正确

2. **如何修改已发布文章的 URL？**
   - 更新 `slug` 字段

3. **NFT 相关字段是必填的吗？**
   - 不是必填的，只有在需要 NFT 功能时才需要配置
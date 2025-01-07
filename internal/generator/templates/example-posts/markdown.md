---
date: 2025-01-07
description: 详细介绍 Stars 支持的 Markdown 语法和使用方法
series: Stars 教程
seriesOrder: 3
slug: markdown-guide
tags:
  - 教程
  - Markdown
  - 内容创作
title: Stars Markdown 语法
verification:
    arweaveId: 
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

Stars 支持标准的 Markdown 语法，并添加了一些扩展功能。本文将详细介绍如何使用这些语法来创建丰富的内容。

## 基础语法

### 标题

使用 `#` 创建标题，支持 1-6 级标题：

# 一级标题
## 二级标题
### 三级标题
#### 四级标题
##### 五级标题
###### 六级标题

### 文本格式

- **粗体文本** (`**粗体文本**`)
- *斜体文本* (`*斜体文本*`)
- ***粗斜体*** (`***粗斜体***`)
- ~~删除线~~ (`~~删除线~~`)
- `行内代码` (`` `行内代码` ``)
- ==高亮文本== (`==高亮文本==`)

### 列表

无序列表：
- 项目 1
- 项目 2
  - 子项目 2.1
  - 子项目 2.2
- 项目 3

有序列表：
1. 第一步
2. 第二步
   1. 子步骤 2.1
   2. 子步骤 2.2
3. 第三步

任务列表：
- [x] 已完成任务
- [ ] 未完成任务
- [ ] 待办事项

### 引用

> 这是一段引用文本
> 
> 可以包含多个段落
>> 也可以嵌套引用

## 扩展功能

### 代码块

支持多种编程语言的语法高亮：

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Stars!")
}
```

```javascript
function greet(name) {
    console.log(`Hello, ${name}!`);
}
```

### 表格

| 功能 | 语法 | 说明 |
|------|------|------|
| 标题 | `#` | 1-6 级标题 |
| 粗体 | `**text**` | 加粗文本 |
| 斜体 | `*text*` | 倾斜文本 |
| 代码块 | ` ```language ` | 支持语法高亮 |

### 脚注

这是一个带有脚注的文本[^1]。

[^1]: 这是脚注的内容。

### 链接和图片

[Stars 项目主页](https://github.com/jiangjiax/stars)

图片：
![Stars Logo](https://example.com/stars-logo.png)
*图片说明文字*

### Emoji 表情

Stars 支持 Emoji 表情符号：

:smile: :rocket: :star: :heart:

## 最佳实践

1. **文档结构**
   - 使用合适的标题层级
   - 保持内容层次清晰
   - 适当使用空行分隔段落

2. **格式规范**
   - 列表项保持一致的格式
   - 代码块指定正确的语言
   - 表格对齐美观

3. **图片处理**
   - 使用合适大小的图片
   - 添加有意义的替代文本
   - 适当添加图片说明

4. **链接使用**
   - 确保链接可访问
   - 使用描述性的链接文本
   - 考虑是否需要外部链接

## 相关资源

- [Markdown 官方指南](https://markdown.com.cn)

---

*文末可以添加作者信息或其他说明*
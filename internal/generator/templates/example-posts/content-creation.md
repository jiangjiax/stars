---
title: "Stars 内容创作指南"
date: 2024-03-21
description: "学习如何使用 Stars 创作和管理你的内容"
tags: ["教程", "Markdown", "内容创作"]
slug: "content-creation-guide"
series: "Stars 教程"
seriesOrder: 2
---

本文将介绍如何使用 Stars 创建和组织你的内容，并展示所有支持的 Markdown 语法。

## 基础语法

### 标题

使用 `#` 创建标题，支持 1-6 级标题：

# 一级标题
## 二级标题
### 三级标题
#### 四级标题
##### 五级标题
###### 六级标题

### 文本样式

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
[x] 已完成任务
[ ] 未完成任务
[ ] 待办事项

### 引用

> 这是一段引用文本
> 
> 可以���含多个段落
>> 也可以嵌套引用

## 高级功能

### 代码块

带语法高亮的代码块：

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

| 功能     | 语法                | 说明         |
|---------|-------------------|------------|
| 标题     | `#`              | 1-6 级标题   |
| 粗体     | `**text**`       | 加粗文本     |
| 斜体     | `*text*`         | 倾斜文本     |
| 代码块    | ` ```language ` | 支持语法高亮   |

### 脚注

这是一个带有脚注的文本[^1]。

[^1]: 这是脚注的内容。

### 链接和图片

[Stars 项目主页](https://github.com/yourusername/stars)

图片：
![Stars Logo](https://informedainews.com/assets/images/ai1-b50e820e78cfc0a3ee336925f65a5161.jpeg)
*图片说明文字*

### Emoji 表情

:smile: :rocket: :star: :heart:

### 水平分割线

---

## Front Matter

每篇文章开头的 Front Matter 用于定义文章的元数据：

```yaml
---
title: "文章标题"
date: 2024-03-21
description: "文章描述"
tags: ["标签1", "标签2"]
slug: "url-slug"
series: "系列名称"
seriesOrder: 1
---
```

## 最佳实践

1. 使用合适的标题层级
2. 添加有意义的描述和标签
3. 适当使用图片和代码示例
4. 保持文章结构清晰
5. 使用系列功能组织相关文章

## 扩展阅读

- [Markdown 官方指南](https://markdown.com.cn)
- [Stars 文档](https://stars-docs.com)

---

*文末可以添加作者信息或其他说明*
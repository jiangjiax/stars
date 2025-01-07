# 贡献指南

感谢你对 Stars 项目感兴趣！我们欢迎任何形式的贡献。

## 如何贡献

### 报告 Bug

1. 确保 bug 没有在 [issues](https://github.com/jiangjiax/stars/issues) 中被报告过
2. 打开一个新的 issue，使用 Bug 报告模板
3. 清晰地描述问题，包括:
   - 复现步骤
   - 预期行为
   - 实际行为
   - 环境信息

### 提交新功能

1. 先在 issues 中讨论新功能
2. 获得维护者同意后再开始开发
3. 遵循项目的代码规范
4. 提交 PR 时附上完整的功能描述

### Pull Request 流程

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/**`)
3. 提交改动 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/**`)
5. 打开 Pull Request，并附上完整的功能描述
6. 等待维护者审核

## 开发设置

### 环境要求

- Go 1.21+
- Node.js 18+
- npm 9+

### 本地开发 

## 代码规范

### Go 代码规范

- 使用 `gofmt` 格式化代码
- 使用 `gocyclo` 检查代码复杂度
- 遵循 [Effective Go](https://golang.org/doc/effective_go) 指南
- 添加必要的注释和文档
- 编写单元测试

### 提交信息规范

使用语义化的提交信息:

- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

## 版本发布

我们使用 [语义化版本](https://semver.org/lang/zh-CN/) 进行版本控制:

- 主版本号: 不兼容的 API 修改
- 次版本号: 向下兼容的功能性新增
- 修订号: 向下兼容的问题修正

## 行为准则

### 我们的承诺

为了建设开放和友好的环境，我们承诺：

- 包容不同的观点和经验
- 友善地接受建设性批评
- 以社区的最大利益为重
- 展现同理心

### 不可接受的行为

不可接受的行为包括：

- 使用性别、种族、宗教等方面的歧视性语言
- 公开或私下的骚扰
- 未经他人明确许可而发布他人的私人信息
- 其他不道德或不专业的行为

## 版本管理策略

### 分支策略

- main      - 稳定的发布分支
- dev       - 开发分支，新特性集成测试
- feature/* - 新功能开发 (例如: feature/ipfs-storage)
- fix/*     - Bug 修复 (例如: fix/toc-scroll)
- docs/*    - 文档更新 (例如: docs/api-reference)
- release/* - 版本发布准备 (例如: release/v0.1.0)

### 版本号管理

采用语义化版本 (Semantic Versioning): v(major).(minor).(patch) [例如: v0.1.0, v1.0.0, v1.2.3]

- major: 不兼容的 API 修改
- minor: 向后兼容的功能新增
- patch: 向后兼容的问题修复

### 发布流程

dev -> release/v*.*.* -> master

1. 在 dev 分支开发并测试
2. 创建 release 分支准备发布
3. 完成测试后合并到 master
4. 在 master 打 tag 发布

```bash
# 开发新功能
git checkout dev
git checkout -b feature/new-feature
# 开发完成后
git checkout dev
git merge feature/new-feature

# 准备发布
git checkout -b release/v1.0.0
# 测试和修复
git checkout master
git merge release/v1.0.0
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

## 许可证

通过贡献代码，你同意你的贡献将按照项目的 MIT 许可证进行授权。

## 联系我们

邮箱: jiangjiaxingogogo@gmail.com
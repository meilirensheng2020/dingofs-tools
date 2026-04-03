# CLAUDE.md

此文件用于配置 Claude Code 在此项目中的行为。

## 自动提交配置

Claude Code 会自动将生成的代码变更提交到仓库，提交人信息：
- Name: Claude Code
- Email: claude-code@anthropic.com

所有提交会自动添加 `[claude-code]` 前缀以便识别。

## 使用方式

直接向 Claude Code 发出指令，它会自动：
1. 生成/修改代码
2. 创建 git commit
3. 提交到当前分支

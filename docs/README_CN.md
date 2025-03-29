# Perplexity Search MCP 服务器

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

基于 Go 语言实现的 Perplexity Search MCP 服务器，允许大型语言模型 (LLM) 通过 [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) 访问 Perplexity 搜索 API。

## 功能特性

- **perplexity_search**: 执行网络搜索并返回结果，包括引用链接
  - **参数**:
    - `query` (字符串, 必填): 搜索查询内容
    - `search_recency_filter` (字符串, 可选): 按时间筛选结果 (`month`, `week`, `day`, `hour`)
    - `max_tokens` (整数, 可选): 返回的最大令牌数
    - `temperature` (数值, 可选, 默认值: 0.2): 控制响应的随机性
    - `top_p` (数值, 可选, 默认值: 0.9): 核采样阈值
    - `search_domain_filter` (数组, 可选): 限制搜索结果的域名列表
    - `return_images` (布尔值, 可选): 在结果中包含图片链接
    - `return_related_questions` (布尔值, 可选): 包含相关问题
    - `top_k` (数值, 可选, 默认值: 0): 用于 top-k 过滤的令牌数
    - `stream` (布尔值, 可选): 增量流式传输响应
    - `presence_penalty` (数值, 可选, 默认值: 0): 调整讨论新主题的可能性
    - `frequency_penalty` (数值, 可选, 默认值: 1): 减少重复内容
    - `web_search_options` (对象, 可选): Web 搜索的配置选项

## 设置与使用

### 前提条件

- Go 1.19 或更高版本
- Perplexity API 密钥

### 安装

1. 克隆仓库:

```bash
git clone https://github.com/chenxilol/perplexity-mcp-go.git
cd perplexity-mcp-go
```

2. 构建应用:

```bash
go build -o perplexity-search-mcp
```

### 本地运行

1. 设置 Perplexity API 密钥:

```bash
export PERPLEXITY_API_KEY="你的API密钥"
```

2. 运行服务器:

```bash
./perplexity-search-mcp
```

### 与 Claude 集成

1. 将提供的 `claude_desktop_config.json` 复制到 Claude 配置目录:
   - Windows: `%USERPROFILE%\AppData\Roaming\Claude\claude_desktop_config.json`
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Linux: `~/.config/Claude/claude_desktop_config.json`

2. 编辑配置文件以包含您的 API 密钥:

```json
{
  "mcpServers": {
    "perplexity-search": {
      "command": "/path/to/perplexity-search-mcp",
      "env": {
        "PERPLEXITY_API_KEY": "你的API密钥"
      }
    }
  }
}
```

### Docker 支持

1. 构建 Docker 镜像:

```bash
docker build -t perplexity-search-mcp:latest .
```

2. 运行容器:

```bash
docker run -i --rm -e PERPLEXITY_API_KEY=你的API密钥 perplexity-search-mcp:latest
```

## 使用示例

配置完成后，Claude 可以通过 MCP 协议使用 perplexity_search 工具执行实时网络搜索。

带参数的搜索示例:
```json
{
  "query": "最新的人工智能研究进展",
  "search_recency_filter": "week",
  "temperature": 0.5,
  "return_related_questions": true,
  "web_search_options": {
    "search_context_size": "high"
  }
}
```

## 故障排除

如果遇到问题:
1. 验证您的 API 密钥是否正确设置
2. 检查网络连接
3. 查看 stderr 日志了解错误信息

## 许可证

本项目采用 MIT 许可证 - 详情请查看 [LICENSE](LICENSE) 文件。

## 致谢

- [Model Context Protocol](https://modelcontextprotocol.io/) 提供的 MCP 规范
- [MCP-Go](https://github.com/mark3labs/mcp-go) 提供的 Go MCP 实现
- [Perplexity](https://www.perplexity.ai/) 提供的搜索 API 

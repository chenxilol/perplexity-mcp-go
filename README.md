# Perplexity Search MCP Server

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go implementation of a Perplexity Search MCP server that allows large language models (LLMs) to access the Perplexity search API through the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/).

## Features

- **perplexity_search**: Perform web searches and return results, including citations
  - **Parameters**:
    - `query` (string, required): The search query
    - `search_recency_filter` (string, optional): Filter results by time (`month`, `week`, `day`, `hour`)
    - `max_tokens` (integer, optional): Maximum number of tokens to return
    - `temperature` (number, optional, default: 0.2): Controls randomness in response
    - `top_p` (number, optional, default: 0.9): Nucleus sampling threshold
    - `search_domain_filter` (array, optional): List of domains to limit search results
    - `return_images` (boolean, optional): Include image links in results
    - `return_related_questions` (boolean, optional): Include related questions
    - `top_k` (number, optional, default: 0): Number of tokens for top-k filtering
    - `stream` (boolean, optional): Stream response incrementally
    - `presence_penalty` (number, optional, default: 0): Adjust likelihood of new topics
    - `frequency_penalty` (number, optional, default: 1): Reduce repetition
    - `web_search_options` (object, optional): Configuration options for web search

## Setup & Usage

### Prerequisites

- Go 1.23 or higher
- Perplexity API key

### Installation

1. Clone the repository:

```bash
git clone https://github.com/chenxilol/perplexity-mcp-go.git
cd perplexity-mcp-go
```

2. Build the application:

```bash
go build -o perplexity-search-mcp
```

### Running Locally

1. Set your Perplexity API key:

```bash
export PERPLEXITY_API_KEY="your-api-key-here"
```

2. Run the server:

```bash
./perplexity-search-mcp
```

### Integrating with Claude

1. Copy the provided `claude_desktop_config.json` to your Claude configuration directory:
   - Windows: `%USERPROFILE%\AppData\Roaming\Claude\claude_desktop_config.json`
   - macOS: `~/Library/Application Support/Claude/claude_desktop_config.json`
   - Linux: `~/.config/Claude/claude_desktop_config.json`

2. Edit the configuration file to include your API key:

```json
{
  "mcpServers": {
    "perplexity-search": {
      "command": "/path/to/perplexity-search-mcp",
      "env": {
        "PERPLEXITY_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

### Docker Support

1. Build the Docker image:

```bash
docker build -t perplexity-search-mcp:latest .
```

2. Run the container:

```bash
docker run -i --rm -e PERPLEXITY_API_KEY=your-api-key-here perplexity-search-mcp:latest
```

## Example Usage

Once configured, Claude can use the perplexity_search tool via MCP to perform real-time web searches.

Example search with parameters:
```json
{
  "query": "latest AI research developments",
  "search_recency_filter": "week",
  "temperature": 0.5,
  "return_related_questions": true,
  "web_search_options": {
    "search_context_size": "high"
  }
}
```

## Troubleshooting

If you encounter issues:
1. Verify your API key is correctly set
2. Check network connectivity
3. Examine stderr logs for error messages

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Model Context Protocol](https://modelcontextprotocol.io/) for the MCP specification
- [MCP-Go](https://github.com/mark3labs/mcp-go) for the Go MCP implementation
- [Perplexity](https://www.perplexity.ai/) for their search API 

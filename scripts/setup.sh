#!/bin/bash

# Check if API key is provided
if [ -z "$1" ]; then
    echo "Usage: ./setup.sh <your-perplexity-api-key>"
    echo "Please provide your Perplexity API key as an argument."
    exit 1
fi

API_KEY=$1

# Build application
echo "===> Building Perplexity Search MCP Server..."
go build -o perplexity-search-mcp

# Check if build was successful
if [ $? -ne 0 ]; then
    echo "Error: Build failed. Please make sure Go is installed correctly."
    exit 1
fi

# Create configuration file
echo "===> Creating Claude desktop configuration file..."

cat > claude_desktop_config.json << EOL
{
  "mcpServers": {
    "perplexity-search": {
      "command": "$(pwd)/perplexity-search-mcp",
      "env": {
        "PERPLEXITY_API_KEY": "${API_KEY}"
      }
    }
  }
}
EOL

# Create a run script with provided API key
echo "===> Creating run script..."

cat > run_with_key.sh << EOL
#!/bin/bash
export PERPLEXITY_API_KEY="${API_KEY}"
./perplexity-search-mcp
EOL

chmod +x run_with_key.sh

echo "===> Setup completed successfully!"
echo ""
echo "To run the server:"
echo "  ./run_with_key.sh"
echo ""
echo "To integrate with Claude, copy the claude_desktop_config.json to your Claude configuration directory:"
echo "  - Windows: %USERPROFILE%\\AppData\\Roaming\\Claude\\claude_desktop_config.json"
echo "  - macOS: ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "  - Linux: ~/.config/Claude/claude_desktop_config.json" 
#!/bin/bash

# Set API key
export PERPLEXITY_API_KEY="your-perplexity-api-key"

# Compile application
echo "===> Compiling application..."
go build -o perplexity-search-mcp

# Check if compilation was successful
if [ $? -ne 0 ]; then
  echo "Error: Application compilation failed"
  exit 1
fi

echo "===> Application compiled successfully"

# Create local configuration file
cat > claude_local_config.json << EOL
{
  "mcpServers": {
    "perplexity-search": {
      "command": "$(pwd)/perplexity-search-mcp",
      "env": {
        "PERPLEXITY_API_KEY": "${PERPLEXITY_API_KEY}"
      }
    }
  }
}
EOL

echo "===> Local configuration file created: claude_local_config.json"
echo "Please copy this configuration file to your Claude desktop configuration directory:"
echo "  - Windows: %USERPROFILE%\\AppData\\Roaming\\Claude\\claude_desktop_config.json"
echo "  - macOS: ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "  - Linux: ~/.config/Claude/claude_desktop_config.json"
echo ""
echo "===> Run the application:"
echo "$(pwd)/perplexity-search-mcp" 
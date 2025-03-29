#!/bin/bash

# Set variables
IMAGE_NAME="perplexity-search-mcp"
API_KEY="your-perplexity-api-key"

echo "===> Building Docker image..."
docker build --no-cache -t $IMAGE_NAME:latest .

# Check if build was successful
if [ $? -ne 0 ]; then
  echo "Error: Docker image build failed"
  exit 1
fi

echo "===> Image built successfully: $IMAGE_NAME:latest"

echo "===> You can run the container manually with:"
echo "docker run -i --rm -e PERPLEXITY_API_KEY=$API_KEY $IMAGE_NAME:latest"

echo ""
echo "===> Claude desktop configuration has been generated in claude_desktop_config.json"
echo "Please copy this configuration file to your Claude desktop configuration directory:"
echo "  - Windows: %USERPROFILE%\\AppData\\Roaming\\Claude\\claude_desktop_config.json"
echo "  - macOS: ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "  - Linux: ~/.config/Claude/claude_desktop_config.json" 
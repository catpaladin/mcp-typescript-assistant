package main

import (
	"context"
	"log"
	"os"

	"mcp-typescript-assistant/internal/server"
)

// main entry point for the TypeScript MCP server
func main() {
	// Create context
	ctx := context.Background()

	// Create and run the server
	mcpServer := server.NewTypeScriptMCPServer()
	
	if err := mcpServer.Run(ctx); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
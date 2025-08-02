package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mcp-typescript-assistant/internal/server"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		<-sigCh
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Create and start the MCP server
	mcpServer := server.NewTypeScriptMCPServer()
	
	log.Println("TypeScript MCP Server starting...")
	
	if err := mcpServer.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	
	log.Println("TypeScript MCP Server stopped")
}

func init() {
	// Check for required environment variables or configuration
	if os.Getenv("DEBUG") == "true" {
		log.Println("Debug mode enabled")
	}
	
	// Print usage information to stderr (since stdout is used for MCP communication)
	fmt.Fprintln(os.Stderr, "TypeScript MCP Server - Model Context Protocol server for TypeScript development tools")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Available tools:")
	fmt.Fprintln(os.Stderr, "  - type-check: Run TypeScript type checking")
	fmt.Fprintln(os.Stderr, "  - get-types: Extract type information")
	fmt.Fprintln(os.Stderr, "  - lint-check: Run ESLint checking")
	fmt.Fprintln(os.Stderr, "  - suggest-improvements: Suggest code improvements")
	fmt.Fprintln(os.Stderr, "  - load-guidelines: Load custom coding guidelines")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Prerequisites:")
	fmt.Fprintln(os.Stderr, "  - TypeScript: npm install -g typescript")
	fmt.Fprintln(os.Stderr, "  - ESLint: npm install -g eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin")
	fmt.Fprintln(os.Stderr, "")
}
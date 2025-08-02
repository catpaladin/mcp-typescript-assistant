package server

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TypeScriptMCPServer represents the main MCP server for TypeScript tools
type TypeScriptMCPServer struct {
	server   *mcp.Server
	handlers *Handlers
}

// NewTypeScriptMCPServer creates a new TypeScript MCP server
func NewTypeScriptMCPServer() *TypeScriptMCPServer {
	handlers := NewHandlers()
	
	server := mcp.NewServer("typescript-analyzer", "1.0.0", nil)

	mcpServer := &TypeScriptMCPServer{
		server:   server,
		handlers: handlers,
	}

	mcpServer.registerTools()
	return mcpServer
}

// registerTools registers all the TypeScript tools with the MCP server
func (s *TypeScriptMCPServer) registerTools() {
	// Create tools using NewServerTool
	typeCheckTool := mcp.NewServerTool("type-check", "Run TypeScript type checking on files or projects", s.handlers.TypeCheckHandler)
	getTypesTool := mcp.NewServerTool("get-types", "Extract type information for symbols in TypeScript files", s.handlers.GetTypesHandler)
	lintCheckTool := mcp.NewServerTool("lint-check", "Run ESLint checking on TypeScript files", s.handlers.LintCheckHandler)
	suggestImprovementsTool := mcp.NewServerTool("suggest-improvements", "Analyze TypeScript code and suggest improvements following best practices", s.handlers.SuggestImprovementsHandler)
	loadGuidelinesTool := mcp.NewServerTool("load-guidelines", "Load custom coding guidelines from markdown files", s.handlers.LoadGuidelinesHandler)

	// Add tools to server
	s.server.AddTools(typeCheckTool, getTypesTool, lintCheckTool, suggestImprovementsTool, loadGuidelinesTool)

	log.Println("Registered TypeScript MCP tools:")
	log.Println("- type-check: TypeScript type checking")
	log.Println("- get-types: Type information extraction")
	log.Println("- lint-check: ESLint checking")
	log.Println("- suggest-improvements: Code improvement suggestions")
	log.Println("- load-guidelines: Custom guideline loading")
}

// Run starts the MCP server with stdio transport
func (s *TypeScriptMCPServer) Run(ctx context.Context) error {
	log.Println("Starting TypeScript MCP Server...")
	log.Println("Server name: typescript-analyzer")
	log.Println("Version: 1.0.0")
	log.Println("Transport: stdio")
	
	// Check tool availability and log status
	s.logToolStatus()
	
	return s.server.Run(ctx, mcp.NewStdioTransport())
}

// logToolStatus logs the availability status of external tools
func (s *TypeScriptMCPServer) logToolStatus() {
	log.Println("Checking external tool availability...")
	
	if err := s.handlers.tscTool.CheckTSCAvailable(); err != nil {
		log.Printf("Warning: TypeScript compiler not available: %v", err)
		log.Println("  Make sure 'tsc' is installed (npm install -g typescript)")
	} else {
		if version, err := s.handlers.tscTool.GetVersion(); err == nil {
			log.Printf("TypeScript compiler available: %s", version)
		} else {
			log.Println("TypeScript compiler available")
		}
	}
	
	if err := s.handlers.eslintTool.CheckESLintAvailable(); err != nil {
		log.Printf("Warning: ESLint not available: %v", err)
		log.Println("  Make sure 'eslint' is installed (npm install -g eslint)")
	} else {
		if version, err := s.handlers.eslintTool.GetVersion(); err == nil {
			log.Printf("ESLint available: %s", version)
		} else {
			log.Println("ESLint available")
		}
	}
}

// Shutdown gracefully shuts down the server
func (s *TypeScriptMCPServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down TypeScript MCP Server...")
	return nil
}
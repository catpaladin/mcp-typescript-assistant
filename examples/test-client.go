package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp-typescript-assistant/pkg/types"
)

// TestClient provides a simple CLI client to test the TypeScript MCP server
type TestClient struct {
	client *mcp.Client
}

// NewTestClient creates a new test client
func NewTestClient() *TestClient {
	return &TestClient{}
}

// Connect connects to the MCP server
func (tc *TestClient) Connect() error {
	// Start the MCP server as a subprocess
	cmd := exec.Command("./mcp-typescript-assistant")
	
	transport, err := mcp.NewStdioTransport(cmd)
	if err != nil {
		return fmt.Errorf("failed to create transport: %w", err)
	}

	client := mcp.NewClient(&mcp.ClientImplementation{
		Name:    "typescript-test-client",
		Version: "1.0.0",
	})

	tc.client = client
	
	ctx := context.Background()
	return client.Connect(ctx, transport)
}

// TestTypeCheck tests the type-check tool
func (tc *TestClient) TestTypeCheck(filePath string) error {
	fmt.Printf("Testing type-check for file: %s\n", filePath)
	
	params := types.TypeCheckParams{
		FilePath: filePath,
	}
	
	ctx := context.Background()
	result, err := tc.client.CallTool(ctx, "type-check", params)
	if err != nil {
		return fmt.Errorf("type-check failed: %w", err)
	}
	
	fmt.Println("Type Check Result:")
	tc.printResult(result)
	return nil
}

// TestLintCheck tests the lint-check tool
func (tc *TestClient) TestLintCheck(filePath string) error {
	fmt.Printf("Testing lint-check for file: %s\n", filePath)
	
	params := types.LintCheckParams{
		FilePath: filePath,
	}
	
	ctx := context.Background()
	result, err := tc.client.CallTool(ctx, "lint-check", params)
	if err != nil {
		return fmt.Errorf("lint-check failed: %w", err)
	}
	
	fmt.Println("Lint Check Result:")
	tc.printResult(result)
	return nil
}

// TestSuggestImprovements tests the suggest-improvements tool
func (tc *TestClient) TestSuggestImprovements(codeSnippet string) error {
	fmt.Println("Testing suggest-improvements...")
	
	params := types.SuggestImprovementsParams{
		CodeSnippet: codeSnippet,
		Context:     "test code",
	}
	
	ctx := context.Background()
	result, err := tc.client.CallTool(ctx, "suggest-improvements", params)
	if err != nil {
		return fmt.Errorf("suggest-improvements failed: %w", err)
	}
	
	fmt.Println("Improvement Suggestions:")
	tc.printResult(result)
	return nil
}

// TestLoadGuidelines tests the load-guidelines tool
func (tc *TestClient) TestLoadGuidelines(guidelinePath string) error {
	fmt.Printf("Testing load-guidelines for file: %s\n", guidelinePath)
	
	params := types.LoadGuidelinesParams{
		GuidelinePath: guidelinePath,
		GuidelineType: "team-standards",
	}
	
	ctx := context.Background()
	result, err := tc.client.CallTool(ctx, "load-guidelines", params)
	if err != nil {
		return fmt.Errorf("load-guidelines failed: %w", err)
	}
	
	fmt.Println("Guidelines Loading Result:")
	tc.printResult(result)
	return nil
}

// printResult prints the result in a formatted way
func (tc *TestClient) printResult(result interface{}) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting result: %v\n", err)
		return
	}
	fmt.Println(string(data))
	fmt.Println()
}

// Disconnect closes the connection to the MCP server
func (tc *TestClient) Disconnect() error {
	if tc.client != nil {
		ctx := context.Background()
		return tc.client.Close(ctx)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: test-client <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  type-check <file>")
		fmt.Println("  lint-check <file>")
		fmt.Println("  suggest-improvements <code>")
		fmt.Println("  load-guidelines <file>")
		fmt.Println("  interactive")
		return
	}

	client := NewTestClient()
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Disconnect()

	command := os.Args[1]
	
	switch command {
	case "type-check":
		if len(os.Args) < 3 {
			log.Fatal("type-check requires a file path")
		}
		if err := client.TestTypeCheck(os.Args[2]); err != nil {
			log.Fatalf("Test failed: %v", err)
		}
		
	case "lint-check":
		if len(os.Args) < 3 {
			log.Fatal("lint-check requires a file path")
		}
		if err := client.TestLintCheck(os.Args[2]); err != nil {
			log.Fatalf("Test failed: %v", err)
		}
		
	case "suggest-improvements":
		if len(os.Args) < 3 {
			log.Fatal("suggest-improvements requires code snippet")
		}
		if err := client.TestSuggestImprovements(os.Args[2]); err != nil {
			log.Fatalf("Test failed: %v", err)
		}
		
	case "load-guidelines":
		if len(os.Args) < 3 {
			log.Fatal("load-guidelines requires a file path")
		}
		if err := client.TestLoadGuidelines(os.Args[2]); err != nil {
			log.Fatalf("Test failed: %v", err)
		}
		
	case "interactive":
		runInteractiveMode(client)
		
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// runInteractiveMode runs an interactive testing session
func runInteractiveMode(client *TestClient) {
	fmt.Println("=== TypeScript MCP Server Interactive Test ===")
	fmt.Println()
	
	// Test sample code improvements
	sampleCode := `
function getUserData(id) {
  return fetch('/api/users/' + id)
    .then(response => response.json())
    .then(data => data)
    .catch(err => console.log(err));
}

interface user {
  name: string;
  age: number;
}
`
	
	fmt.Println("1. Testing code improvement suggestions...")
	if err := client.TestSuggestImprovements(sampleCode); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("2. Testing guideline loading...")
	if err := client.TestLoadGuidelines("./examples/sample-guidelines.md"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("3. Testing improvements with loaded guidelines...")
	if err := client.TestSuggestImprovements(sampleCode); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Println("=== Interactive Test Complete ===")
}
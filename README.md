# TypeScript MCP Server

A Model Context Protocol (MCP) server that provides comprehensive TypeScript development
tools and best practices analysis. This server enables LLMs to analyze TypeScript code,
provide type information, run linting checks, and suggest improvements following
TypeScript best practices.

## Features

### Core Tools

1. **type-check** - TypeScript type checking

   - Run `tsc --noEmit` on files or entire projects
   - Parse and structure compiler errors and warnings
   - Support for both single files and project-wide checking

2. **get-types** - Type information extraction

   - Extract type information for symbols in TypeScript files
   - Analyze type definitions and interfaces
   - Provide detailed type metadata

3. **lint-check** - ESLint integration

   - Run ESLint with TypeScript-specific rules
   - Parse linting results with fix suggestions
   - Support for custom rule configurations

4. **suggest-improvements** - Code analysis and suggestions

   - Analyze TypeScript code for best practices
   - Provide actionable improvement recommendations
   - Support for custom coding guidelines

5. **load-guidelines** - Custom guideline support
   - Load coding standards from markdown files
   - Parse team-specific rules and conventions
   - Apply custom guidelines in code analysis

### Key Capabilities

- **TypeScript Integration**: Direct integration with TypeScript compiler (tsc) and
  ESLint
- **Best Practices Analysis**: Built-in knowledge of TypeScript conventions and patterns
- **Custom Guidelines**: Support for team-specific coding standards
- **Structured Output**: JSON-formatted results for easy integration
- **Error Handling**: Robust error handling and logging
- **Tool Availability Checking**: Automatic detection of required tools

## Installation

### Prerequisites

1. **Go 1.24.0+** - Required for building the MCP server
2. **TypeScript** - For type checking functionality
   ```bash
   npm install -g typescript
   ```
3. **ESLint** - For linting functionality
   ```bash
   npm install -g eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
   ```

### Building the Server

```bash
# Clone and build
git clone <repository-url>
cd mcp-typescript-assistant
go mod tidy
go build -o mcp-typescript-assistant .
```

### MCP Client Configuration

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "typescript-analyzer": {
      "command": "./mcp-typescript-assistant",
      "args": [],
      "env": {
        "DEBUG": "false"
      }
    }
  }
}
```

## Usage

### Tool Examples

#### Type Checking

```json
{
  "tool": "type-check",
  "arguments": {
    "file_path": "./src/app.ts",
    "project_root": "./src"
  }
}
```

#### Lint Checking

```json
{
  "tool": "lint-check",
  "arguments": {
    "file_path": "./src/app.ts",
    "rules": ["@typescript-eslint/no-unused-vars", "prefer-const"]
  }
}
```

#### Code Improvement Suggestions

```json
{
  "tool": "suggest-improvements",
  "arguments": {
    "code_snippet": "function getData(id) { return fetch('/api/' + id); }",
    "context": "API utility function"
  }
}
```

#### Loading Custom Guidelines

```json
{
  "tool": "load-guidelines",
  "arguments": {
    "guideline_path": "./team-standards.md",
    "guideline_type": "team-conventions"
  }
}
```

### Comprehensive Example Prompts

For detailed examples of how to use each tool effectively, see these example files:

**[`examples/tool-prompts.md`](./examples/tool-prompts.md)** - Technical tool usage:
- Formal tool call syntax and JSON examples
- Parameter specifications and expected responses
- Workflow combinations and chaining

**[`examples/chat-prompts.md`](./examples/chat-prompts.md)** - Natural conversation examples:
- Real chat conversations showing how users naturally request help
- Casual language that triggers appropriate tool calls
- Complete workflow examples with back-and-forth dialogue
- Migration and debugging scenarios
- Tips for effective natural language interaction

### Testing with CLI Client

Build and run the test client:

```bash
go build -o test-client ./examples/test-client.go

# Test individual tools
./test-client type-check ./examples/sample.ts
./test-client lint-check ./examples/sample.ts
./test-client suggest-improvements "let data = fetchData();"
./test-client load-guidelines ./examples/sample-guidelines.md

# Run interactive test session
./test-client interactive
```

## Project Structure

```
mcp-typescript-assistant/
├── cmd/server/          # Server executable
├── internal/
│   ├── server/          # MCP server implementation
│   ├── tools/           # TypeScript and ESLint integrations
│   ├── typescript/      # Code analysis and improvement engine
│   └── guidelines/      # Guideline parsing and management
├── pkg/types/           # Shared type definitions
├── examples/            # Example configurations and test client
├── go.mod
├── main.go             # Main server entry point
└── README.md
```

## Custom Guidelines

Create markdown files with your coding standards:

````markdown
# Team TypeScript Guidelines

## Naming Conventions

Use PascalCase for types and interfaces, camelCase for variables.

### Good Example

```typescript
interface UserProfile {
  firstName: string;
}
```
````

### Rules

- prefer-interface
- camelcase-variables

````

The server will parse these guidelines and apply them during code analysis.

## Development

### Building from Source

```bash
go mod tidy
go build -o mcp-typescript-assistant .
````

### Running in Debug Mode

```bash
DEBUG=true ./mcp-typescript-assistant
```

### Running Tests

```bash
go test ./...
```

## Best Practices Implemented

The server enforces and suggests:

- **Type Safety**: Explicit type annotations, avoiding `any`
- **Modern TypeScript**: ES modules, utility types, strict checking
- **Naming Conventions**: PascalCase for types, camelCase for variables
- **Async Patterns**: Proper async/await usage and error handling
- **Import/Export**: Named exports, organized imports
- **Performance**: Tree-shaking friendly patterns

## Troubleshooting

### Common Issues

1. **"TypeScript compiler not available"**

   ```bash
   npm install -g typescript
   ```

2. **"ESLint not available"**

   ```bash
   npm install -g eslint @typescript-eslint/parser @typescript-eslint/eslint-plugin
   ```

3. **"No such file or directory"**

   - Ensure file paths are absolute or relative to the working directory
   - Check file permissions

4. **MCP Connection Issues**
   - Verify the server binary is executable
   - Check MCP client configuration
   - Review server logs (stderr output)

### Debug Mode

Enable debug logging:

```bash
DEBUG=true ./mcp-typescript-assistant
```


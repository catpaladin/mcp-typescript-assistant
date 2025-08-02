# Example Tool Call Prompts

This document provides example prompts and scenarios for using each tool in the
TypeScript MCP server.

## 1. type-check Tool

### Basic File Type Checking

**Prompt:** "Please run type checking on the TypeScript file located at
`./src/components/UserProfile.ts`"

**Tool Call:**

```json
{
  "tool": "type-check",
  "arguments": {
    "file_path": "./src/components/UserProfile.ts"
  }
}
```

### Project-Wide Type Checking

**Prompt:** "Check all TypeScript files in my project for type errors. The project root
is `./frontend` and has a tsconfig.json file."

**Tool Call:**

```json
{
  "tool": "type-check",
  "arguments": {
    "file_path": "",
    "project_root": "./frontend"
  }
}
```

### Specific File in Project Context

**Prompt:** "I need to verify the types in `./src/utils/api.ts` but check it against the
project configuration in `./src`"

**Tool Call:**

```json
{
  "tool": "type-check",
  "arguments": {
    "file_path": "./src/utils/api.ts",
    "project_root": "./src"
  }
}
```

## 2. get-types Tool

### Extract All Types from File

**Prompt:** "Show me all the type information available in the file `./types/user.ts`"

**Tool Call:**

```json
{
  "tool": "get-types",
  "arguments": {
    "file_path": "./types/user.ts"
  }
}
```

### Get Specific Symbol Type

**Prompt:** "What is the type definition for the `UserProfile` interface in
`./src/models/user.ts`?"

**Tool Call:**

```json
{
  "tool": "get-types",
  "arguments": {
    "file_path": "./src/models/user.ts",
    "symbol_name": "UserProfile"
  }
}
```

### Analyze Function Types

**Prompt:** "I want to understand the type signature of the `processData` function in
`./utils/data-processor.ts`"

**Tool Call:**

```json
{
  "tool": "get-types",
  "arguments": {
    "file_path": "./utils/data-processor.ts",
    "symbol_name": "processData"
  }
}
```

## 3. lint-check Tool

### Basic Linting

**Prompt:** "Run ESLint on `./src/app.ts` to check for code quality issues"

**Tool Call:**

```json
{
  "tool": "lint-check",
  "arguments": {
    "file_path": "./src/app.ts"
  }
}
```

### Linting with Specific Rules

**Prompt:** "Check `./components/Button.tsx` specifically for unused variables and
missing return types"

**Tool Call:**

```json
{
  "tool": "lint-check",
  "arguments": {
    "file_path": "./components/Button.tsx",
    "rules": [
      "@typescript-eslint/no-unused-vars",
      "@typescript-eslint/explicit-function-return-type"
    ]
  }
}
```

### Multiple Rule Checking

**Prompt:** "I want to enforce strict TypeScript rules on `./services/auth.ts` - check
for any usage, missing types, and prefer const assertions"

**Tool Call:**

```json
{
  "tool": "lint-check",
  "arguments": {
    "file_path": "./services/auth.ts",
    "rules": [
      "@typescript-eslint/no-explicit-any",
      "@typescript-eslint/explicit-module-boundary-types",
      "prefer-const",
      "@typescript-eslint/prefer-as-const"
    ]
  }
}
```

## 4. suggest-improvements Tool

### Code Snippet Analysis

**Prompt:** "Can you analyze this TypeScript code and suggest improvements based on best
practices?"

**Code:**

```typescript
function fetchUser(id) {
  return fetch("/api/users/" + id)
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
      return data;
    })
    .catch((err) => console.error(err));
}
```

**Tool Call:**

```json
{
  "tool": "suggest-improvements",
  "arguments": {
    "code_snippet": "function fetchUser(id) {\n  return fetch('/api/users/' + id)\n    .then(response => response.json())\n    .then(data => {\n      console.log(data);\n      return data;\n    })\n    .catch(err => console.error(err));\n}",
    "context": "API utility function for fetching user data"
  }
}
```

### Interface Definition Review

**Prompt:** "Review this interface definition and suggest TypeScript best practices
improvements:"

**Code:**

```typescript
interface user {
  name: string;
  age?: number;
  email: string;
  preferences: any;
}
```

**Tool Call:**

```json
{
  "tool": "suggest-improvements",
  "arguments": {
    "code_snippet": "interface user {\n  name: string;\n  age?: number;\n  email: string;\n  preferences: any;\n}",
    "context": "User data model for application state"
  }
}
```

### Class Implementation Analysis

**Prompt:** "Analyze this TypeScript class and suggest improvements for better type
safety and modern patterns:"

**Code:**

```typescript
class DataManager {
  private data: any[] = [];

  addItem(item) {
    this.data.push(item);
  }

  getItem(index) {
    return this.data[index];
  }
}
```

**Tool Call:**

```json
{
  "tool": "suggest-improvements",
  "arguments": {
    "code_snippet": "class DataManager {\n  private data: any[] = [];\n  \n  addItem(item) {\n    this.data.push(item);\n  }\n  \n  getItem(index) {\n    return this.data[index];\n  }\n}",
    "context": "Generic data management utility class"
  }
}
```

## 5. load-guidelines Tool

### Loading Team Guidelines

**Prompt:** "Load our team's TypeScript coding standards from the guidelines file so
they can be applied to code analysis"

**Tool Call:**

```json
{
  "tool": "load-guidelines",
  "arguments": {
    "guideline_path": "./docs/typescript-guidelines.md",
    "guideline_type": "team-standards"
  }
}
```

### Loading Project-Specific Rules

**Prompt:** "Load the specific coding conventions for this React project from our style
guide"

**Tool Call:**

```json
{
  "tool": "load-guidelines",
  "arguments": {
    "guideline_path": "./project-conventions.md",
    "guideline_type": "react-project-rules"
  }
}
```

### Loading Company-Wide Standards

**Prompt:** "Import our company's enterprise TypeScript standards to ensure compliance"

**Tool Call:**

```json
{
  "tool": "load-guidelines",
  "arguments": {
    "guideline_path": "/company-standards/typescript-enterprise-guidelines.md",
    "guideline_type": "enterprise-standards"
  }
}
```

## Common Workflow Examples

### Complete Code Review Workflow

**Prompt:** "I want to do a complete review of `./src/user-service.ts` - check types,
run linting, and get improvement suggestions"

**Workflow:**

```json
// Step 1: Type check
{
  "tool": "type-check",
  "arguments": {
    "file_path": "./src/user-service.ts"
  }
}

// Step 2: Lint check
{
  "tool": "lint-check",
  "arguments": {
    "file_path": "./src/user-service.ts"
  }
}

// Step 3: Get improvement suggestions (assuming you have the file content)
{
  "tool": "suggest-improvements",
  "arguments": {
    "code_snippet": "[file content here]",
    "context": "User service with authentication and data management"
  }
}
```

### Setting Up Project Analysis

**Prompt:** "Set up comprehensive analysis for our TypeScript project - load our
guidelines and then analyze the main application file"

**Workflow:**

```json
// Step 1: Load guidelines
{
  "tool": "load-guidelines",
  "arguments": {
    "guideline_path": "./coding-standards.md",
    "guideline_type": "project-guidelines"
  }
}

// Step 2: Analyze with loaded guidelines
{
  "tool": "suggest-improvements",
  "arguments": {
    "code_snippet": "[application code]",
    "context": "Main application entry point"
  }
}
```

## Expected Response Formats

### Type Check Response

```json
{
  "success": true,
  "errors": [
    {
      "file": "./src/app.ts",
      "line": 15,
      "column": 23,
      "message": "Property 'name' does not exist on type 'User'",
      "code": "TS2339",
      "severity": "error"
    }
  ],
  "warnings": [],
  "compile_time": "1.2s"
}
```

### Improvement Suggestions Response

```json
{
  "improvements": [
    {
      "type": "type_annotation",
      "description": "Add explicit type annotations to function parameters",
      "reasoning": "Explicit types improve code readability and catch type errors early",
      "priority": "high"
    }
  ],
  "summary": "Found 3 improvement suggestions: 1 high priority, 2 medium priority",
  "applied_rules": ["typescript-standard-practices", "team-guidelines"]
}
```

## Tips for Effective Tool Usage

1. **Use absolute paths** when possible to avoid confusion
2. **Provide context** in improvement suggestions for better analysis
3. **Load guidelines first** if you want custom rules applied
4. **Chain tool calls** for comprehensive analysis
5. **Check file permissions** if you get "file not found" errors


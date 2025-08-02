package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp-typescript-assistant/internal/guidelines"
	"mcp-typescript-assistant/internal/tools"
	"mcp-typescript-assistant/internal/typescript"
	"mcp-typescript-assistant/pkg/types"
)

// Handlers contains all the tool handlers for the MCP server
type Handlers struct {
	tscTool     *tools.TypeScriptCompiler
	eslintTool  *tools.ESLintTool
	analyzer    *typescript.Analyzer
	parser      *guidelines.Parser
}

// NewHandlers creates a new handlers instance
func NewHandlers() *Handlers {
	return &Handlers{
		tscTool:     tools.NewTypeScriptCompiler(),
		eslintTool:  tools.NewESLintTool(),
		analyzer:    typescript.NewAnalyzer(),
		parser:      guidelines.NewParser(),
	}
}

// TypeCheckHandler handles TypeScript type checking requests
func (h *Handlers) TypeCheckHandler(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[types.TypeCheckParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := h.tscTool.TypeCheck(params.Arguments)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error performing type check: %v", err),
				},
			},
		}, nil
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error marshaling result: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil
}

// GetTypesHandler handles type information extraction requests
func (h *Handlers) GetTypesHandler(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[types.GetTypesParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := h.tscTool.GetTypes(params.Arguments)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error extracting type information: %v", err),
				},
			},
		}, nil
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error marshaling result: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil
}

// LintCheckHandler handles ESLint checking requests
func (h *Handlers) LintCheckHandler(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[types.LintCheckParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := h.eslintTool.LintCheck(params.Arguments)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error performing lint check: %v", err),
				},
			},
		}, nil
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error marshaling result: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil
}

// SuggestImprovementsHandler handles code improvement suggestion requests
func (h *Handlers) SuggestImprovementsHandler(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[types.SuggestImprovementsParams]) (*mcp.CallToolResultFor[any], error) {
	result, err := h.analyzer.SuggestImprovements(params.Arguments)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error suggesting improvements: %v", err),
				},
			},
		}, nil
	}

	resultJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error marshaling result: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil
}

// LoadGuidelinesHandler handles guideline loading requests
func (h *Handlers) LoadGuidelinesHandler(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[types.LoadGuidelinesParams]) (*mcp.CallToolResultFor[any], error) {
	guidelineSet, err := h.parser.ParseGuidelinesFromFile(params.Arguments.GuidelinePath, params.Arguments.GuidelineType)
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error loading guidelines: %v", err),
				},
			},
		}, nil
	}

	// Validate guidelines
	warnings := h.parser.ValidateGuidelines(guidelineSet)
	
	// Load guidelines into analyzer
	h.analyzer.LoadGuidelines(guidelineSet)

	response := map[string]interface{}{
		"success":        true,
		"guideline_set":  guidelineSet,
		"warnings":       warnings,
		"message":        fmt.Sprintf("Successfully loaded %d guidelines from %s", len(guidelineSet.Guidelines), guidelineSet.Name),
	}

	resultJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error marshaling result: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil
}

// GetServerInfoHandler provides information about the server capabilities
func (h *Handlers) GetServerInfoHandler(ctx context.Context, cc *mcp.ServerSession) (*mcp.CallToolResultFor[any], error) {
	info := map[string]interface{}{
		"name":        "typescript-analyzer",
		"version":     "1.0.0",
		"description": "TypeScript development tools and best practices analyzer",
		"tools": []string{
			"type-check",
			"get-types", 
			"lint-check",
			"suggest-improvements",
			"load-guidelines",
		},
		"capabilities": map[string]bool{
			"typescript_compilation": true,
			"eslint_integration":     true,
			"code_analysis":          true,
			"custom_guidelines":      true,
			"type_extraction":        true,
		},
	}

	// Check tool availability
	toolStatus := make(map[string]bool)
	toolStatus["typescript"] = h.tscTool.CheckTSCAvailable() == nil
	toolStatus["eslint"] = h.eslintTool.CheckESLintAvailable() == nil

	info["tool_status"] = toolStatus

	// Get versions if available
	versions := make(map[string]string)
	if tscVersion, err := h.tscTool.GetVersion(); err == nil {
		versions["typescript"] = tscVersion
	}
	if eslintVersion, err := h.eslintTool.GetVersion(); err == nil {
		versions["eslint"] = eslintVersion
	}
	info["versions"] = versions

	// Get loaded guidelines
	loadedGuidelines := h.analyzer.GetLoadedGuidelines()
	guidelineNames := make([]string, 0, len(loadedGuidelines))
	for name := range loadedGuidelines {
		guidelineNames = append(guidelineNames, name)
	}
	info["loaded_guidelines"] = guidelineNames

	resultJSON, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error marshaling server info: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(resultJSON),
			},
		},
	}, nil
}
package tools

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"mcp-typescript-assistant/pkg/types"
)

// TypeScriptCompiler provides TypeScript compilation and type checking capabilities
type TypeScriptCompiler struct {
	tscPath string
}

// NewTypeScriptCompiler creates a new TypeScript compiler instance
func NewTypeScriptCompiler() *TypeScriptCompiler {
	tscPath := "tsc"
	if path, err := exec.LookPath("npx"); err == nil {
		tscPath = path
	}
	return &TypeScriptCompiler{tscPath: tscPath}
}

// TypeCheck performs TypeScript type checking on a file or project
func (tsc *TypeScriptCompiler) TypeCheck(params types.TypeCheckParams) (*types.TypeCheckResult, error) {
	startTime := time.Now()

	var args []string
	if tsc.tscPath == "npx" {
		args = append(args, "tsc")
	}

	args = append(args, "--noEmit", "--pretty", "false")

	if params.ProjectRoot != "" {
		// Check for project compilation
		configPath := filepath.Join(params.ProjectRoot, "tsconfig.json")
		args = append(args, "--project", configPath)
	} else {
		// Single file compilation
		args = append(args, params.FilePath)
	}

	cmd := exec.Command(tsc.tscPath, args...)
	if params.ProjectRoot != "" {
		cmd.Dir = params.ProjectRoot
	}

	output, err := cmd.CombinedOutput()
	compileTime := time.Since(startTime).String()

	result := &types.TypeCheckResult{
		Success:     err == nil,
		CompileTime: compileTime,
	}

	if len(output) > 0 {
		errors, warnings := tsc.parseTypeScriptOutput(string(output))
		result.Errors = errors
		result.Warnings = warnings
	}

	return result, nil
}

// GetTypes extracts type information for a symbol in a TypeScript file
func (tsc *TypeScriptCompiler) GetTypes(params types.GetTypesParams) (*types.TypeInfo, error) {
	// This would ideally use the TypeScript Language Service API
	// For now, we'll use a simplified approach with compilation output

	args := []string{"--noEmit", "--listFiles", params.FilePath}
	if tsc.tscPath == "npx" {
		args = append([]string{"tsc"}, args...)
	}

	cmd := exec.Command(tsc.tscPath, args...)
	_, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("failed to analyze types: %w", err)
	}

	// For a more complete implementation, this would parse the TypeScript AST
	// and extract detailed type information
	typeInfo := &types.TypeInfo{
		SymbolName: params.SymbolName,
		Type:       "unknown",  // Would be extracted from AST analysis
		Kind:       "variable", // Would be determined from AST
		Location: &types.SourceLocation{
			File: params.FilePath,
		},
	}

	if params.SymbolName == "" {
		typeInfo.SymbolName = "file_analysis"
		typeInfo.Documentation = fmt.Sprintf("Type analysis for file: %s", params.FilePath)
	}

	return typeInfo, nil
}

// parseTypeScriptOutput parses TypeScript compiler output into structured errors and warnings
func (tsc *TypeScriptCompiler) parseTypeScriptOutput(output string) ([]types.TypeScriptError, []types.TypeScriptError) {
	var errors []types.TypeScriptError
	var warnings []types.TypeScriptError

	// TypeScript error format: file(line,column): error TS####: message
	errorRegex := regexp.MustCompile(`^(.+?)\((\d+),(\d+)\):\s+(error|warning)\s+TS(\d+):\s+(.+)$`)

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := errorRegex.FindStringSubmatch(line)
		if len(matches) == 7 {
			lineNum, _ := strconv.Atoi(matches[2])
			colNum, _ := strconv.Atoi(matches[3])
			severity := matches[4]
			code := "TS" + matches[5]
			message := matches[6]

			tsError := types.TypeScriptError{
				File:     matches[1],
				Line:     lineNum,
				Column:   colNum,
				Message:  message,
				Code:     code,
				Severity: severity,
			}

			if severity == "error" {
				errors = append(errors, tsError)
			} else {
				warnings = append(warnings, tsError)
			}
		}
	}

	return errors, warnings
}

// CheckTSCAvailable checks if TypeScript compiler is available
func (tsc *TypeScriptCompiler) CheckTSCAvailable() error {
	var cmd *exec.Cmd
	if tsc.tscPath == "npx" {
		cmd = exec.Command("npx", "tsc", "--version")
	} else {
		cmd = exec.Command("tsc", "--version")
	}

	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("TypeScript compiler not available: %w", err)
	}
	return nil
}

// GetVersion returns the TypeScript compiler version
func (tsc *TypeScriptCompiler) GetVersion() (string, error) {
	var cmd *exec.Cmd
	if tsc.tscPath == "npx" {
		cmd = exec.Command("npx", "tsc", "--version")
	} else {
		cmd = exec.Command("tsc", "--version")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(output))
	return version, nil
}


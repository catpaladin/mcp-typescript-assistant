package tools

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"mcp-typescript-assistant/pkg/types"
)

// ESLintTool provides ESLint integration for TypeScript files
type ESLintTool struct {
	eslintPath string
}

// NewESLintTool creates a new ESLint tool instance
func NewESLintTool() *ESLintTool {
	eslintPath := "eslint"
	if path, err := exec.LookPath("npx"); err == nil {
		eslintPath = path
	}
	return &ESLintTool{eslintPath: eslintPath}
}

// ESLintOutput represents the JSON output from ESLint
type ESLintOutput struct {
	FilePath            string          `json:"filePath"`
	Messages            []ESLintMessage `json:"messages"`
	ErrorCount          int             `json:"errorCount"`
	WarningCount        int             `json:"warningCount"`
	FixableErrorCount   int             `json:"fixableErrorCount"`
	FixableWarningCount int             `json:"fixableWarningCount"`
}

// ESLintMessage represents a single ESLint message
type ESLintMessage struct {
	RuleID   string `json:"ruleId"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	NodeType string `json:"nodeType,omitempty"`
	Fix      *Fix   `json:"fix,omitempty"`
}

// Fix represents an ESLint fix suggestion
type Fix struct {
	Range []int  `json:"range"`
	Text  string `json:"text"`
}

// LintCheck performs ESLint checking on a TypeScript file
func (eslint *ESLintTool) LintCheck(params types.LintCheckParams) (*types.LintResult, error) {
	var args []string

	if eslint.eslintPath == "npx" {
		args = append(args, "eslint")
	}

	args = append(args, "--format", "json")

	if len(params.Rules) > 0 {
		// Add specific rules
		for _, rule := range params.Rules {
			args = append(args, "--rule", rule)
		}
	}

	args = append(args, params.FilePath)

	cmd := exec.Command(eslint.eslintPath, args...)
	output, err := cmd.Output()

	// ESLint returns non-zero exit code when there are linting errors
	// but we still want to parse the output
	result := &types.LintResult{
		Success: err == nil,
	}

	if len(output) > 0 {
		issues, fixableCount := eslint.parseESLintOutput(output)
		result.Issues = issues
		result.Fixable = fixableCount
		result.Summary = eslint.generateSummary(issues, fixableCount)
	} else if err != nil {
		// If there's an error and no output, ESLint might not be configured properly
		return nil, fmt.Errorf("ESLint execution failed: %w", err)
	}

	return result, nil
}

// parseESLintOutput parses ESLint JSON output into structured issues
func (eslint *ESLintTool) parseESLintOutput(output []byte) ([]types.LintIssue, int) {
	var eslintResults []ESLintOutput
	if err := json.Unmarshal(output, &eslintResults); err != nil {
		return nil, 0
	}

	var issues []types.LintIssue
	fixableCount := 0

	for _, result := range eslintResults {
		for _, message := range result.Messages {
			severity := "warning"
			if message.Severity == 2 {
				severity = "error"
			}

			fixable := message.Fix != nil
			if fixable {
				fixableCount++
			}

			issue := types.LintIssue{
				File:     result.FilePath,
				Line:     message.Line,
				Column:   message.Column,
				Message:  message.Message,
				Rule:     message.RuleID,
				Severity: severity,
				Fixable:  fixable,
			}
			issues = append(issues, issue)
		}
	}

	return issues, fixableCount
}

// generateSummary creates a summary of linting results
func (eslint *ESLintTool) generateSummary(issues []types.LintIssue, fixableCount int) string {
	if len(issues) == 0 {
		return "No linting issues found"
	}

	errorCount := 0
	warningCount := 0

	for _, issue := range issues {
		if issue.Severity == "error" {
			errorCount++
		} else {
			warningCount++
		}
	}

	summary := fmt.Sprintf("Found %d issue(s): %d error(s), %d warning(s)",
		len(issues), errorCount, warningCount)

	if fixableCount > 0 {
		summary += fmt.Sprintf(", %d fixable", fixableCount)
	}

	return summary
}

// AutoFix attempts to automatically fix ESLint issues
func (eslint *ESLintTool) AutoFix(filePath string) (*types.LintResult, error) {
	var args []string

	if eslint.eslintPath == "npx" {
		args = append(args, "eslint")
	}

	args = append(args, "--fix", "--format", "json", filePath)

	cmd := exec.Command(eslint.eslintPath, args...)
	output, err := cmd.Output()

	result := &types.LintResult{
		Success: err == nil,
	}

	if len(output) > 0 {
		issues, fixableCount := eslint.parseESLintOutput(output)
		result.Issues = issues
		result.Fixable = fixableCount
		result.Summary = eslint.generateSummary(issues, fixableCount)
	}

	return result, nil
}

// CheckESLintAvailable checks if ESLint is available
func (eslint *ESLintTool) CheckESLintAvailable() error {
	var cmd *exec.Cmd
	if eslint.eslintPath == "npx" {
		cmd = exec.Command("npx", "eslint", "--version")
	} else {
		cmd = exec.Command("eslint", "--version")
	}

	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("ESLint not available: %w", err)
	}
	return nil
}

// GetVersion returns the ESLint version
func (eslint *ESLintTool) GetVersion() (string, error) {
	var cmd *exec.Cmd
	if eslint.eslintPath == "npx" {
		cmd = exec.Command("npx", "eslint", "--version")
	} else {
		cmd = exec.Command("eslint", "--version")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	version := strings.TrimSpace(string(output))
	return version, nil
}

// GetConfig returns ESLint configuration for a file
func (eslint *ESLintTool) GetConfig(filePath string) (map[string]interface{}, error) {
	var args []string

	if eslint.eslintPath == "npx" {
		args = append(args, "eslint")
	}

	args = append(args, "--print-config", filePath)

	cmd := exec.Command(eslint.eslintPath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get ESLint config: %w", err)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(output, &config); err != nil {
		return nil, fmt.Errorf("failed to parse ESLint config: %w", err)
	}

	return config, nil
}


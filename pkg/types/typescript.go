package types

import "encoding/json"

// TypeCheckParams represents parameters for TypeScript type checking
type TypeCheckParams struct {
	FilePath    string `json:"file_path"`
	ProjectRoot string `json:"project_root,omitempty"`
}

// GetTypesParams represents parameters for getting type information
type GetTypesParams struct {
	FilePath   string `json:"file_path"`
	SymbolName string `json:"symbol_name,omitempty"`
}

// LintCheckParams represents parameters for ESLint checking
type LintCheckParams struct {
	FilePath string   `json:"file_path"`
	Rules    []string `json:"rules,omitempty"`
}

// SuggestImprovementsParams represents parameters for code improvement suggestions
type SuggestImprovementsParams struct {
	CodeSnippet string `json:"code_snippet"`
	Context     string `json:"context,omitempty"`
}

// LoadGuidelinesParams represents parameters for loading coding guidelines
type LoadGuidelinesParams struct {
	GuidelinePath string `json:"guideline_path"`
	GuidelineType string `json:"guideline_type,omitempty"`
}

// TypeCheckResult represents the result of TypeScript type checking
type TypeCheckResult struct {
	Success     bool               `json:"success"`
	Errors      []TypeScriptError  `json:"errors,omitempty"`
	Warnings    []TypeScriptError  `json:"warnings,omitempty"`
	CompileTime string             `json:"compile_time,omitempty"`
}

// TypeScriptError represents a TypeScript compiler error or warning
type TypeScriptError struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Code     string `json:"code,omitempty"`
	Severity string `json:"severity"`
}

// TypeInfo represents type information for a symbol
type TypeInfo struct {
	SymbolName   string            `json:"symbol_name"`
	Type         string            `json:"type"`
	Kind         string            `json:"kind"`
	Documentation string           `json:"documentation,omitempty"`
	Location     *SourceLocation   `json:"location,omitempty"`
	Properties   []PropertyInfo    `json:"properties,omitempty"`
}

// SourceLocation represents a location in source code
type SourceLocation struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column"`
}

// PropertyInfo represents information about a property or method
type PropertyInfo struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Optional      bool   `json:"optional"`
	Documentation string `json:"documentation,omitempty"`
}

// LintResult represents the result of ESLint checking
type LintResult struct {
	Success  bool        `json:"success"`
	Issues   []LintIssue `json:"issues,omitempty"`
	Fixable  int         `json:"fixable_count"`
	Summary  string      `json:"summary"`
}

// LintIssue represents an ESLint issue
type LintIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
	Fixable  bool   `json:"fixable"`
}

// Improvement represents a code improvement suggestion
type Improvement struct {
	Type         string `json:"type"`
	Description  string `json:"description"`
	Before       string `json:"before,omitempty"`
	After        string `json:"after,omitempty"`
	Reasoning    string `json:"reasoning"`
	Priority     string `json:"priority"`
	GuidelineRef string `json:"guideline_ref,omitempty"`
}

// ImprovementResult represents the result of improvement suggestions
type ImprovementResult struct {
	Improvements []Improvement `json:"improvements"`
	Summary      string        `json:"summary"`
	AppliedRules []string      `json:"applied_rules,omitempty"`
}

// Guideline represents a coding guideline
type Guideline struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Priority    string            `json:"priority"`
	Examples    []GuidelineExample `json:"examples,omitempty"`
	Rules       []string          `json:"rules,omitempty"`
}

// GuidelineExample represents an example in a guideline
type GuidelineExample struct {
	Title       string `json:"title"`
	Good        string `json:"good,omitempty"`
	Bad         string `json:"bad,omitempty"`
	Explanation string `json:"explanation"`
}

// GuidelineSet represents a collection of guidelines
type GuidelineSet struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	Description string      `json:"description"`
	Guidelines  []Guideline `json:"guidelines"`
	LoadedAt    string      `json:"loaded_at"`
}

// String methods for better logging
func (tc TypeCheckResult) String() string {
	data, _ := json.MarshalIndent(tc, "", "  ")
	return string(data)
}

func (ti TypeInfo) String() string {
	data, _ := json.MarshalIndent(ti, "", "  ")
	return string(data)
}

func (lr LintResult) String() string {
	data, _ := json.MarshalIndent(lr, "", "  ")
	return string(data)
}
package typescript

import (
	"fmt"
	"regexp"
	"strings"

	"mcp-typescript-assistant/pkg/types"
)

// Analyzer provides TypeScript code analysis and improvement suggestions
type Analyzer struct {
	guidelines map[string]*types.GuidelineSet
}

// NewAnalyzer creates a new TypeScript analyzer
func NewAnalyzer() *Analyzer {
	return &Analyzer{
		guidelines: make(map[string]*types.GuidelineSet),
	}
}

// SuggestImprovements analyzes TypeScript code and suggests improvements
func (a *Analyzer) SuggestImprovements(params types.SuggestImprovementsParams) (*types.ImprovementResult, error) {
	var improvements []types.Improvement
	var appliedRules []string

	// Analyze the code snippet for common TypeScript issues
	improvements = append(improvements, a.analyzeTypeAnnotations(params.CodeSnippet)...)
	improvements = append(improvements, a.analyzeNamingConventions(params.CodeSnippet)...)
	improvements = append(improvements, a.analyzeImportExports(params.CodeSnippet)...)
	improvements = append(improvements, a.analyzeAsyncAwait(params.CodeSnippet)...)
	improvements = append(improvements, a.analyzeTypeAssertions(params.CodeSnippet)...)
	improvements = append(improvements, a.analyzeUtilityTypes(params.CodeSnippet)...)

	// Apply custom guidelines if loaded
	for _, guidelineSet := range a.guidelines {
		guidelineImprovements := a.applyGuidelines(params.CodeSnippet, guidelineSet)
		improvements = append(improvements, guidelineImprovements...)
		appliedRules = append(appliedRules, guidelineSet.Name)
	}

	// Add standard TypeScript best practices
	appliedRules = append(appliedRules, "typescript-standard-practices")

	summary := a.generateImprovementSummary(improvements)

	return &types.ImprovementResult{
		Improvements: improvements,
		Summary:      summary,
		AppliedRules: appliedRules,
	}, nil
}

// analyzeTypeAnnotations checks for missing or incorrect type annotations
func (a *Analyzer) analyzeTypeAnnotations(code string) []types.Improvement {
	var improvements []types.Improvement

	// Check for implicit any types
	anyRegex := regexp.MustCompile(`let\s+\w+\s*=|const\s+\w+\s*=|var\s+\w+\s*=`)
	if anyRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "type_annotation",
			Description: "Consider adding explicit type annotations to variables",
			Reasoning:   "Explicit types improve code readability and catch type errors early",
			Priority:    "medium",
		})
	}

	// Check for function parameters without types
	paramRegex := regexp.MustCompile(`function\s+\w+\s*\([^)]*\w+[^:)]*\)`)
	if paramRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "function_types",
			Description: "Add type annotations to function parameters",
			Reasoning:   "Typed parameters prevent runtime errors and improve IDE support",
			Priority:    "high",
		})
	}

	return improvements
}

// analyzeNamingConventions checks naming conventions
func (a *Analyzer) analyzeNamingConventions(code string) []types.Improvement {
	var improvements []types.Improvement

	// Check for PascalCase interfaces/types
	interfaceRegex := regexp.MustCompile(`interface\s+([a-z][a-zA-Z]*)`)
	matches := interfaceRegex.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			improvements = append(improvements, types.Improvement{
				Type:        "naming_convention",
				Description: fmt.Sprintf("Interface '%s' should use PascalCase", match[1]),
				Before:      match[1],
				After:       strings.Title(match[1]),
				Reasoning:   "TypeScript convention uses PascalCase for types and interfaces",
				Priority:    "low",
			})
		}
	}

	// Check for camelCase functions/variables
	varRegex := regexp.MustCompile(`(?:let|const|var)\s+([A-Z][a-zA-Z]*)`)
	matches = varRegex.FindAllStringSubmatch(code, -1)
	for _, match := range matches {
		if len(match) > 1 {
			improvements = append(improvements, types.Improvement{
				Type:        "naming_convention",
				Description: fmt.Sprintf("Variable '%s' should use camelCase", match[1]),
				Before:      match[1],
				After:       strings.ToLower(match[1][:1]) + match[1][1:],
				Reasoning:   "TypeScript convention uses camelCase for variables and functions",
				Priority:    "low",
			})
		}
	}

	return improvements
}

// analyzeImportExports checks import/export best practices
func (a *Analyzer) analyzeImportExports(code string) []types.Improvement {
	var improvements []types.Improvement

	// Check for default exports that could be named exports
	defaultExportRegex := regexp.MustCompile(`export\s+default\s+(class|function|interface)`)
	if defaultExportRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "export_style",
			Description: "Consider using named exports instead of default exports",
			Reasoning:   "Named exports provide better tree-shaking and refactoring support",
			Priority:    "medium",
		})
	}

	// Check for import without file extension (simplified check)
	importRegex := regexp.MustCompile(`import.*from\s+['"]\.\/[^'"]*['"]`)
	if importRegex.MatchString(code) && !strings.Contains(code, ".ts\"") && !strings.Contains(code, ".js\"") {
		improvements = append(improvements, types.Improvement{
			Type:        "import_style",
			Description: "Consider adding explicit file extensions to relative imports",
			Reasoning:   "Explicit extensions improve compatibility with ES modules",
			Priority:    "low",
		})
	}

	return improvements
}

// analyzeAsyncAwait checks async/await usage
func (a *Analyzer) analyzeAsyncAwait(code string) []types.Improvement {
	var improvements []types.Improvement

	// Check for Promise.then() that could be async/await
	thenRegex := regexp.MustCompile(`\.then\s*\(`)
	if thenRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "async_pattern",
			Description: "Consider using async/await instead of .then()",
			Reasoning:   "Async/await provides better error handling and readability",
			Priority:    "medium",
		})
	}

	// Check for missing error handling in async functions
	asyncRegex := regexp.MustCompile(`async\s+function`)
	tryRegex := regexp.MustCompile(`try\s*{`)
	if asyncRegex.MatchString(code) && !tryRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "error_handling",
			Description: "Add error handling to async functions",
			Reasoning:   "Proper error handling prevents unhandled promise rejections",
			Priority:    "high",
		})
	}

	return improvements
}

// analyzeTypeAssertions checks type assertion usage
func (a *Analyzer) analyzeTypeAssertions(code string) []types.Improvement {
	var improvements []types.Improvement

	// Check for 'as any' assertions
	asAnyRegex := regexp.MustCompile(`as\s+any`)
	if asAnyRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "type_safety",
			Description: "Avoid using 'as any' type assertions",
			Reasoning:   "Type assertions bypass TypeScript's type checking and reduce type safety",
			Priority:    "high",
		})
	}

	// Check for angle bracket assertions (prefer 'as' syntax)
	angleBracketRegex := regexp.MustCompile(`<\w+>`)
	if angleBracketRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "assertion_style",
			Description: "Use 'as' syntax instead of angle bracket assertions",
			Reasoning:   "'as' syntax is preferred and works better with JSX",
			Priority:    "low",
		})
	}

	return improvements
}

// analyzeUtilityTypes suggests utility type usage
func (a *Analyzer) analyzeUtilityTypes(code string) []types.Improvement {
	var improvements []types.Improvement

	// Check for manual partial type definitions
	partialRegex := regexp.MustCompile(`{\s*\w+\?\s*:`)
	if partialRegex.MatchString(code) {
		improvements = append(improvements, types.Improvement{
			Type:        "utility_types",
			Description: "Consider using Partial<T> utility type",
			Reasoning:   "Utility types provide better type safety and maintainability",
			Priority:    "medium",
		})
	}

	// Check for manual pick/omit patterns
	if strings.Contains(code, "Pick") || strings.Contains(code, "Omit") {
		improvements = append(improvements, types.Improvement{
			Type:        "utility_types",
			Description: "Good use of utility types",
			Reasoning:   "Pick and Omit utility types provide excellent type manipulation",
			Priority:    "low",
		})
	}

	return improvements
}

// applyGuidelines applies custom guidelines to the code analysis
func (a *Analyzer) applyGuidelines(code string, guidelineSet *types.GuidelineSet) []types.Improvement {
	var improvements []types.Improvement

	for _, guideline := range guidelineSet.Guidelines {
		// Simple pattern matching for guidelines
		// In a real implementation, this would be more sophisticated
		for _, rule := range guideline.Rules {
			if strings.Contains(code, rule) {
				improvements = append(improvements, types.Improvement{
					Type:         "guideline",
					Description:  guideline.Description,
					Reasoning:    fmt.Sprintf("According to %s guidelines", guidelineSet.Name),
					Priority:     guideline.Priority,
					GuidelineRef: guideline.ID,
				})
			}
		}
	}

	return improvements
}

// generateImprovementSummary creates a summary of all improvements
func (a *Analyzer) generateImprovementSummary(improvements []types.Improvement) string {
	if len(improvements) == 0 {
		return "No improvements suggested. Code follows TypeScript best practices."
	}

	highPriority := 0
	mediumPriority := 0
	lowPriority := 0

	for _, improvement := range improvements {
		switch improvement.Priority {
		case "high":
			highPriority++
		case "medium":
			mediumPriority++
		case "low":
			lowPriority++
		}
	}

	return fmt.Sprintf("Found %d improvement suggestions: %d high priority, %d medium priority, %d low priority",
		len(improvements), highPriority, mediumPriority, lowPriority)
}

// LoadGuidelines loads custom guidelines from a guideline set
func (a *Analyzer) LoadGuidelines(guidelineSet *types.GuidelineSet) {
	a.guidelines[guidelineSet.Name] = guidelineSet
}

// GetLoadedGuidelines returns all loaded guidelines
func (a *Analyzer) GetLoadedGuidelines() map[string]*types.GuidelineSet {
	return a.guidelines
}
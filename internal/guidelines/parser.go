package guidelines

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"mcp-typescript-assistant/pkg/types"
)

// Parser handles parsing of markdown guideline files
type Parser struct {
	headerRegex   *regexp.Regexp
	codeRegex     *regexp.Regexp
	listRegex     *regexp.Regexp
}

// NewParser creates a new guideline parser
func NewParser() *Parser {
	return &Parser{
		headerRegex: regexp.MustCompile(`^#+\s+(.+)$`),
		codeRegex:   regexp.MustCompile("```(typescript|ts|javascript|js)?([\\s\\S]*?)```"),
		listRegex:   regexp.MustCompile(`^[\*\-\+]\s+(.+)$`),
	}
}

// ParseGuidelinesFromFile parses guidelines from a markdown file
func (p *Parser) ParseGuidelinesFromFile(filePath, guidelineType string) (*types.GuidelineSet, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open guideline file: %w", err)
	}
	defer file.Close()

	content, err := p.readFileContent(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read guideline file: %w", err)
	}

	return p.ParseGuidelines(content, filepath.Base(filePath), guidelineType)
}

// ParseGuidelines parses guidelines from markdown content
func (p *Parser) ParseGuidelines(content, name, guidelineType string) (*types.GuidelineSet, error) {
	if guidelineType == "" {
		guidelineType = "general"
	}

	guidelineSet := &types.GuidelineSet{
		Name:        name,
		Version:     "1.0.0",
		Description: fmt.Sprintf("%s coding guidelines", guidelineType),
		LoadedAt:    time.Now().Format(time.RFC3339),
	}

	guidelines := p.parseContent(content)
	guidelineSet.Guidelines = guidelines

	return guidelineSet, nil
}

// readFileContent reads the entire content of a file
func (p *Parser) readFileContent(file *os.File) (string, error) {
	var content strings.Builder
	scanner := bufio.NewScanner(file)
	
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}
	
	if err := scanner.Err(); err != nil {
		return "", err
	}
	
	return content.String(), nil
}

// parseContent parses markdown content into guidelines
func (p *Parser) parseContent(content string) []types.Guideline {
	var guidelines []types.Guideline
	
	sections := p.splitIntoSections(content)
	
	for i, section := range sections {
		guideline := p.parseSection(section, i+1)
		if guideline != nil {
			guidelines = append(guidelines, *guideline)
		}
	}
	
	return guidelines
}

// splitIntoSections splits content by headers
func (p *Parser) splitIntoSections(content string) []string {
	lines := strings.Split(content, "\n")
	var sections []string
	var currentSection strings.Builder
	
	for _, line := range lines {
		if p.headerRegex.MatchString(line) && currentSection.Len() > 0 {
			sections = append(sections, currentSection.String())
			currentSection.Reset()
		}
		currentSection.WriteString(line)
		currentSection.WriteString("\n")
	}
	
	if currentSection.Len() > 0 {
		sections = append(sections, currentSection.String())
	}
	
	return sections
}

// parseSection parses a single section into a guideline
func (p *Parser) parseSection(section string, id int) *types.Guideline {
	lines := strings.Split(section, "\n")
	if len(lines) == 0 {
		return nil
	}
	
	guideline := &types.Guideline{
		ID:       fmt.Sprintf("guideline_%d", id),
		Priority: "medium", // default priority
	}
	
	var currentContent strings.Builder
	var inCodeBlock bool
	var currentExample *types.GuidelineExample
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Parse headers
		if matches := p.headerRegex.FindStringSubmatch(line); len(matches) > 1 {
			if guideline.Title == "" {
				guideline.Title = matches[1]
				guideline.Category = p.inferCategory(matches[1])
				guideline.Priority = p.inferPriority(matches[1])
			}
			continue
		}
		
		// Parse code blocks
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// End of code block
				if currentExample != nil {
					if strings.Contains(strings.ToLower(currentExample.Title), "good") ||
					   strings.Contains(strings.ToLower(currentExample.Title), "correct") ||
					   strings.Contains(strings.ToLower(currentExample.Title), "do") {
						currentExample.Good = strings.TrimSpace(currentContent.String())
					} else if strings.Contains(strings.ToLower(currentExample.Title), "bad") ||
							  strings.Contains(strings.ToLower(currentExample.Title), "incorrect") ||
							  strings.Contains(strings.ToLower(currentExample.Title), "don't") {
						currentExample.Bad = strings.TrimSpace(currentContent.String())
					}
				}
				currentContent.Reset()
				inCodeBlock = false
			} else {
				// Start of code block
				inCodeBlock = true
				if currentExample == nil {
					currentExample = &types.GuidelineExample{
						Title: "Code Example",
					}
				}
			}
			continue
		}
		
		if inCodeBlock {
			currentContent.WriteString(line)
			currentContent.WriteString("\n")
			continue
		}
		
		// Parse list items as rules
		if matches := p.listRegex.FindStringSubmatch(line); len(matches) > 1 {
			rule := strings.TrimSpace(matches[1])
			guideline.Rules = append(guideline.Rules, rule)
			continue
		}
		
		// Parse example headers
		if strings.Contains(strings.ToLower(line), "example") ||
		   strings.Contains(strings.ToLower(line), "good") ||
		   strings.Contains(strings.ToLower(line), "bad") {
			if currentExample != nil && (currentExample.Good != "" || currentExample.Bad != "") {
				guideline.Examples = append(guideline.Examples, *currentExample)
			}
			currentExample = &types.GuidelineExample{
				Title: line,
			}
			continue
		}
		
		// Accumulate description
		if line != "" && guideline.Description == "" {
			guideline.Description = line
		} else if line != "" && currentExample != nil && currentExample.Explanation == "" {
			currentExample.Explanation = line
		}
	}
	
	// Add final example
	if currentExample != nil && (currentExample.Good != "" || currentExample.Bad != "") {
		guideline.Examples = append(guideline.Examples, *currentExample)
	}
	
	// Set description from title if empty
	if guideline.Description == "" && guideline.Title != "" {
		guideline.Description = guideline.Title
	}
	
	return guideline
}

// inferCategory infers the category from the title
func (p *Parser) inferCategory(title string) string {
	titleLower := strings.ToLower(title)
	
	if strings.Contains(titleLower, "type") || strings.Contains(titleLower, "interface") {
		return "typing"
	}
	if strings.Contains(titleLower, "naming") || strings.Contains(titleLower, "convention") {
		return "naming"
	}
	if strings.Contains(titleLower, "function") || strings.Contains(titleLower, "method") {
		return "functions"
	}
	if strings.Contains(titleLower, "import") || strings.Contains(titleLower, "export") {
		return "modules"
	}
	if strings.Contains(titleLower, "error") || strings.Contains(titleLower, "exception") {
		return "error_handling"
	}
	if strings.Contains(titleLower, "async") || strings.Contains(titleLower, "promise") {
		return "async"
	}
	
	return "general"
}

// inferPriority infers the priority from the title
func (p *Parser) inferPriority(title string) string {
	titleLower := strings.ToLower(title)
	
	if strings.Contains(titleLower, "must") || strings.Contains(titleLower, "required") ||
	   strings.Contains(titleLower, "critical") || strings.Contains(titleLower, "error") {
		return "high"
	}
	if strings.Contains(titleLower, "should") || strings.Contains(titleLower, "recommend") {
		return "medium"
	}
	if strings.Contains(titleLower, "could") || strings.Contains(titleLower, "consider") ||
	   strings.Contains(titleLower, "style") || strings.Contains(titleLower, "formatting") {
		return "low"
	}
	
	return "medium"
}

// ValidateGuidelines validates that guidelines follow expected format
func (p *Parser) ValidateGuidelines(guidelineSet *types.GuidelineSet) []string {
	var warnings []string
	
	if guidelineSet.Name == "" {
		warnings = append(warnings, "Guideline set name is empty")
	}
	
	if len(guidelineSet.Guidelines) == 0 {
		warnings = append(warnings, "No guidelines found in the set")
	}
	
	for i, guideline := range guidelineSet.Guidelines {
		if guideline.Title == "" {
			warnings = append(warnings, fmt.Sprintf("Guideline %d has no title", i+1))
		}
		if guideline.Description == "" {
			warnings = append(warnings, fmt.Sprintf("Guideline %d (%s) has no description", i+1, guideline.Title))
		}
		if len(guideline.Rules) == 0 && len(guideline.Examples) == 0 {
			warnings = append(warnings, fmt.Sprintf("Guideline %d (%s) has no rules or examples", i+1, guideline.Title))
		}
	}
	
	return warnings
}
package rulechecker

import (
	"strings"

	"github.com/BytecodeAgency/import-boundary-checker/parser"
)

type ImportChart map[string][]string

type Violation struct {
	Filename     string
	CannotImport string
	ImportLine   string
}

type RuleChecker struct {
	Rules       []parser.Rule
	ImportChart ImportChart
	Violations  []Violation
}

func New(rules []parser.Rule, importChart ImportChart) RuleChecker {
	return RuleChecker{
		Rules:       rules,
		ImportChart: importChart,
	}
}

func (rc *RuleChecker) Check() (isValid bool) {
	for file, imports := range rc.ImportChart {
		forbiddenImports := rc.findForbiddenImportsForFilename(file)
		allowedImports := rc.findAllowedImportsForFilename(file)
		for _, importLine := range imports {
			for _, forbiddenImport := range forbiddenImports {
				if strings.HasPrefix(importLine, forbiddenImport) {
					skipViolation := false
					for _, allowedImport := range allowedImports {
						if strings.HasPrefix(importLine, allowedImport) {
							skipViolation = true
						}
					}
					if !skipViolation {
						rc.addViolation(file, importLine, forbiddenImport)
					}
				}
			}
		}
	}
	return len(rc.Violations) == 0
}

func (rc *RuleChecker) findForbiddenImportsForFilename(fileName string) []string {
	var forbiddenImports []string
	for _, rule := range rc.Rules {
		if strings.HasPrefix(fileName, rule.RuleFor) {
			forbiddenImports = append(forbiddenImports, rule.CannotImport...)
		}
	}
	return forbiddenImports
}

func (rc *RuleChecker) findAllowedImportsForFilename(fileName string) []string {
	var allowedImports []string
	for _, rule := range rc.Rules {
		if strings.HasPrefix(fileName, rule.RuleFor) {
			allowedImports = append(allowedImports, rule.AllowImportExceptions...)
		}
	}
	return allowedImports
}

func (rc *RuleChecker) addViolation(filename string, importLine string, cannotImport string) {
	violation := Violation{
		Filename:     filename,
		CannotImport: cannotImport,
		ImportLine:   importLine,
	}
	rc.Violations = append(rc.Violations, violation)
}

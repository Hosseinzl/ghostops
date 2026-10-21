package incident

import (
	"fmt"
	"ghostops/analyzer"
	"ghostops/jsonutil"
	"regexp"
)

type IncidentRule struct {
	Name           string                   `json:"name"`
	Pattern        string                   `json:"pattern"`
	SeverityLevels []analyzer.SeverityLevel `json:"severity_levels"`
	regex          *regexp.Regexp
}

func ReadIncidentRulesFromFile(filePath string) ([]IncidentRule, error) {

	rules, err := jsonutil.LoadFromJSONFile[[]IncidentRule](filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load incident rules from file: %w", err)
	}

	if len(rules) == 0 {
		return nil, fmt.Errorf("no incident rules found in file: %s", filePath)
	}

	for i := range rules {
		re, err := regexp.Compile(rules[i].Pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid regex %q in rule %q: %w", rules[i].Pattern, rules[i].Name, err)
		}
		rules[i].regex = re
	}

	return rules, nil
}

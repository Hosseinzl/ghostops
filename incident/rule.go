package incident

import (
	"fmt"
	"ghostops/analyzer"
	"ghostops/jsonutil"
)

type IncidentRule struct {
	Name           string                  `json:"name"`
	Pattern        string                  `json:"pattern"`
	SeverityLevels []analyzer.SeverityLevel `json:"severity_levels"`
}

func ReadIncidentRulesFromFile(filePath string) ([]IncidentRule, error) {

	if filePath == "" {
		filePath = "incident/incident_rules.json"
	}

	rules, err := jsonutil.LoadFromJSONFile[[]IncidentRule](filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load incident rules from file: %w", err)
	}

	return rules, nil
}

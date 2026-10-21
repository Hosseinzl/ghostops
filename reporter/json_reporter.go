package reporter

import (
	"fmt"
	"ghostops/analyzer"
	"ghostops/jsonutil"
	"ghostops/incident"
)

type jsonReportFormat struct {
	Incidents      []*incident.Incident     `json:"incidents"`
	AnalysisResult *analyzer.AnalysisManager `json:"overall_analysis"`
}

type JsonReporter struct {
	filePath string
}

func NewJsonReporter(filePath string) *JsonReporter {
	return &JsonReporter{
		filePath: filePath,
	}
}

func (jr *JsonReporter) WriteReport(analysisOverall *analyzer.AnalysisManager, incidents []*incident.Incident) error {

	data := jsonReportFormat{
		Incidents:      incidents,
		AnalysisResult: analysisOverall,
	}

	err := jsonutil.SaveToJSONFile(jr.filePath, data)
	if err != nil {
		return fmt.Errorf("failed to save report in JSON: %w", err)
	}

	return nil
}

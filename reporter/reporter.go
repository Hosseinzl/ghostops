package reporter

import (
	"ghostops/analyzer"
	"ghostops/incident"
)

type Reporter interface {
	WriteReport(analysisOverall *analyzer.AnalysisManager, incidents []*incident.Incident) error
}

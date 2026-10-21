package analyzer

type AnalysisManager struct {
	TotalValidEntries     int                   `json:"total_valid_entries"`
	TotalMalformedEntries int                   `json:"total_malformed_entries"`
	EntriesBySeverity     map[SeverityLevel]int `json:"entries_by_severity"`
	EntriesByService      map[string]int        `json:"entries_by_service"`
}

func NewAnalysisManager() *AnalysisManager {
	return &AnalysisManager{
		TotalValidEntries:     0,
		TotalMalformedEntries: 0,
		EntriesBySeverity:     make(map[SeverityLevel]int),
		EntriesByService:      make(map[string]int),
	}
}

func (ar *AnalysisManager) addValidEntries(entry *LogEntry) {
	ar.TotalValidEntries++
	ar.EntriesByService[entry.Service]++
	ar.EntriesBySeverity[entry.SeverityLevel]++
}

func (ar *AnalysisManager) addInvalidEntries() {
	ar.TotalMalformedEntries++
}

func (ar *AnalysisManager) AnalyzeLogEntry(logEntry *LogEntry) {
	if logEntry == nil {
		ar.addInvalidEntries()
		return
	}
	ar.addValidEntries(logEntry)
}

package incident

import (
	"fmt"
	"ghostops/analyzer"
	"sort"
	"time"
)

type Incident struct {
	ID        uint32                   `json:"id"`
	RuleName  string                   `json:"rule_name"`
	Service   string                   `json:"service"`
	Severity  []analyzer.SeverityLevel `json:"severities"`
	Count     int                      `json:"count"`
	StartTime time.Time                `json:"start_time"`
	EndTime   time.Time                `json:"end_time"`
}

type IncidentManager struct {
	incidents      map[string]*Incident
	rules          []IncidentRule
	lastIncidentID uint32
}

func NewIncidentManager(incidentRulesPath string) (*IncidentManager, error) {

	rules, err := ReadIncidentRulesFromFile(incidentRulesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize incident manager: %w", err)
	}

	return &IncidentManager{
		incidents:      make(map[string]*Incident),
		rules:          rules,
		lastIncidentID: 0,
	}, nil
}

func (im *IncidentManager) NextIncidentID() uint32 {
	im.lastIncidentID++
	return im.lastIncidentID
}

func (im *IncidentManager) AppendIncidents(entry *analyzer.LogEntry, key string, rule IncidentRule) {
	im.incidents[key] = &Incident{
		ID:        im.NextIncidentID(),
		RuleName:  rule.Name,
		Service:   entry.Service,
		Severity:  []analyzer.SeverityLevel{entry.SeverityLevel},
		Count:     1,
		StartTime: entry.TimeStamp,
		EndTime:   entry.TimeStamp,
	}

}

func (im *IncidentManager) UpdateIncident(entry *analyzer.LogEntry, key string) {
	im.incidents[key].Count++
	if entry.TimeStamp.After(im.incidents[key].EndTime) {
		im.incidents[key].EndTime = entry.TimeStamp
	}
	if entry.TimeStamp.Before(im.incidents[key].StartTime) {
		im.incidents[key].StartTime = entry.TimeStamp
	}
	if !im.hasSeverity(im.incidents[key].Severity, entry.SeverityLevel) {
		im.incidents[key].Severity = append(im.incidents[key].Severity, entry.SeverityLevel)
	}
}

func (im *IncidentManager) ProcessLogEntry(entry *analyzer.LogEntry) {
	for _, rule := range im.rules {
		if !im.hasSeverity(rule.SeverityLevels, entry.SeverityLevel) {
			continue
		}

		if rule.regex != nil && !rule.regex.MatchString(entry.Message) {
			continue
		}

		key := fmt.Sprintf("%s-%s", entry.Service, rule.Name)
		if _, exists := im.incidents[key]; exists {
			// update the existing incident
			im.UpdateIncident(entry, key)
		} else {
			// create a new incident
			im.AppendIncidents(entry, key, rule)
		}

		break
	}
}

func (im *IncidentManager) hasSeverity(levels []analyzer.SeverityLevel, target analyzer.SeverityLevel) bool {
	for _, level := range levels {
		if level == target {
			return true
		}
	}
	return false
}

func (im *IncidentManager) IncidentsList(sorted bool) []*Incident {
	incidents := make([]*Incident, 0, len(im.incidents))
	for _, incident := range im.incidents {
		incidents = append(incidents, incident)
	}
	if sorted {
		sort.Slice(incidents, func(i, j int) bool {
			if incidents[i].Count != incidents[j].Count {
				return incidents[i].Count > incidents[j].Count
			}
			if !incidents[i].StartTime.Equal(incidents[j].StartTime) {
				return incidents[i].StartTime.Before(incidents[j].StartTime)
			}
			return incidents[i].RuleName < incidents[j].RuleName
		})
	}
	return incidents
}

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
	Incidents      map[string]*Incident
	Rules          []IncidentRule
	LastIncidentID uint32
}

func NewIncidentManager(incidentRulesPath string) *IncidentManager {

	rules, err := ReadIncidentRulesFromFile(incidentRulesPath)
	if err != nil {
		panic(err)
	}

	return &IncidentManager{
		Incidents:      make(map[string]*Incident),
		Rules:          rules,
		LastIncidentID: 0,
	}
}

func (im *IncidentManager) NextIncidentID() uint32 {
	im.LastIncidentID++
	return im.LastIncidentID
}

func (im *IncidentManager) AppendIncidents(entry *analyzer.LogEntry, key string, rule IncidentRule) {
	im.Incidents[key] = &Incident{
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
	im.Incidents[key].Count++
	if entry.TimeStamp.After(im.Incidents[key].EndTime) {
		im.Incidents[key].EndTime = entry.TimeStamp
	}
	if entry.TimeStamp.Before(im.Incidents[key].StartTime) {
		im.Incidents[key].StartTime = entry.TimeStamp
	}
	if !im.hasSeverity(im.Incidents[key].Severity, entry.SeverityLevel) {
		im.Incidents[key].Severity = append(im.Incidents[key].Severity, entry.SeverityLevel)
	}
}

func (im *IncidentManager) ProcessLogEntry(entry *analyzer.LogEntry) {
	for _, rule := range im.Rules {
		if !im.hasSeverity(rule.SeverityLevels, entry.SeverityLevel) {
			continue
		}

		if rule.regex != nil && rule.regex.MatchString(entry.Message) {
			continue
		}

		key := fmt.Sprintf("%s-%s", entry.Service, rule.Name)
		if _, exists := im.Incidents[key]; exists {
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
	incidents := make([]*Incident, 0, len(im.Incidents))
	for _, incident := range im.Incidents {
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

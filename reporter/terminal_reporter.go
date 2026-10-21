package reporter

import (
	"fmt"
	"ghostops/analyzer"
	"ghostops/incident"
	"io"
	"sort"
	"strings"
)

type TerminalReporter struct {
	out io.Writer
}

func NewTerminalReporter(out io.Writer) *TerminalReporter {
	return &TerminalReporter{
		out: out,
	}
}

func (tr *TerminalReporter) WriteReport(analysisOverall *analyzer.AnalysisManager, incidents []*incident.Incident) error {
	printOveralAnalysis(analysisOverall, tr.out)
	printTopPriorityIncidents(incidents, tr.out)
	return nil
}

func printOveralAnalysis(result *analyzer.AnalysisManager, out io.Writer) {

	fmt.Fprintf(out, "Overall Analysis Result:\n")
	fmt.Fprintf(out, "Total Valid Entries: %d\n", result.TotalValidEntries)
	fmt.Fprintf(out, "Total Malformed Entries: %d\n", result.TotalMalformedEntries)

	// Sort the severity levels for consistent output
	fmt.Fprintf(out, "Entries by Severity:\n")

	severityLevels := make([]analyzer.SeverityLevel, 0, len(result.EntriesBySeverity))
	for severity := range result.EntriesBySeverity {
		severityLevels = append(severityLevels, severity)
	}
	sort.Slice(severityLevels, func(i, j int) bool {
		return severityLevels[i] < severityLevels[j]
	})

	for _, severity := range severityLevels {
		count := result.EntriesBySeverity[severity]
		fmt.Fprintf(out, "  %v: %d\n", severity.String(), count)
	}

	// Sort the services for consistent output
	fmt.Fprintf(out, "Entries by Service:\n")

	services := make([]string, 0, len(result.EntriesByService))
	for service := range result.EntriesByService {
		services = append(services, service)
	}

	sort.Strings(services)

	for _, service := range services {
		count := result.EntriesByService[service]
		fmt.Fprintf(out, "  %s: %d\n", service, count)
	}
}

func printTopPriorityIncidents(incidents []*incident.Incident, out io.Writer) {

	fmt.Fprintf(out, "\nDetected Incidents (%d):\n", len(incidents))

	if len(incidents) == 0 {
		fmt.Fprintf(out, "  No incidents detected.\n")
		return
	}

	for _, incident := range incidents {

		severityStrings := make([]string, len(incident.Severity))
		for i, severity := range incident.Severity {
			severityStrings[i] = severity.String()
		}

		fmt.Fprintf(out, "-----------------------------\n")
		fmt.Fprintf(out, "Incident ID: %d\n", incident.ID)
		fmt.Fprintf(out, "Rule Name: %s\n", incident.RuleName)
		fmt.Fprintf(out, "Service: %s\n", incident.Service)
		fmt.Fprintf(out, "Severities: %s\n", strings.Join(severityStrings, ", "))
		fmt.Fprintf(out, "Count: %d\n", incident.Count)
		fmt.Fprintf(out, "Start Time: %s\n", incident.StartTime.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(out, "End Time: %s\n", incident.EndTime.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(out, "-----------------------------\n")
	}
}

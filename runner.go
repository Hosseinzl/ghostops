package main

import (
	"fmt"
	"ghostops/analyzer"
	"ghostops/incident"
)

type Streamer interface {
	EachLine(fn func(line string) error) error
}

type Reporter interface {
	WriteReport(analysisOverall *analyzer.AnalysisManager, incidents []*incident.Incident) error
}

type Runner struct {
	streamer        Streamer
	analyzeTool     *analyzer.AnalysisManager
	incidentManager *incident.IncidentManager
	reporters       []Reporter
	flagConfig      *FlagConfig
}

func NewRunner(streamer Streamer,
	analyzeTool *analyzer.AnalysisManager,
	incidentManager *incident.IncidentManager,
	reporters []Reporter,
	flagConfig *FlagConfig) *Runner {

	return &Runner{
		streamer:        streamer,
		analyzeTool:     analyzeTool,
		incidentManager: incidentManager,
		reporters:       reporters,
		flagConfig:      flagConfig,
	}
}

func (r *Runner) Run() error {
	err := r.streamer.EachLine(func(line string) error {
		logEntry, err := analyzer.VerifyLogEntry(line)

		if err != nil {
			r.analyzeTool.RecordMalformed()
			return nil
		}

		if filterErr := filterBaseFlags(logEntry, r.flagConfig); filterErr != nil {
			return nil
		}

		r.analyzeTool.RecordValid(logEntry)

		r.incidentManager.ProcessLogEntry(logEntry)

		return nil
	})

	if err != nil {
		return err
	}

	err = r.GenerateReport()

	return err

}

func filterBaseFlags(logEntry *analyzer.LogEntry, config *FlagConfig) error {

	if !config.after.IsZero() && logEntry.TimeStamp.Before(config.after) {
		return fmt.Errorf("log entry is before the specified --after time")
	}

	if !config.before.IsZero() && logEntry.TimeStamp.After(config.before) {
		return fmt.Errorf("log entry is after the specified --before time")
	}

	if config.service != "" && logEntry.Service != config.service {
		return fmt.Errorf("log entry is not from the specified --service")
	}

	if config.minSeverityLevel != 0 && logEntry.SeverityLevel < config.minSeverityLevel {
		return fmt.Errorf("log entry severity level is below the specified --min-level")
	}

	return nil

}

func (r *Runner) GenerateReport() error {
	for _, reporter := range r.reporters {
		err := reporter.WriteReport(r.analyzeTool, r.incidentManager.IncidentsList())
		if err != nil {
			return fmt.Errorf("error in generate report: %w", err)
		}
	}

	return nil
}

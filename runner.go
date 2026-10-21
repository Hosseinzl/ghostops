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
}

func NewRunner(streamer Streamer,
	analyzeTool *analyzer.AnalysisManager,
	incidentManager *incident.IncidentManager,
	reporters []Reporter) *Runner {

	return &Runner{
		streamer:        streamer,
		analyzeTool:     analyzeTool,
		incidentManager: incidentManager,
		reporters:       reporters,
	}
}

func (r *Runner) Run() error {
	err := r.streamer.EachLine(func(line string) error {
		logEntry, err := analyzer.VerifyLogEntry(line)
		if err != nil {
			r.analyzeTool.RecordMalformed()
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

func (r *Runner) GenerateReport() error {
	for _, reporter := range r.reporters {
		err := reporter.WriteReport(r.analyzeTool, r.incidentManager.IncidentsList())
		if err != nil {
			return fmt.Errorf("error in generate report: %w", err)
		}
	}

	return nil
}

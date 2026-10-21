package main

import (
	"fmt"
	"ghostops/analyzer"
	"ghostops/incident"
	"ghostops/reporter"
	"ghostops/streamer"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {

	flagConfiguration, err := NewFlagConfig()
	if err != nil {
		return err
	}

	file, err := os.Open(flagConfiguration.logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	logStreamer, err := streamer.NewLogStreamer(file)
	if err != nil {
		return fmt.Errorf("failed to create log streamer: %w", err)
	}

	incidentManager, err := incident.NewIncidentManager(flagConfiguration.incidentRulesPath)
	if err != nil {
		return fmt.Errorf("failed to create incident manager: %w", err)
	}

	runner := NewRunner(
		logStreamer,
		analyzer.NewAnalysisManager(),
		incidentManager,
		initReporters(flagConfiguration.outputPath),
	)

	return runner.Run()
}

func initReporters(outputPath string) []Reporter {
	reporters := []Reporter{
		reporter.NewTerminalReporter(os.Stdout),
	}
	if outputPath != "" {
		reporters = append(reporters, reporter.NewJsonReporter(outputPath))
	}
	return reporters
}

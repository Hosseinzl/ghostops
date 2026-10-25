package main

import (
	"flag"
	"fmt"
	"ghostops/analyzer"
	"time"
)

type FlagConfig struct {
	logPath           string
	outputPath        string
	incidentRulesPath string
	after             time.Time
	before            time.Time
	service           string
	minSeverityLevel  analyzer.SeverityLevel
}

func NewFlagConfig() (*FlagConfig, error) {

	logPath := flag.String("p", "", "Path to the input log file")
	outputPath := flag.String("o", "", "Path to the output JSON file")
	after := flag.String("after", "", "Show logs after a specified time")
	before := flag.String("before", "", "Show logs before a specified time")
	service := flag.String("service", "", "Show logs of a specified service")
	minLevel := flag.String("min-level", "", "Show logs with a minimum log level")

	incidentRulesPath := flag.String("r", "incident/incident_rules.json", "Path to the incident rules file")

	flag.Parse()

	// Validate required flags.
	if *logPath == "" {
		return nil, fmt.Errorf("-p is required. please provide log file path")
	}

	// Validate positional arguments.
	if flag.NArg() > 0 {
		return nil, fmt.Errorf("unexpected positional arguments")
	}

	// Validate time flags.
	var afterTime, beforeTime time.Time
	var err error

	if *after != "" {
		afterTime, err = parseTime(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid --after value: %w", err)
		}
	}

	if *before != "" {
		beforeTime, err = parseTime(*before)
		if err != nil {
			return nil, fmt.Errorf("invalid --before value: %w", err)
		}
	}	

	// Validate severity level flag.
	var minSeverityLevel analyzer.SeverityLevel
	if *minLevel != "" {
		minSeverityLevel, err = analyzer.VerifySeverityLevel(*minLevel)
		if err != nil {
			return nil, fmt.Errorf("invalid --min-level value: %w", err)
		}
	}
	return &FlagConfig{
		logPath:           *logPath,
		outputPath:        *outputPath,
		after:             afterTime,
		before:            beforeTime,
		service:           *service,
		minSeverityLevel:  minSeverityLevel,
		incidentRulesPath: *incidentRulesPath,
	}, nil

}

func parseTime(value string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, value)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"invalid time %q: expected YYYY-MM-DD or RFC3339 format",
		value,
	)
}
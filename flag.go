package main

import (
	"flag"
	"fmt"
)

type FlagConfig struct {
	logPath           string
	outputPath        string
	incidentRulesPath string
}

func NewFlagConfig() (*FlagConfig, error) {

	logPath := flag.String("p", "", "Path to the input log file")
	outputPath := flag.String("o", "", "Path to the output JSON file")
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

	return &FlagConfig{
		logPath:           *logPath,
		outputPath:        *outputPath,
		incidentRulesPath: *incidentRulesPath,
	}, nil

}

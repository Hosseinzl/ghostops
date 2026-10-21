package analyzer

import (
	"fmt"
	"strings"
	"time"
)

type SeverityLevel uint8

const (
	Unknown SeverityLevel = iota
	Info
	Warn
	Error
	Fatal
)

var severityLevelName = map[string]SeverityLevel{
	"INFO":  Info,
	"WARN":  Warn,
	"ERROR": Error,
	"FATAL": Fatal,
}

func (s SeverityLevel) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *SeverityLevel) UnmarshalText(text []byte) error {
	level, ok := severityLevelName[string(text)]
	if !ok {
		return fmt.Errorf("unknown severity level: %s", string(text))
	}
	*s = level
	return nil
}

func (s SeverityLevel) String() string {
	switch s {
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type LogEntry struct {
	TimeStamp     time.Time
	SeverityLevel SeverityLevel
	Service       string
	Message       string
}

func verifyTimeStamp(timeStampStr string) (time.Time, error) {

	timeStamp, err := time.Parse(time.RFC3339, timeStampStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error in parse timestamp: %w", err)
	}

	return timeStamp, nil

}

func verifySeverityLevel(level string) (SeverityLevel, error) {
	severityLevel := severityLevelName[level]
	if severityLevel == 0 {
		return 0, fmt.Errorf("severity level is invalid")
	}

	return severityLevel, nil
}

func verifyServiceName(serviceName string) (string, error) {
	return serviceName, nil
}

func verifyMessage(message string) (string, error) {
	return message, nil
}

func VerifyLogEntry(entry string) (*LogEntry, error) {
	if entry == "" {
		return nil, fmt.Errorf("log entry is empty")
	}

	entryParts := strings.SplitN(entry, " ", 4)

	if len(entryParts) < 4 {
		return nil, fmt.Errorf("log entry is malformed")
	}

	timeStamp, err := verifyTimeStamp(entryParts[0])
	if err != nil {
		return nil, fmt.Errorf("log entry is malformed: %w", err)
	}

	severityLevel, err := verifySeverityLevel(entryParts[1])
	if err != nil {
		return nil, fmt.Errorf("log entry is malformed: %w", err)
	}

	serviceName, err := verifyServiceName(entryParts[2])
	if err != nil {
		return nil, fmt.Errorf("log entry is malformed: %w", err)
	}

	message, err := verifyMessage(strings.Join(entryParts[3:], " "))
	if err != nil {
		return nil, fmt.Errorf("log entry is malformed: %w", err)
	}

	logEntry := &LogEntry{
		TimeStamp:     timeStamp,
		SeverityLevel: severityLevel,
		Service:       serviceName,
		Message:       message,
	}

	return logEntry, nil
}

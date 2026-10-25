# GhostOps: Streaming Log Analysis CLI

GhostOps is a command-line tool that analyzes application logs, extracts useful statistics, detects predefined incident patterns, and generates a report. It processes structured application logs sequentially via buffered streaming, extracts metrics, detects operational incidents using configurable precompiled regular expressions, and exports deterministic reports in both formatted terminal and JSON outputs.

---


## Core Features

- **Stream-Based Processing:** Processes large scale log files sequentially using `bufio.Scanner` with a dynamically scalable buffer up to 1 MB, keeping resident memory usage exceptionally low.
- **Precompiled RegEx Matcher:** Precompiles all pattern matching rules during initialization to ensure maximum CPU throughput per line.
- **Incident Clustering & Aggregation:** Groups related events by service and rule type, tracking frequency count and timestamps 
- **Flexible Log Filtering:** Filter by service name, minimum severity level, and historical time ranges.
- **Deterministic Reporting:** Guarantees consistent sorting across map iterations for both terminal stdout and structured JSON files.

---
## Architecture & Data Pipeline

GhostOps follows a processing pipeline designed for memory-efficient log analysis and clear separation of responsibilities.

### Processing Pipeline

1. **Log Streaming (`streamer.LogStreamer`)**

   Reads the input log file incrementally using buffered streaming. This approach avoids loading the entire file into memory and supports processing large log files efficiently.

2. **Entry Validation (`analyzer.VerifyEntry`)**

   Each raw log line is passed to the analyzer for validation. The entry is checked against the expected schema, including RFC3339 timestamp parsing. Valid entries are returned as `LogEntry` objects, while malformed lines are handled separately.

3. **Malformed Entry Handling**

   Lines that fail validation are recorded through `analyzer.RecordMalformed()`. This allows GhostOps to track malformed entries without interrupting the processing of subsequent lines.

4. **Base Filtering (`runner.filterBaseFlags`)**

   Valid entries are checked against the configured filters, including:

   - Minimum severity level (`--min-level`)
   - Service name (`--service`)
   - Start time (`--after`)
   - End time (`--before`)

   Entries that do not match the filters are excluded from further analysis.

5. **Analysis and Incident Processing**

   Entries that pass all filters are processed by two components:

   - `analyzer.RecordValid()` updates the analysis statistics, including severity and service counts.
   - `incident.ProcessLogEntry()` evaluates entries against the incident detection rules and updates the incident results.

6. **Orchestration (`runner.Runner`)**

   The `Runner` coordinates the processing pipeline, connecting the streamer, analyzer, filtering logic, and incident processor. It also manages the overall execution flow and prepares the results for reporting.

7. **Report Generation**

   After processing is complete, GhostOps generates reports through the appropriate reporter implementations:

   - **`reporter.TerminalReporter`** produces a human-readable terminal report with formatted tables and an overview of the analysis results.
   - **`reporter.JsonReporter`** exports the results as indented JSON for further processing or integration with other tools.

### Design Principles

- **Memory Efficiency:** Log entries are processed incrementally rather than loading the entire file into memory.
- **Separation of Concerns:** Streaming, validation, filtering, analysis, incident detection, and reporting have distinct responsibilities.
- **Unidirectional Data Flow:** Log entries move through a defined sequence of processing stages.
- **Extensibility:** Separate components make it easier to modify filtering rules, incident detection, or reporting without tightly coupling the entire pipeline.
- **Error Resilience:** Malformed entries are tracked separately so that individual invalid lines do not prevent the analysis of the remaining file.


---

## Getting Started

### Building from Source

Clone the repository and build the binary:

```bash
# Clone the repository
git clone https://github.com/Hosseinzl/ghostops.git
cd ghostops

# Run GhostOps
go run . -p <log-file-path> [options]
```

### CLI Flags & Options

| Flag | Required | Default | Description |
| :--- | :--- | :--- | :--- |
| `-p` | **Yes** | — | Path to the target input log file |
| `-o` | No | `""` | Output path to export the report in JSON format |
| `-r` | No | `incident/incident_rules.json` | Path to custom incident detection rules JSON file |
| `--service` |  No | `""` | Filter entries by a specific service name |
| `--min-level` | No | `""` | Filter entries with a minimum severity (`INFO`, `WARN`, `ERROR`, `FATAL`) |
| `--after` | No | `""` | Filter entries after timestamp (`YYYY-MM-DD` or `RFC3339`) |
| `--before` | No | `""` | Filter entries before timestamp (`YYYY-MM-DD` or `RFC3339`) |

---

## Filtering Capabilities (Extensions)

GhostOps supports analytical filtering flags that restrict analysis without violating the streaming paradigm:

```bash
# 1. Basic log analysis printed to terminal
./ghostops -p sample.log

# 2. Analyze logs and export to JSON
./ghostops -p sample.log -o report.json

# 3. Filter by microservice and severity threshold
./ghostops -p sample.log --service auth-service --min-level WARN

# 4. Filter within a historical time window
./ghostops -p sample.log --after 2026-10-01T08:00:00Z --before 2026-10-01T09:30:00Z

# 5. Combined analysis with custom rules and JSON output
./ghostops -p sample.log -r incident/incident_rules.json --service db-service --min-level ERROR -o db_incidents.json
```

---

## Incident Detection Engine

Incident rules are defined via external JSON configurations. Each rule specifies a pattern and a slice of target severity levels:

```json
[
  {
    "name": "Database Connection Failure",
    "pattern": "Failed to connect to database(:\\s*timeout after \\d+ms)?",
    "severity_levels": [
      "ERROR",
      "FATAL"
    ]
  },
  {
    "name": "Rate Limit Warning",
    "pattern": "Rate limit approaching( for client IP [\\d\\.]+)?",
    "severity_levels": [
      "WARN"
    ]
  }
]
```

---

## Key Engineering Decisions

### 1. Memory-Efficient Streaming with Controlled Buffers
Rather than loading large files into memory with `os.ReadFile` (which causes Out-Of-Memory exceptions on large inputs), `LogStreamer` wraps an `io.Reader` using `bufio.Scanner`.  
the buffer starts at **64 KB** and dynamically scales up to a safety ceiling of **1 MB**:
```go
maxBufferSize := 1024 * 1024 // 1 MB
scanner.Buffer(make([]byte, 64*1024), maxBufferSize)
```

### 2. Precompiled Regular Expressions
Invoking `regexp.MatchString(pattern, message)` on every line forces the regex engine to parse and compile the state machine repeatedly, causing severe CPU spikes. GhostOps precompiles all patterns **once** upon reading the configuration:
```go
for i := range rules {
    re, err := regexp.Compile(rules[i].Pattern)
    if err != nil {
        return nil, fmt.Errorf("invalid regex %q in rule %q: %w", rules[i].Pattern, rules[i].Name, err)
    }
    rules[i].regex = re
}
```
During line evaluation, the precompiled `*regexp.Regexp` is matched directly.

### 3. Consumer-Driven Interfaces
According to *"Consumer owns the interface"*, abstractions (`Streamer`, `Reporter`) are declared directly inside `runner.go` where they are consumed, rather than inside provider packages:
- Providers (`LogStreamer`, `TerminalReporter`, `JsonReporter`) return concrete structs.
- `Runner` defines the exact interface surface it requires, keeping coupling minimal and making testing straightforward.

---
## Sample Outputs

### Terminal Output

```text
Overall Analysis Result:
Total Valid Entries: 130
Total Malformed Entries: 468
Entries by Severity:
  ERROR: 97
  FATAL: 33
Entries by Service:
  payment-service: 130

Detected Incidents (2):
-----------------------------
Incident ID: 1
Rule Name: Payment Gateway Issue
Service: payment-service
Severities: ERROR
Count: 97
Start Time: 2026-10-03 03:23:02
End Time: 2026-10-06 17:19:29
-----------------------------
-----------------------------
Incident ID: 2
Rule Name: Database Pool Exhaustion
Service: payment-service
Severities: FATAL
Count: 33
Start Time: 2026-10-03 01:07:14
End Time: 2026-10-06 14:58:35
-----------------------------
```

### JSON Export

```json
{
  "incidents": [
    {
      "id": 21,
      "rule_name": "Payment Gateway Issue",
      "service": "auth-service",
      "severities": [
        "ERROR"
      ],
      "count": 30,
      "start_time": "2026-10-01T02:15:44Z",
      "end_time": "2026-10-06T15:42:34Z"
    },
    {
      "id": 12,
      "rule_name": "Payment Gateway Issue",
      "service": "notification-service",
      "severities": [
        "ERROR"
      ],
      "count": 26,
      "start_time": "2026-10-01T08:49:29Z",
      "end_time": "2026-10-06T17:59:02Z"
    },
    ...
  ],
  "overall_analysis": {
    "total_valid_entries": 5532,
    "total_malformed_entries": 468,
    "entries_by_severity": {
      "ERROR": 820,
      "FATAL": 271,
      "INFO": 3335,
      "WARN": 1106
    },
    "entries_by_service": {
      "auth-service": 1106,
      "notification-service": 1130,
      "order-service": 1082,
      "payment-service": 1090,
      "user-service": 1124
    }
  }
}

```
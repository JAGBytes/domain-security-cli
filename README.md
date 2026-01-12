# Domain Security CLI

A command-line tool to analyze SSL/TLS security of domains using the SSL Labs API.

## Features

- Analyze SSL/TLS configuration of any domain
- Display security grades for each endpoint
- Output in table or JSON format
- Automatic polling for in-progress analyses
- Rate limit handling

## Installation

Requires Go 1.21 or higher.

Clone the repository:

```bash
git clone https://github.com/JAGBytes/domain-security-cli.git
```

Build the executable:

```bash
cd domain-security-cli
go build -o domain-security-cli.exe
./domain-security-cli analyze <domain>
```

Or run directly:

```bash
go run main.go analyze <domain>
```

## Usage

```bash
# Analyze a domain
go run main.go analyze amazon.com

# JSON output
go run main.go analyze amazon.com --json
go run main.go analyze amazon.com --json > results.json

# Force new analysis, takes 1-2 minutes
go run main.go analyze amazon.com --new

# Analyze a domain with executable
./domain-security-cli.exe analyze github.com

# using executable JSON output
./domain-security-cli.exe analyze github.com --json
go run main.go analyze github.com --json > results.json

# Force new analysis with built executable
./domain-security-cli.exe analyze github.com --new



```

## Output

Table format shows the overall grade and details for each endpoint:

```
SSL Labs Analysis Results
==================================================

Domain: google.com
Overall Grade: B
Status: READY

ENDPOINTS: 10
--------------------------------------------------

[1] IP: 2607:f8b0:4002:c0c:0:0:0:71
Grade: B
Server: yi-in-f113.1e100.net
Status: Ready

...
```

JSON format provides complete data for all analyzed endpoints.

## Project structure

- `cmd/` - Command-line interface using Cobra
- `internal/client/` - Communicates with SSL Labs API
- `internal/service/` - Handles analysis logic and polling
- `pkg/formatter/` - Formats output as table or JSON

## Requirements

- Go 1.21 or higher
- Internet connection to reach SSL Labs API

## License

MIT

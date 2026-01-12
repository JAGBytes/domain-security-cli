package formatter

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/JAGBytes/domain-security-cli/internal/models"
)

func FormatAsTable(result *models.CheckResult) string {
	output := "SSL Labs Analysis Results\n"
	output += strings.Repeat("=", 50) + "\n\n"

	output += fmt.Sprintf("Domain: %s\n", result.Domain)
	output += fmt.Sprintf("Overall Grade: %s\n", result.Grade)
	output += fmt.Sprintf("Status: %s\n\n", result.Status)

	output += fmt.Sprintf("ENDPOINTS: %d\n", len(result.Endpoints))
	output += strings.Repeat("-", 50) + "\n\n"

	if len(result.Endpoints) == 0 {
		output += "No endpoints found.\n"
	} else {
		for i, ep := range result.Endpoints {
			output += fmt.Sprintf("[%d] IP: %s\n", i+1, ep.IPAddress)
			output += fmt.Sprintf("Grade: %s\n", ep.Grade)

			if ep.ServerName != "" {
				output += fmt.Sprintf("Server: %s\n", ep.ServerName)
			}

			if ep.HasWarnings {
				output += "Warnings: Yes\n"
			}

			if ep.StatusMessage != "" {
				output += fmt.Sprintf("Status: %s\n", ep.StatusMessage)
			}
			output += "\n\n"
		}
	}

	return output
}

func FormatAsJSON(result *models.CheckResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

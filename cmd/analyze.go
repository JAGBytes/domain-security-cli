package cmd

import (
	"fmt"
	"os"

	"github.com/JAGBytes/domain-security-cli/internal/service"
	"github.com/JAGBytes/domain-security-cli/pkg/formatter"
	"github.com/spf13/cobra"
)

var (
	jsonOutput  bool
	newAnalysis bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [domain]",
	Short: "Analyze SSL/TLS security of a domain",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		checker := service.NewCheckerService()
		results, errors := checker.CheckDomains(args, newAnalysis)

		for _, err := range errors {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

		if jsonOutput {
			output, err := formatter.FormatAsJSON(results)
			if err != nil {
				return err
			}
			fmt.Println(output)
		} else {
			fmt.Println(formatter.FormatAsTables(results))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	analyzeCmd.Flags().BoolVar(&newAnalysis, "new", false, "Force new analysis")
}

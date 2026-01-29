package cmd

import (
	"fmt"

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
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		domain := args[0]

		checker := service.NewCheckerService()
		result, err := checker.CacheDomain(domain, newAnalysis)
		if err != nil {
			return fmt.Errorf("analysis failed: %w", err)
		}

		if jsonOutput {
			output, err := formatter.FormatAsJSON(result)
			if err != nil {
				return err
			}
			fmt.Println(output)
		} else {
			fmt.Println(formatter.FormatAsTable(result))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	analyzeCmd.Flags().BoolVar(&newAnalysis, "new", false, "Force new analysis")
}

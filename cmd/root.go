package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	baseFile   string
	targetFile string
	outputFmt  string
	outputFile string
)

var rootCmd = &cobra.Command{
	Use:   "driftctl-diff",
	Short: "Compare Terraform state files and surface configuration drift",
	Long: `driftctl-diff compares two Terraform state files and produces a
readable report highlighting added, removed, and changed resources.`,
	RunE: runDiff,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&baseFile, "base", "b", "", "Path to the base Terraform state file (required)")
	rootCmd.Flags().StringVarP(&targetFile, "target", "t", "", "Path to the target Terraform state file (required)")
	rootCmd.Flags().StringVarP(&outputFmt, "format", "f", "text", "Output format: text or json")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write report to file instead of stdout")

	_ = rootCmd.MarkFlagRequired("base")
	_ = rootCmd.MarkFlagRequired("target")
}

package cmd

import (
	"fmt"
	"os"

	"github.com/owner/driftctl-diff/internal/diff"
	"github.com/owner/driftctl-diff/internal/lint"
	"github.com/owner/driftctl-diff/internal/state"
	"github.com/spf13/cobra"
)

var lintConfigPath string

func init() {
	lintCmd := &cobra.Command{
		Use:   "lint <source.tfstate> <target.tfstate>",
		Short: "Apply lint rules to detected drift",
		Args:  cobra.ExactArgs(2),
		RunE:  runLint,
	}
	lintCmd.Flags().StringVar(&lintConfigPath, "lint-config", "", "path to lint config YAML file")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	src, err := state.ParseStateFile(args[0])
	if err != nil {
		return fmt.Errorf("lint: parse source: %w", err)
	}
	tgt, err := state.ParseStateFile(args[1])
	if err != nil {
		return fmt.Errorf("lint: parse target: %w", err)
	}

	srcIdx := state.IndexResources(src)
	tgtIdx := state.IndexResources(tgt)
	result := diff.Compare(srcIdx, tgtIdx)

	var cfg lint.Config
	if lintConfigPath != "" {
		cfg, err = lint.LoadConfig(lintConfigPath)
		if err != nil {
			return fmt.Errorf("lint: load config: %w", err)
		}
	} else {
		cfg = lint.DefaultLintConfig()
	}

	rules := lint.FilterRules(lint.DefaultRules(), cfg)
	findings := lint.Run(result, rules)

	if len(findings) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "✓ No lint findings.")
		return nil
	}

	for _, f := range findings {
		fmt.Fprintln(cmd.OutOrStdout(), f.String())
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\n%d finding(s) reported.\n", len(findings))

	if lint.ShouldFail(findings, cfg) {
		os.Exit(1)
	}
	return nil
}

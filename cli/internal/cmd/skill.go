package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cogfoundry-ai/loomloom/cli/internal/skill"
	"github.com/spf13/cobra"
)

func newSkillCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Install LoomLoom templates as local agent skills",
	}
	install := &cobra.Command{
		Use:   "install",
		Short: "Install a LoomLoom template wrapper as a local agent skill",
	}
	install.AddCommand(
		newSkillInstallMarketCmd(opts),
		newSkillInstallTemplateSpecCmd(opts),
	)
	cmd.AddCommand(install)
	return cmd
}

type skillInstallOptions struct {
	agent     string
	outputDir string
	dryRun    bool
}

func addSkillInstallFlags(cmd *cobra.Command, installOpts *skillInstallOptions) {
	cmd.Flags().StringVar(&installOpts.agent, "agent", "", "Target agent: codex|claude|openclaw")
	cmd.Flags().StringVar(&installOpts.outputDir, "output-dir", "", "Output directory for this single generated skill")
	cmd.Flags().BoolVar(&installOpts.dryRun, "dry-run", false, "Validate and preview without writing final skill files")
	_ = cmd.MarkFlagRequired("agent")
	_ = cmd.MarkFlagRequired("output-dir")
}

func newSkillInstallMarketCmd(opts *rootOptions) *cobra.Command {
	installOpts := &skillInstallOptions{}
	cmd := &cobra.Command{
		Use:   "market <listing-id>",
		Short: "Install a Market SkillBot listing as a local agent skill",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSkillInstall(cmd, opts, skill.Options{
				SourceType: skill.SourceMarketListing,
				Agent:      skill.Agent(strings.TrimSpace(installOpts.agent)),
				OutputDir:  strings.TrimSpace(installOpts.outputDir),
				DryRun:     installOpts.dryRun,
				ListingID:  strings.TrimSpace(args[0]),
			})
		},
	}
	addSkillInstallFlags(cmd, installOpts)
	return cmd
}

func newSkillInstallTemplateSpecCmd(opts *rootOptions) *cobra.Command {
	installOpts := &skillInstallOptions{}
	cmd := &cobra.Command{
		Use:   "template-spec <template-id> <version-id>",
		Short: "Install a private TemplateSpec version as a local agent skill",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSkillInstall(cmd, opts, skill.Options{
				SourceType:        skill.SourceUserTemplate,
				Agent:             skill.Agent(strings.TrimSpace(installOpts.agent)),
				OutputDir:         strings.TrimSpace(installOpts.outputDir),
				DryRun:            installOpts.dryRun,
				TemplateID:        strings.TrimSpace(args[0]),
				TemplateVersionID: strings.TrimSpace(args[1]),
			})
		},
	}
	addSkillInstallFlags(cmd, installOpts)
	return cmd
}

func runSkillInstall(cmd *cobra.Command, opts *rootOptions, skillOpts skill.Options) error {
	httpClient, err := newHTTPClient(opts)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
	defer cancel()

	exporter := skill.NewExporter(httpClient)
	result, err := exporter.Install(ctx, skillOpts)
	if err != nil {
		return err
	}
	if opts.output == "json" {
		return writeIndentedJSON(cmd.OutOrStdout(), result)
	}
	if skillOpts.DryRun {
		if !result.Installable {
			return printSkillPreview(cmd.OutOrStdout(), result)
		}
		return printSkillPreview(cmd.OutOrStdout(), result)
	}
	return printSkillInstallResult(cmd.OutOrStdout(), result)
}

func printSkillPreview(w interface {
	Write([]byte) (int, error)
}, result *skill.InstallResult) error {
	if _, err := fmt.Fprintf(w, "installable\t%t\n", result.Installable); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "skill_name\t%s\nsource_type\t%s\nagent\t%s\noutput_dir\t%s\n", result.SkillName, result.SourceType, result.Agent, result.OutputDir); err != nil {
		return err
	}
	if result.SourceID != "" {
		if _, err := fmt.Fprintf(w, "source_id\t%s\n", result.SourceID); err != nil {
			return err
		}
	}
	for _, warning := range result.Warnings {
		if _, err := fmt.Fprintf(w, "warning\t%s\t%s\n", warning.Code, warning.Message); err != nil {
			return err
		}
	}
	for _, issue := range result.Errors {
		if _, err := fmt.Fprintf(w, "error\t%s\t%s\n", issue.Code, issue.Message); err != nil {
			return err
		}
	}
	if result.BlockingReason != "" {
		if _, err := fmt.Fprintf(w, "blocking_reason\t%s\nrecommended_action\t%s\n", result.BlockingReason, result.RecommendedAction); err != nil {
			return err
		}
	}
	return nil
}

func printSkillInstallResult(w interface {
	Write([]byte) (int, error)
}, result *skill.InstallResult) error {
	if _, err := fmt.Fprintf(w, "installed\t%t\nskill_name\t%s\nsource_type\t%s\nagent\t%s\noutput_dir\t%s\n", result.Installed, result.SkillName, result.SourceType, result.Agent, result.OutputDir); err != nil {
		return err
	}
	if result.SourceID != "" {
		if _, err := fmt.Fprintf(w, "source_id\t%s\n", result.SourceID); err != nil {
			return err
		}
	}
	switch result.SourceType {
	case skill.SourceMarketListing:
		if _, err := fmt.Fprintf(w, "listing_id\t%s\nlisting_version_id\t%s\n", result.ListingID, result.ListingVersionID); err != nil {
			return err
		}
	case skill.SourceUserTemplate:
		if _, err := fmt.Fprintf(w, "template_id\t%s\ntemplate_version_id\t%s\n", result.TemplateID, result.TemplateVersionID); err != nil {
			return err
		}
	}
	if result.Metadata != "" {
		if _, err := fmt.Fprintf(w, "metadata\t%s\nskill_file\t%s\n", result.Metadata, result.SkillFile); err != nil {
			return err
		}
	}
	if result.SourceType == skill.SourceMarketListing {
		if _, err := fmt.Fprintln(w, "notice\tMarket executions will use the current listing price, permission, and transaction rules."); err != nil {
			return err
		}
	}
	if result.Trigger != "" {
		if _, err := fmt.Fprintf(w, "trigger_example\t%s\n", result.Trigger); err != nil {
			return err
		}
	}
	for _, warning := range result.Warnings {
		if _, err := fmt.Fprintf(w, "warning\t%s\t%s\n", warning.Code, warning.Message); err != nil {
			return err
		}
	}
	return nil
}

/*
Copyright 2021 The terraform-docs Authors.

Licensed under the MIT license (the "License"); you may not
use this file except in compliance with the License.

You may obtain a copy of the License at the LICENSE file in
the root directory of this source tree.
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/terraform-docs/terraform-docs/cmd/asciidoc"
	"github.com/terraform-docs/terraform-docs/cmd/completion"
	"github.com/terraform-docs/terraform-docs/cmd/json"
	"github.com/terraform-docs/terraform-docs/cmd/markdown"
	"github.com/terraform-docs/terraform-docs/cmd/pretty"
	"github.com/terraform-docs/terraform-docs/cmd/tfvars"
	"github.com/terraform-docs/terraform-docs/cmd/toml"
	"github.com/terraform-docs/terraform-docs/cmd/version"
	"github.com/terraform-docs/terraform-docs/cmd/xml"
	"github.com/terraform-docs/terraform-docs/cmd/yaml"
	"github.com/terraform-docs/terraform-docs/internal/cli"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	if err := NewCommand().Execute(); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return err
	}
	return nil
}

// NewCommand returns a new cobra.Command for 'root' command
func NewCommand() *cobra.Command {
	config := cli.DefaultConfig()
	cmd := &cobra.Command{
		Args:          cobra.MaximumNArgs(1),
		Use:           "terraform-docs [PATH]",
		Short:         "A utility to generate documentation from Terraform modules in various output formats",
		Long:          "A utility to generate documentation from Terraform modules in various output formats",
		Version:       version.Full(),
		SilenceUsage:  true,
		SilenceErrors: true,
		Annotations:   cli.Annotations("root"),
		PreRunE:       cli.PreRunEFunc(config),
		RunE:          cli.RunEFunc(config),
	}

	// flags
	cmd.PersistentFlags().StringVarP(&config.File, "config", "c", ".terraform-docs.yml", "config file name")

	cmd.PersistentFlags().StringSliceVar(&config.Sections.Show, "show", []string{}, "show section [header, inputs, modules, outputs, providers, requirements, resources]")
	cmd.PersistentFlags().StringSliceVar(&config.Sections.Hide, "hide", []string{}, "hide section [header, inputs, modules, outputs, providers, requirements, resources]")
	cmd.PersistentFlags().BoolVar(&config.Sections.ShowAll, "show-all", true, "show all sections")
	cmd.PersistentFlags().BoolVar(&config.Sections.HideAll, "hide-all", false, "hide all sections (default false)")

	cmd.PersistentFlags().BoolVar(&config.Sort.Enabled, "sort", true, "sort items")
	cmd.PersistentFlags().BoolVar(&config.Sort.By.Required, "sort-by-required", false, "sort items by name and print required ones first (default false)")
	cmd.PersistentFlags().BoolVar(&config.Sort.By.Type, "sort-by-type", false, "sort items by type of them (default false)")

	cmd.PersistentFlags().StringVar(&config.HeaderFrom, "header-from", "main.tf", "relative path of a file to read header from")

	cmd.PersistentFlags().BoolVar(&config.OutputValues.Enabled, "output-values", false, "inject output values into outputs (default false)")
	cmd.PersistentFlags().StringVar(&config.OutputValues.From, "output-values-from", "", "inject output values from file into outputs (default \"\")")

	// deprecation
	cmd.PersistentFlags().BoolVar(&config.Sections.Deprecated.NoHeader, "no-header", false, "do not show module header")
	cmd.PersistentFlags().BoolVar(&config.Sections.Deprecated.NoInputs, "no-inputs", false, "do not show inputs")
	cmd.PersistentFlags().BoolVar(&config.Sections.Deprecated.NoOutputs, "no-outputs", false, "do not show outputs")
	cmd.PersistentFlags().BoolVar(&config.Sections.Deprecated.NoProviders, "no-providers", false, "do not show providers")
	cmd.PersistentFlags().BoolVar(&config.Sections.Deprecated.NoRequirements, "no-requirements", false, "do not show module requirements")
	cmd.PersistentFlags().BoolVar(&config.Sort.Deprecated.NoSort, "no-sort", false, "do no sort items")

	cmd.PersistentFlags().MarkDeprecated("no-header", "use '--hide header' instead")             //nolint:errcheck
	cmd.PersistentFlags().MarkDeprecated("no-inputs", "use '--hide inputs' instead")             //nolint:errcheck
	cmd.PersistentFlags().MarkDeprecated("no-outputs", "use '--hide outputs' instead")           //nolint:errcheck
	cmd.PersistentFlags().MarkDeprecated("no-providers", "use '--hide providers' instead")       //nolint:errcheck
	cmd.PersistentFlags().MarkDeprecated("no-requirements", "use '--hide requirements' instead") //nolint:errcheck
	cmd.PersistentFlags().MarkDeprecated("no-sort", "use '--sort=false' instead")                //nolint:errcheck

	// formatter subcommands
	cmd.AddCommand(asciidoc.NewCommand(config))
	cmd.AddCommand(json.NewCommand(config))
	cmd.AddCommand(markdown.NewCommand(config))
	cmd.AddCommand(pretty.NewCommand(config))
	cmd.AddCommand(tfvars.NewCommand(config))
	cmd.AddCommand(toml.NewCommand(config))
	cmd.AddCommand(xml.NewCommand(config))
	cmd.AddCommand(yaml.NewCommand(config))

	// other subcommands
	cmd.AddCommand(completion.NewCommand())
	cmd.AddCommand(version.NewCommand())

	return cmd
}

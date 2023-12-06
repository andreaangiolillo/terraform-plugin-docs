// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-docs/internal/provider"
)

type migrateCmd struct {
	commonCmd

	flagProviderDir         string
	flagOldWebsiteSourceDir string
	flagNewWebsiteSourceDir string
}

func (cmd *migrateCmd) Synopsis() string {
	return "migrates website files from old directory structure (website/docs/r) to tfplugindocs supported structure (templates/)"
}

func (cmd *migrateCmd) Help() string {
	strBuilder := &strings.Builder{}

	longestName := 0
	longestUsage := 0
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if len(f.Name) > longestName {
			longestName = len(f.Name)
		}
		if len(f.Usage) > longestUsage {
			longestUsage = len(f.Usage)
		}
	})

	strBuilder.WriteString("\nUsage: tfplugindocs migrate [<args>]\n\n")
	cmd.Flags().VisitAll(func(f *flag.Flag) {
		if f.DefValue != "" {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s  (default: %q)\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
				f.DefValue,
			))
		} else {
			strBuilder.WriteString(fmt.Sprintf("    --%s <ARG> %s%s%s\n",
				f.Name,
				strings.Repeat(" ", longestName-len(f.Name)+2),
				f.Usage,
				strings.Repeat(" ", longestUsage-len(f.Usage)+2),
			))
		}
	})
	strBuilder.WriteString("\n")

	return strBuilder.String()
}

func (cmd *migrateCmd) Flags() *flag.FlagSet {
	fs := flag.NewFlagSet("migrate", flag.ExitOnError)

	fs.StringVar(&cmd.flagProviderDir, "provider-dir", "", "relative or absolute path to the root provider code directory when running the command outside the root provider code directory")
	fs.StringVar(&cmd.flagOldWebsiteSourceDir, "old-website-source-dir", "website", "old website directory based on provider-dir")
	fs.StringVar(&cmd.flagNewWebsiteSourceDir, "new-website-source-dir", "templates", "new website templates directory based on provider-dir")
	return fs
}

func (cmd *migrateCmd) Run(args []string) int {
	fs := cmd.Flags()
	err := fs.Parse(args)
	if err != nil {
		cmd.ui.Error(fmt.Sprintf("unable to parse flags: %s", err))
		return 1
	}

	return cmd.run(cmd.runInternal)
}

func (cmd *migrateCmd) runInternal() error {
	err := provider.Migrate(
		cmd.ui,
		cmd.flagProviderDir,
		cmd.flagOldWebsiteSourceDir,
		cmd.flagNewWebsiteSourceDir,
	)
	if err != nil {
		return fmt.Errorf("unable to migrate website: %w", err)
	}

	return nil
}

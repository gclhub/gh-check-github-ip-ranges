package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// For testing purposes
var osExit = os.Exit

func main() {
	cmd := &cobra.Command{
		Use:   "gh-check-github-ip-ranges <ip-address>",
		Short: "Check if an IP address is within GitHub's published IP ranges",
		Long: `Check if a given IP address is within GitHub's published IP ranges.
The ranges are fetched from GitHub's /meta API endpoint. Only IPv4 addresses
are supported at this time.`,
		Args: cobra.ExactArgs(1),
		RunE: runCommand,
		SilenceUsage: true,
		SilenceErrors: true,
	}

	cmd.Flags().BoolP("silent", "s", false, "Silent mode - only use exit codes")

	if err := cmd.Execute(); err != nil {
		if !cmd.Flags().Changed("silent") {
			if err.Error() == "The provided IP address is not a GitHub-owned address" {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			} else {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
		}
		
		// Determine exit code based on error type
		switch {
		case err.Error() == "The provided IP address is not a GitHub-owned address":
			osExit(1)
		default:
			osExit(2)
		}
	}
}

func runCommand(cmd *cobra.Command, args []string) error {
	ipAddress := args[0]
	silent, _ := cmd.Flags().GetBool("silent")

	checker := NewIPChecker()
	result, err := checker.CheckIP(ipAddress)
	if err != nil {
		return err
	}

	if !result.IsGitHubIP {
		return fmt.Errorf("The provided IP address is not a GitHub-owned address")
	}

	if !silent {
		fmt.Printf("IP %s belongs to GitHub's %s range (%s)\n", 
			ipAddress, result.FunctionalArea, result.Range)
	}
	return nil
}
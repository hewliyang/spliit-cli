package cmd

import (
	"github.com/spf13/cobra"
)

var groupID string

var rootCmd = &cobra.Command{
	Use:   "spliit",
	Short: "CLI tool for managing shared expenses using Spliit",
	Long: `Spliit CLI is a command-line interface for interacting with the Spliit API.

Manage shared expenses, track balances, and handle reimbursements for your groups.

Examples:
  spliit --group <id> balances
  spliit -g <id> expenses --limit 10
  spliit -g <id> add-expense "Groceries" "Alice" 5000`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&groupID, "group", "g", "", "group ID (required for most commands)")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func GetGroupID() string {
	return groupID
}

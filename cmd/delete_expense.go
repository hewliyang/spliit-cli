package cmd

import (
	"fmt"

	"github.com/m1a1/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var deleteExpenseCmd = &cobra.Command{
	Use:   "delete-expense <expense-id>",
	Short: "Delete an expense",
	Long: `Delete an expense by its ID.
Use 'spliit expenses --show-ids' to see expense IDs.

Example:
  spliit delete-expense 88Yc54rAzQFx8bWXmxGU9`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetGroupID() == "" {
			return fmt.Errorf("--group flag is required")
		}
		expenseID := args[0]

		client := api.NewClient(GetGroupID())

		if err := client.DeleteExpense(expenseID); err != nil {
			return fmt.Errorf("failed to delete expense: %w", err)
		}

		fmt.Println("✓ Expense deleted successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteExpenseCmd)
}

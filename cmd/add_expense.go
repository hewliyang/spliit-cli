package cmd

import (
	"fmt"
	"strconv"

	"github.com/hewliyang/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var expenseCategory int

var addExpenseCmd = &cobra.Command{
	Use:   "add-expense <title> <payer> <amount>",
	Short: "Add a new expense",
	Long: `Add a new expense to the group. Amount is in cents (5000 = $50.00).
The expense is automatically split equally among all participants.

Categories:
  0 = General (default)
  1 = Entertainment
  2 = Food
  3 = Transport

Examples:
  spliit add-expense "Groceries" "Alice" 5000
  spliit add-expense "Movie tickets" "Bob" 3500 --category 1`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetGroupID() == "" {
			return fmt.Errorf("--group flag is required")
		}
		title := args[0]
		payerName := args[1]
		amount, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid amount: %w", err)
		}

		client := api.NewClient(GetGroupID())

		// Get payer ID
		payerID, err := client.GetParticipantID(payerName)
		if err != nil {
			return fmt.Errorf("failed to find participant: %w", err)
		}
		if payerID == "" {
			return fmt.Errorf("participant %q not found", payerName)
		}

		// Get all participants for equal split
		participants, err := client.GetParticipants()
		if err != nil {
			return fmt.Errorf("failed to get participants: %w", err)
		}

		paidFor := make([]api.PaidForInput, len(participants))
		for i, p := range participants {
			paidFor[i] = api.PaidForInput{
				ParticipantID: p.ID,
				Shares:        1,
			}
		}

		expense, err := client.AddExpense(title, payerID, paidFor, amount, expenseCategory)
		if err != nil {
			return fmt.Errorf("failed to add expense: %w", err)
		}

		fmt.Println("✓ Expense added successfully!")
		fmt.Printf("\nTitle: %s\n", title)
		fmt.Printf("Amount: $%.2f\n", float64(amount)/100)
		fmt.Printf("Paid by: %s\n", payerName)

		if expense.ID != "" {
			fmt.Printf("\nView: https://spliit.app/groups/%s/expenses/%s/edit\n", GetGroupID(), expense.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addExpenseCmd)

	addExpenseCmd.Flags().IntVar(&expenseCategory, "category", 0, "expense category (0=General, 1=Entertainment, 2=Food, 3=Transport)")
}

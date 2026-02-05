package cmd

import (
	"fmt"
	"strconv"

	"github.com/hewliyang/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var expenseCategory int
var isReimbursement bool
var reimbursementTo string

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

Use --reimbursement flag with --to to record a debt settlement.

Examples:
  spliit add-expense "Groceries" "Alice" 5000
  spliit add-expense "Movie tickets" "Bob" 3500 --category 1
  spliit add-expense "Settle up" "Bob" 2500 --reimbursement --to "Alice"`,
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

		payerID, err := client.GetParticipantID(payerName)
		if err != nil {
			return fmt.Errorf("failed to find participant: %w", err)
		}
		if payerID == "" {
			return fmt.Errorf("participant %q not found", payerName)
		}

		var paidFor []api.PaidForInput

		if isReimbursement {
			if reimbursementTo == "" {
				return fmt.Errorf("--to flag is required for reimbursements")
			}
			toID, err := client.GetParticipantID(reimbursementTo)
			if err != nil {
				return fmt.Errorf("failed to find recipient: %w", err)
			}
			if toID == "" {
				return fmt.Errorf("recipient %q not found", reimbursementTo)
			}
			paidFor = []api.PaidForInput{
				{ParticipantID: toID, Shares: 1},
			}
		} else {
			participants, err := client.GetParticipants()
			if err != nil {
				return fmt.Errorf("failed to get participants: %w", err)
			}
			paidFor = make([]api.PaidForInput, len(participants))
			for i, p := range participants {
				paidFor[i] = api.PaidForInput{
					ParticipantID: p.ID,
					Shares:        1,
				}
			}
		}

		expense, err := client.AddExpense(title, payerID, paidFor, amount, expenseCategory, isReimbursement)
		if err != nil {
			return fmt.Errorf("failed to add expense: %w", err)
		}

		if isReimbursement {
			fmt.Println("✓ Reimbursement added successfully!")
			fmt.Printf("\nTitle: %s\n", title)
			fmt.Printf("Amount: $%.2f\n", float64(amount)/100)
			fmt.Printf("From: %s → To: %s\n", payerName, reimbursementTo)
		} else {
			fmt.Println("✓ Expense added successfully!")
			fmt.Printf("\nTitle: %s\n", title)
			fmt.Printf("Amount: $%.2f\n", float64(amount)/100)
			fmt.Printf("Paid by: %s\n", payerName)
		}

		if expense.ID != "" {
			fmt.Printf("\nView: https://spliit.app/groups/%s/expenses/%s/edit\n", GetGroupID(), expense.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addExpenseCmd)

	addExpenseCmd.Flags().IntVar(&expenseCategory, "category", 0, "expense category (0=General, 1=Entertainment, 2=Food, 3=Transport)")
	addExpenseCmd.Flags().BoolVar(&isReimbursement, "reimbursement", false, "mark expense as a reimbursement (settling a debt)")
	addExpenseCmd.Flags().StringVar(&reimbursementTo, "to", "", "recipient of reimbursement (required with --reimbursement)")
}

package cmd

import (
	"fmt"
	"time"

	"github.com/m1a1/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	expenseStartDate string
	expenseEndDate   string
	expenseLimit     int
	expenseOffset    int
	expensePage      int
	showExpenseIDs   bool
)

var expensesCmd = &cobra.Command{
	Use:   "expenses",
	Short: "List expenses",
	Long: `List expenses with optional date filtering.

Examples:
  spliit expenses
  spliit expenses --page 2
  spliit expenses --limit 10 --page 3
  spliit expenses --from 2025-11-01
  spliit expenses --from 2025-11-01 --to 2025-11-30
  spliit expenses --limit 10 --show-ids`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetGroupID() == "" {
			return fmt.Errorf("--group flag is required")
		}
		client := api.NewClient(GetGroupID())

		// Calculate offset from page if provided
		offset := expenseOffset
		if expensePage > 0 {
			offset = (expensePage - 1) * expenseLimit
		}

		expenses, err := client.GetExpenses(expenseLimit, offset)
		if err != nil {
			return fmt.Errorf("failed to get expenses: %w", err)
		}

		// Filter by date if specified
		var filtered []api.Expense
		for _, exp := range expenses {
			include := true

			if expenseStartDate != "" {
				start, err := time.Parse("2006-01-02", expenseStartDate)
				if err != nil {
					return fmt.Errorf("invalid start date format: %w", err)
				}
				if exp.ExpenseDate.Before(start) {
					include = false
				}
			}

			if expenseEndDate != "" {
				end, err := time.Parse("2006-01-02", expenseEndDate)
				if err != nil {
					return fmt.Errorf("invalid end date format: %w", err)
				}
				// Include the end date
				end = end.Add(24 * time.Hour)
				if exp.ExpenseDate.After(end) {
					include = false
				}
			}

			if include {
				filtered = append(filtered, exp)
			}
		}

		if len(filtered) == 0 {
			fmt.Println("No expenses found.")
			return nil
		}

		// Show pagination info
		pageNum := 1
		if expensePage > 0 {
			pageNum = expensePage
		} else if offset > 0 {
			pageNum = (offset / expenseLimit) + 1
		}
		fmt.Printf("Expenses (showing %d, page %d):\n\n", len(filtered), pageNum)

		var total int64
		for _, exp := range filtered {
			date := exp.ExpenseDate.Format("2006-01-02")
			idSuffix := ""
			if showExpenseIDs {
				idSuffix = fmt.Sprintf(" (ID: %s)", exp.ID)
			}

			fmt.Printf("[%s] %s%s\n", date, exp.Title, idSuffix)
			fmt.Printf("  Amount: $%.2f\n", float64(exp.Amount)/100)
			fmt.Printf("  Paid by: %s\n", exp.PaidBy.Name)

			category := "Uncategorized"
			if exp.Category != nil {
				category = exp.Category.Name
			}
			fmt.Printf("  Category: %s\n", category)

			if len(exp.PaidFor) > 0 {
				names := ""
				for i, pf := range exp.PaidFor {
					if i > 0 {
						names += ", "
					}
					names += pf.Participant.Name
				}
				fmt.Printf("  Split between: %s\n", names)
			}
			fmt.Println()

			total += exp.Amount
		}

		fmt.Printf("Total: $%.2f\n", float64(total)/100)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(expensesCmd)

	expensesCmd.Flags().StringVar(&expenseStartDate, "from", "", "filter expenses from this date (YYYY-MM-DD)")
	expensesCmd.Flags().StringVar(&expenseEndDate, "to", "", "filter expenses until this date (YYYY-MM-DD)")
	expensesCmd.Flags().IntVar(&expenseLimit, "limit", 20, "number of expenses per page")
	expensesCmd.Flags().IntVar(&expenseOffset, "offset", 0, "number of expenses to skip (for pagination)")
	expensesCmd.Flags().IntVar(&expensePage, "page", 0, "page number (1-indexed, overrides --offset)")
	expensesCmd.Flags().BoolVar(&showExpenseIDs, "show-ids", false, "show expense IDs (for deletion)")
}

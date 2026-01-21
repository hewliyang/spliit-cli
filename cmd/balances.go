package cmd

import (
	"fmt"

	"github.com/m1a1/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var balancesCmd = &cobra.Command{
	Use:   "balances",
	Short: "Show balances and suggested reimbursements",
	Long:  `Display current balances for all participants and suggested reimbursements to settle debts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetGroupID() == "" {
			return fmt.Errorf("--group flag is required")
		}
		client := api.NewClient(GetGroupID())

		balances, err := client.GetBalances()
		if err != nil {
			return fmt.Errorf("failed to get balances: %w", err)
		}

		participants, err := client.GetParticipants()
		if err != nil {
			return fmt.Errorf("failed to get participants: %w", err)
		}

		// Create ID to name mapping
		idToName := make(map[string]string)
		for _, p := range participants {
			idToName[p.ID] = p.Name
		}

		fmt.Println("Balances:")
		fmt.Println()

		for participantID, balance := range balances.Balances {
			name := idToName[participantID]
			if name == "" {
				name = participantID
			}

			status := "(to receive)"
			if balance.Total < 0 {
				status = "(to pay)"
			}

			fmt.Printf("%s:\n", name)
			fmt.Printf("  Paid: $%.2f\n", float64(balance.Paid)/100)
			fmt.Printf("  Owed: $%.2f\n", float64(balance.PaidFor)/100)
			fmt.Printf("  Balance: $%.2f %s\n", float64(balance.Total)/100, status)
			fmt.Println()
		}

		if len(balances.Reimbursements) > 0 {
			fmt.Println("Suggested Reimbursements:")
			for _, r := range balances.Reimbursements {
				fromName := idToName[r.From]
				toName := idToName[r.To]
				if fromName == "" {
					fromName = r.From
				}
				if toName == "" {
					toName = r.To
				}
				fmt.Printf("  %s → %s: $%.2f\n", fromName, toName, float64(r.Amount)/100)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(balancesCmd)
}

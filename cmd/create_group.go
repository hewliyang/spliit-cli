package cmd

import (
	"fmt"
	"strings"

	"github.com/hewliyang/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	currencySymbol string
	currencyCode   string
)

var createGroupCmd = &cobra.Command{
	Use:   "create-group <name> <participants>",
	Short: "Create a new group",
	Long: `Create a new Spliit group with the given name and participants.
Participants should be a comma-separated list of names.

Examples:
  spliit create-group "Trip to Japan" "Alice,Bob,Charlie"
  spliit create-group "NYC Trip" "Alice,Bob" --currency "$" --currency-code USD`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		participantNames := strings.Split(args[1], ",")

		var cleaned []string
		for _, n := range participantNames {
			n = strings.TrimSpace(n)
			if len(n) >= 2 {
				cleaned = append(cleaned, n)
			}
		}

		if len(cleaned) == 0 {
			return fmt.Errorf("at least one participant with 2+ characters required")
		}

		fmt.Printf("Creating group: %s\n", name)
		fmt.Printf("Participants: %s\n", strings.Join(cleaned, ", "))
		fmt.Printf("Currency: %s (%s)\n\n", currencySymbol, currencyCode)

		client := api.NewClient("")

		newGroupID, err := client.CreateGroup(name, cleaned, currencySymbol, currencyCode)
		if err != nil {
			return fmt.Errorf("failed to create group: %w", err)
		}

		fmt.Println("✓ Group created successfully!")
		fmt.Println()
		fmt.Printf("Group ID: %s\n", newGroupID)
		fmt.Printf("URL: https://spliit.app/groups/%s\n", newGroupID)
		fmt.Println()
		fmt.Println("To use this group, set SPLIIT_GROUP_ID environment variable or use --group flag")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(createGroupCmd)

	createGroupCmd.Flags().StringVar(&currencySymbol, "currency", "$", "currency symbol")
	createGroupCmd.Flags().StringVar(&currencyCode, "currency-code", "SGD", "ISO-4217 currency code")
}

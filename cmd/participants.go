package cmd

import (
	"fmt"

	"github.com/m1a1/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var participantsCmd = &cobra.Command{
	Use:   "participants",
	Short: "List all participants in the group",
	Long:  `Display all participants with their IDs.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetGroupID() == "" {
			return fmt.Errorf("--group flag is required")
		}
		client := api.NewClient(GetGroupID())

		participants, err := client.GetParticipants()
		if err != nil {
			return fmt.Errorf("failed to get participants: %w", err)
		}

		fmt.Println("Participants:")
		for _, p := range participants {
			fmt.Printf("  %s: %s\n", p.Name, p.ID)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(participantsCmd)
}

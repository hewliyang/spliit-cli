package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/hewliyang/spliit-cli/internal/api"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Get group details",
	Long:  `Retrieve and display full group information including name, currency, and participants.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if GetGroupID() == "" {
			return fmt.Errorf("--group flag is required")
		}
		client := api.NewClient(GetGroupID())

		group, err := client.GetGroup()
		if err != nil {
			return fmt.Errorf("failed to get group: %w", err)
		}

		output, err := json.MarshalIndent(group, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)
}

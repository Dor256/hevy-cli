package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"hevy_cli/internal/hevyapi"
	"hevy_cli/internal/middleware"
)

func listWorkouts(hevyClient hevyapi.HevyAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all workouts.",
		Long:  "List all of the finished workouts under the user.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pretty, err := cmd.Flags().GetBool("pretty")
			if err != nil {
				return err
			}
			page, err := cmd.Flags().GetInt("page")
			if err != nil {
				return err
			}
			pageSize, err := cmd.Flags().GetInt("pageSize")
			if err != nil {
				return err
			}

			result, err := hevyClient.ListWorkouts(cmd.Context(), page, pageSize)
			if err != nil {
				if errors.Is(err, middleware.ErrUnauthenticated) {
					return errLogin
				}
				return fmt.Errorf("Error listing workouts %w", err)
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			if pretty {
				enc.SetIndent("", "  ")
			}
			return enc.Encode(result.Workouts)
		},
	}

	addPrettyFlag(cmd)
	addPaginationFlags(cmd)
	return cmd
}

func WorkoutsCmd(hevyClient hevyapi.HevyAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workouts",
		Short: "Manage Hevy workouts.",
		Long:  `Use to perform actions on Hevy workouts.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := hevyClient.ListWorkouts(cmd.Context(), 1, 10)
			if err != nil {
				return fmt.Errorf("Error listing workouts %w", err)
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			return enc.Encode(result.Workouts)
		},
	}
	cmd.AddCommand(listWorkouts(hevyClient))
	return cmd
}

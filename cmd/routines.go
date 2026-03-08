package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/spf13/cobra"
	"hevy_cli/internal/hevyapi"
	"hevy_cli/internal/middleware"
)

func listRoutines(hevyClient hevyapi.HevyAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all workout routines.",
		Long:  "List all of the workout routines under the user.",
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

			result, err := hevyClient.ListRoutines(cmd.Context(), page, pageSize)
			if err != nil {
				if errors.Is(err, middleware.ErrUnauthenticated) {
					return errLogin
				}
				return fmt.Errorf("Error listing routines %w", err)
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			if pretty {
				enc.SetIndent("", "  ")
			}
			return enc.Encode(result.Routines)
		},
	}
	addPrettyFlag(cmd)
	addPaginationFlags(cmd)
	return cmd
}

func getRoutine(hevyClient hevyapi.HevyAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single workout routine.",
		Long:  "Get the details of a single workout routine by id. Provide the id argument to this command.",
		RunE: func(cmd *cobra.Command, args []string) error {
			pretty, err := cmd.Flags().GetBool("pretty")
			if err != nil {
				return err
			}
			id, err := cmd.Flags().GetString("id")
			if err != nil {
				return err
			}

			result, err := hevyClient.GetRoutine(cmd.Context(), id)
			if err != nil {
				if errors.Is(err, middleware.ErrUnauthenticated) {
					return errLogin
				}
				return fmt.Errorf("Error getting routine with id: %s", id)
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			if pretty {
				enc.SetIndent("", "  ")
			}
			return enc.Encode(result.Routine)
		},
	}
	addPrettyFlag(cmd)
	cmd.Flags().String("id", "", "The id of the routine to fetch. Must be in UUID form.")
	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}
	return cmd
}

func updateRoutineLongDescription() string {
	schema, err := jsonschema.For[hevyapi.UpdateRoutineRequest](nil)
	if err != nil {
		panic(err)
	}
	schemaJSON, err := json.MarshalIndent(schema, "  ", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(`Update a workout routine on Hevy.

		The --body flag accepts a JSON string matching this JSON Schema:

		%s`, string(schemaJSON))
}

func updateRoutine(hevyClient hevyapi.HevyAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a routine.",
		Long:  updateRoutineLongDescription(),
		RunE: func(cmd *cobra.Command, args []string) error {
			pretty, err := cmd.Flags().GetBool("pretty")
			if err != nil {
				return err
			}
			id, err := cmd.Flags().GetString("id")
			if err != nil {
				return err
			}
			body, err := cmd.Flags().GetString("body")
			if err != nil {
				return err
			}

			var request hevyapi.UpdateRoutineRequest
			err = json.Unmarshal([]byte(body), &request)
			if err != nil {
				return err
			}
			result, err := hevyClient.UpdateRoutine(cmd.Context(), id, &request)
			if err != nil {
				if errors.Is(err, middleware.ErrUnauthenticated) {
					return errLogin
				}
				return fmt.Errorf("Error listing routines %w", err)
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			if pretty {
				enc.SetIndent("", "  ")
			}
			return enc.Encode(result.Routine)
		},
	}
	addPrettyFlag(cmd)
	cmd.Flags().String("id", "", "The id of the routine to fetch. Must be in UUID form.")
	cmd.Flags().String("body", "", "JSON string representing the routine update payload. See --help for schema.")
	return cmd
}

func RoutinesCmd(hevyClient hevyapi.HevyAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routines",
		Short: "Manage Hevy workout routines.",
		Long:  `Use to perform actions on Hevy workout routines.`,
	}
	cmd.AddCommand(
		listRoutines(hevyClient),
		getRoutine(hevyClient),
		updateRoutine(hevyClient),
	)
	return cmd
}

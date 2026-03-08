package cmd

import "github.com/spf13/cobra"

func addPrettyFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("pretty", "p", false, "Pretty print result JSON with indentation.")
}

func addPaginationFlags(cmd *cobra.Command) {
	cmd.Flags().Int("page", 1, "Page number to fetch for paginated endpoint (default is 1).")
	cmd.Flags().Int("pageSize", 10, "Max number of results per page for paginated endpoint (default is 10).")
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"hevy_cli/internal/assets"
)

func SetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup skill and references.",
		Long:  "Use to set up the SKILL.md and all references for AI agents to use.",
		RunE: func(cmd *cobra.Command, args []string) error {
			subFS, err := fs.Sub(assets.Skill, "skill")
			if err != nil {
				return err
			}
			home, _ := os.UserHomeDir()
			skillPath := fmt.Sprintf("%s/.claude/skills/workout-coach", home)
			err = os.MkdirAll(skillPath, 0755)
			if err != nil {
				return err
			}
			err = os.CopyFS(skillPath, subFS)
			if err != nil {
				return err
			}
			fmt.Printf("Successfuly set up agent skill in %s\n", skillPath)
			return nil
		},
	}

	return cmd
}

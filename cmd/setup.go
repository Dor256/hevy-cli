package cmd

import (
	"embed"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
)

var embeddedSkills embed.FS

func SetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Setup skill and references.",
		Long:  "Use to set up the SKILL.md and all references for AI agents to use.",
		RunE: func(cmd *cobra.Command, args []string) error {
			subFS, err := fs.Sub(embeddedSkills, "skill")
			// entries, err := embeddedSkills.ReadDir("skill")
			if err != nil {
				return err
			}
			home, _ := os.UserHomeDir()
			skillPath := fmt.Sprintf("%s/.claude/skills/workout-coach", home)
			err = os.MkdirAll(skillPath, 0755)
			if err != nil {
				return err
			}
			return os.CopyFS(skillPath, subFS)
		},
	}

	return cmd
}

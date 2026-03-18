package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetup(test *testing.T) {
	test.Run("copies the skill directory to the correct path", func(test *testing.T) {
		tmpHome := test.TempDir()
		test.Setenv("HOME", tmpHome)

		cmd := SetupCmd()
		_, err := executeCommand(cmd)
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}

		skillFile := filepath.Join(tmpHome, ".claude", "skills", "workout-coach", "SKILL.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			test.Fatalf("expected %s to exist", skillFile)
		}
	})

	test.Run("removes the existing skill directory and copies the new one", func(test *testing.T) {
		tmpHome := test.TempDir()
		test.Setenv("HOME", tmpHome)

		skillDir := filepath.Join(tmpHome, ".claude", "skills", "workout-coach")
		err := os.MkdirAll(skillDir, 0755)
		if err != nil {
			test.Fatalf("failed to create skill dir: %v", err)
		}
		staleFile := filepath.Join(skillDir, "stale.txt")
		err = os.WriteFile(staleFile, []byte("old content"), 0644)
		if err != nil {
			test.Fatalf("failed to write stale file: %v", err)
		}

		cmd := SetupCmd()
		_, err = executeCommand(cmd)
		if err != nil {
			test.Fatalf("unexpected error: %v", err)
		}

		if _, err := os.Stat(staleFile); !os.IsNotExist(err) {
			test.Fatalf("expected stale file %s to be removed", staleFile)
		}

		skillFile := filepath.Join(skillDir, "SKILL.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			test.Fatalf("expected %s to exist", skillFile)
		}
	})
}

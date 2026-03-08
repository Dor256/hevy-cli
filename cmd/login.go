package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"os"
	"path/filepath"
)

func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log into Hevy via API key.",
		Long:  "Provide an API key from your Hevy app settings to be used for requests to the Hevy servers.",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, _ := os.UserHomeDir()
			configDir := filepath.Join(home, ".config", "hevy")
			configPath := filepath.Join(configDir, ".hevy.yaml")

			if err := os.MkdirAll(configDir, 0700); err != nil {
				return fmt.Errorf("Error creating config directory %w", err)
			}

			fmt.Print("Enter your Hevy API key: ")
			byteKey, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return fmt.Errorf("Error reading input %w", err)
			}
			apiKey := string(byteKey)
			fmt.Println("\nKey received.")

			viper.Set("api_key", apiKey)
			err = viper.WriteConfigAs(configPath)
			if err != nil {
				err = viper.SafeWriteConfigAs(configPath)
			}
			if err != nil {
				return fmt.Errorf("Failed to save config %w", err)
			}
			os.Chmod(configPath, 0600)
			fmt.Printf("Successfuly logged in!")
			return nil
		},
	}

	return cmd
}

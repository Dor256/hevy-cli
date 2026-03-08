package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"hevy_cli/cmd"
	"hevy_cli/internal/hevyapi"
	"os"
	"os/signal"
)

var rootCmd = &cobra.Command{
	Use:   "hevy",
	Short: "A powerful CLI for managing Hevy.",
	Long: `Hevy is a workout management application.

	This CLI allows reading and modifying workouts and routines from the hevy app for a given user.
	The user will be inferenced by the HEVY_API_KEY.`,
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	hevyClient := hevyapi.NewClient()
	home, _ := os.UserHomeDir()
	viper.AddConfigPath(fmt.Sprintf("%s/.config/hevy", home))
	viper.SetConfigName(".hevy")
	viper.SetConfigType("yaml")
	_ = viper.ReadInConfig()

	rootCmd.SetContext(ctx)
	rootCmd.AddCommand(
		cmd.LoginCmd(),
		cmd.WorkoutsCmd(hevyClient),
		cmd.RoutinesCmd(hevyClient),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

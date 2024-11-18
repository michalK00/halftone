package cmd

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) int {
	_ = godotenv.Load("app.env")

	rootCmd := &cobra.Command{
		Use:   "halftone",
		Short: "Halftone is a service designed for photographers to share their work with their clients and receive feedback",
	}
	rootCmd.AddCommand(APICmd(ctx))
	rootCmd.AddCommand(SchedulerCmd(ctx))

	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

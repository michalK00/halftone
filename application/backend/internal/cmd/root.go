package cmd

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"os"
)

func Execute(ctx context.Context) int {
	_ = godotenv.Load(".env")

	rootCmd := &cobra.Command{
		Use:   "halftone",
		Short: "Halftone is a service designed for photographers to share their work with their clients and receive feedback",
	}
	rootCmd.AddCommand(APICmd(ctx))
	rootCmd.AddCommand(SchedulerCmd(ctx))

	if err := rootCmd.Execute(); err != nil {
		log.Error("command failed ", err)
		log.Error(os.Getenv("MONGODB_URI"))
		return 1
	}
	return 0
}

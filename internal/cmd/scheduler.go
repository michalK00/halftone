package cmd

import (
	"context"
	"fmt"
	"github.com/michalK00/sg-qr/internal/cmdutil"
	"github.com/spf13/cobra"
)

func SchedulerCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scheduler",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the job scheduler",
		RunE: func(cmd *cobra.Command, args []string) error {

			logger := cmdutil.NewLogger("scheduler")
			defer func() { _ = logger.Sync() }()

			db, err := cmdutil.NewMongoClient()
			if err != nil {
				return fmt.Errorf("failed to connect to mongodb: %w", err)
			}
			defer func() { _ = db.Client().Disconnect(context.Background()) }()

			rdb, err := cmdutil.NewRedisClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to connect to redis: %w", err)
			}
			defer func() { _ = rdb.Close() }()

			return nil
		},
	}
	return cmd
}

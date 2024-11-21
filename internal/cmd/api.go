package cmd

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/michalK00/sg-qr/internal/api"
	"github.com/michalK00/sg-qr/internal/cmdutil"
	"github.com/spf13/cobra"
	"os"
)

func APICmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Args:  cobra.ExactArgs(0),
		Short: "Runs the RESTful API",
		RunE: func(cmd *cobra.Command, args []string) error {
			port := "8080"
			if os.Getenv("PORT") != "" {
				port = os.Getenv("PORT")
			}

			//logger := cmdutil.NewLogger("api")
			//defer func() { _ = logger.Sync() }()

			db, err := cmdutil.NewMongoClient()
			if err != nil {
				return fmt.Errorf("could not connect to mongodb: %w", err)
			}
			defer func() { _ = db.Client().Disconnect(context.Background()) }()

			rdb, err := cmdutil.NewRedisClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to connect to redis: %w", err)
			}
			defer func() { _ = rdb.Close() }()

			a := api.NewApi(db, rdb)
			app := a.Server()

			go func() {
				_ = app.Listen("0.0.0.0:" + port)
			}()

			log.Info("started api ", "port ", port)

			<-ctx.Done()

			_ = app.Shutdown()

			return nil
		},
	}
	return cmd
}

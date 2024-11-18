package cmd

import (
	"context"
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

			db, err := cmdutil.NewMongoClient()
			if err != nil {
				return err
			}
			defer db.Client().Disconnect(context.Background())

			a := api.NewApi(db)
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

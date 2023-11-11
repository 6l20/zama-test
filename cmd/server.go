package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/6l20/zama-test/common/log/zap"
	"github.com/6l20/zama-test/infra/api"
	"github.com/6l20/zama-test/server"
	"github.com/6l20/zama-test/server/config"
	"github.com/6l20/zama-test/server/usecases"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server",
	Long: `Will start the http server on the port specified in the config env.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		
		logConfig, err := zap.EnvConfig()
		logger, err := zap.NewLogger(logConfig)

		serverConfig, err := config.EnvServerConfig()
		if err != nil {
			fmt.Println("Error getting server config:", err)
			return
		}

		serverUseCases := usecases.NewServerUseCases(logger, server.NewServer(logger, *serverConfig), nil)

		r := api.Router(*serverUseCases)

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

		srv := &http.Server{
			Addr:    ":" + serverConfig.Port,
			Handler: r,
		}
		go func() {
			logger.Fatal(srv.ListenAndServe().Error())
		}()
		logger.Info("The service is ready to listen and serve.")
	
		killSignal := <-interrupt
		switch killSignal {
		case os.Kill:
			logger.Info("Got SIGKILL...")
		case os.Interrupt:
			logger.Info("Got SIGINT...")
		case syscall.SIGTERM:
			logger.Info("Got SIGTERM...")
		}
	
		logger.Info("The service is shutting down...")
		srv.Shutdown(context.Background())
		logger.Info("Done")

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

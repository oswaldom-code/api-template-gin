package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
	"github.com/oswaldom-code/api-template-gin/pkg/log"
	"github.com/oswaldom-code/api-template-gin/src/adapters/http/rest/infrastructure"
	"github.com/spf13/cobra"
)

func Initialize() {
	config.LoadConfiguration()
	config.SetEnvironment(config.GetEnvironmentConfig().Environment)
}

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run the server",
		Long:  `The web server hosts the API and manages the authentication middleware.`,
		Run: func(cmd *cobra.Command, args []string) {
			StartServer()
		},
	}
}

func StartServer() {
	r := infrastructure.NewServer()
	uri := config.GetServerConfig().AsUri()

	srv := &http.Server{
		Addr:    uri,
		Handler: r,
	}
	
	logConfig := log.LogConfig{
		LogToFile: true,
		FilePath:  "./error.log",
	}

	log.ConfigureLogger(logConfig)

	go func() {
		log.Info("Server running", log.Fields{"uri": uri})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed:", log.Fields{"error": err.Error()})
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", log.Fields{"error": err.Error()})
	}

	log.Info("Server gracefully stopped.")
}

func Execute() {
	Initialize()

	rootCmd := &cobra.Command{
		Use: "api-template", // Replace with your desired command name
	}

	rootCmd.AddCommand(newServeCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Command execution failed:", log.Fields{"error": err.Error()})
		os.Exit(1)
	}
}

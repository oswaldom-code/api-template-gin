package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oswaldom-code/api-template-gin/pkg/config"
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

	go func() {
		log.Printf("Server running at: %s", uri)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped.")
}

func Execute() {
	Initialize()

	rootCmd := &cobra.Command{
		Use: "api-template", // Replace with your desired command name
	}

	rootCmd.AddCommand(newServeCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
		os.Exit(1)
	}
}

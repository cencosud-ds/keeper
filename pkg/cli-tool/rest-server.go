package cli_tool

import (
	"context"
	"github.com/spf13/cobra"
	"keeper/pkg/repository"
	"keeper/pkg/rest-server/handler"
	"keeper/pkg/rest-server/server"
	"keeper/pkg/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// restServer starts the REST server for encrypting
var restServer = &cobra.Command{
	Use:   "server",
	Short: "Starts the REST server",
	Run: func(_ *cobra.Command, _ []string) {
		client, err := createKMSClient()
		if err != nil {
			log.Println(err)
		}

		r := repository.NewRepository(client, encryptionKey)
		s := service.NewService(r, nil)
		h := handler.NewHandler(s)
		srv := server.NewServer(h)

		wait := time.Second * 30
		// Run our server in a goroutine so that it doesn't block.
		go func() {
			err := srv.ListenAndServe()
			if err != nil {
				log.Println(err)
			}
		}()
		log.Println("server started in port 8080")

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		// Block until we receive our signal.
		<-c

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), wait)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			return
		}

		log.Println("shutting down")
		os.Exit(0)
	},
}

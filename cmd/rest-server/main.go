package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	flag "github.com/spf13/pflag"
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
a
func main() {
	awsProfile := flag.StringP("aws-profile", "p", "", "aws profile to use, if not set, uses default")
	encryptionKey := flag.StringP("key", "k", "", "arn of the KMS key to be used for encrypting and decrypting values")
	flag.Parse()

	if *awsProfile == "" {
		err := flag.Set("aws-profile", os.Getenv("AWS_PROFILE"))
		if err != nil {
			log.Printf("error setting aws-profile: %v", err)
		}
	}

	if *encryptionKey == "" {
		err := flag.Set("key", os.Getenv("ENCRYPTION_KEY_ARN"))
		if err != nil {
			log.Printf("error setting encryption key: %v", err)
		}
	}

	log.Printf("Creating KMS client with profile: %v\n", *awsProfile)
	log.Printf("Using KMS key: %v\n", *encryptionKey)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(*awsProfile),
	)
	if err != nil {
		log.Println(err)
	}

	kmsClient := kms.NewFromConfig(cfg)

	r := repository.NewRepository(kmsClient, *encryptionKey)
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
}

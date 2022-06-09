package cli_tool

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/spf13/cobra"
	"keeper/pkg/repository"
	"keeper/pkg/service"
	"log"
	"os"
)

// decrypt represents the decrypt command
var decrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypts terraform.tfvars.encrypted files and generates a plain file",
	Run: func(_ *cobra.Command, _ []string) {
		cfg, err := config.LoadDefaultConfig(
			context.TODO(),
			config.WithSharedConfigProfile(awsProfile),
		)
		if err != nil {
			log.Println(err)
		}
		kmsClient := kms.NewFromConfig(cfg)

		r := repository.NewRepository(kmsClient, encryptionKey)

		s := service.NewService(nil, r)

		log.Println("Reading file: " + encryptedFile)
		fileData, err := os.ReadFile(encryptedFile)
		if err != nil {
			log.Fatalf("error reading file %v: %v", encryptedFile, err)
		}

		log.Println("Decrypting file")
		decryptedDataString, err := s.DecryptTerraform(string(fileData))
		if err != nil {
			log.Fatalf("error decrypting file: %v", err)
		}

		log.Println("Saving into file: " + plainFile)
		file, err := os.Create(plainFile)
		if err != nil {
			log.Fatalf("error creating creating %v: %v", plainFile, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Println(fmt.Errorf("error closing %v: %w", plainFile, err))
			}
		}()

		_, err = file.WriteString(decryptedDataString)
		if err != nil {
			log.Fatalf("error writing to %v: %v", plainFile, err)
		}
	},
}

package cli_tool

import (
	"context"
	"github.com/spf13/cobra"
	"keeper/pkg/repository"
	"keeper/pkg/service"
	"log"
	"os"
)

// decrypt represents the decrypt command
var encrypt = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypts terraform.tfvars input files and generates an encrypted file",
	Run: func(_ *cobra.Command, _ []string) {
		client, err := createKMSClient()
		if err != nil {
			log.Println(err)
		}

		r := repository.NewRepository(client, encryptionKey)
		s := service.NewService(r, nil)

		log.Println("Reading file: " + plainFile)
		fileData, err := os.ReadFile(plainFile)
		if err != nil {
			log.Fatalf("error reading file %v: %v", plainFile, err)
		}

		log.Println("Encrypting into: " + encryptedFile)
		encryptedData, err := s.Encrypt(context.TODO(), string(fileData))
		if err != nil {
			log.Fatalf("error encrypting file: %v", err)
		}

		file, err := os.Create(encryptedFile)
		if err != nil {
			log.Fatalf("error creating %v: %v", encryptedFile, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		_, err = file.WriteString(encryptedData)
		if err != nil {
			log.Fatalf("error writing to %v: %v", encryptedFile, err)
		}
	},
}

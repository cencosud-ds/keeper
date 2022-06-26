package cli_tool

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

// root represents the base command, it is called when there are no subcommands
var root = &cobra.Command{
	Use:  "keeper",
	Long: "Tool to encrypt sensitive data (API keys, credentials, etc) to be used in Terraform repositories and CI/CD pipelines",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the root.
func Execute() {
	cobra.CheckErr(root.Execute())
}

var awsProfile string
var encryptionKey string

func init() {
	// Makes error line shows on "log" lib usage
	log.SetFlags(log.LstdFlags | log.Llongfile)

	root.PersistentFlags().StringVarP(&awsProfile, "aws-profile", "p", "", "aws profile to use, if not set, uses default")
	root.PersistentFlags().StringVarP(&encryptionKey, "encryption-key", "k", "", "arn of the KMS key to be used for encrypting and decrypting values")

	if awsProfile == "" {
		awsProfile = os.Getenv("AWS_PROFILE")
	}

	if encryptionKey == "" {
		encryptionKey = os.Getenv("ENCRYPTION_KEY_ARN")
	}

	log.Printf("Creating KMS client with profile: %v\n", awsProfile)
	log.Printf("Using KMS key: %v\n", encryptionKey)

	root.AddCommand(restServer)
	root.AddCommand(decrypt)
	root.AddCommand(encrypt)
	root.AddCommand(version)
}

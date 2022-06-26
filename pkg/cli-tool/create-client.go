package cli_tool

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func createKMSClient() (*kms.Client, error) {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(awsProfile),
	)
	if err != nil {
		return nil, err
	}

	return kms.NewFromConfig(cfg), nil
}

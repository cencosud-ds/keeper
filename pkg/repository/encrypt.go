package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (r Repository) Encrypt(ctx context.Context, data string) (string, error) {
	input := &kms.EncryptInput{
		KeyId: aws.String(r.encryptionKeyName),
	}

	input.Plaintext = []byte(data)

	output, err := r.client.Encrypt(ctx, input)
	if err != nil {
		return "", fmt.Errorf("error in repository encrypting string: %w", err)
	}

	return r.encoder.EncodeToString(output.CiphertextBlob), err
}

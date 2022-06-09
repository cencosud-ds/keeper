package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func (r Repository) Decrypt(data string) (string, error) {
	input := &kms.DecryptInput{
		KeyId: aws.String(r.encryptionKeyName),
	}

	decodedInput, err := r.encoder.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("error in repository decoding: %w", err)
	}

	input.CiphertextBlob = decodedInput

	output, err := r.client.Decrypt(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("error in repository decrypting: %w", err)
	}

	return string(output.Plaintext), nil
}

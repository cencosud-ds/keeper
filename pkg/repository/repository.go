package repository

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type client interface {
	Encrypt(ctx context.Context, params *kms.EncryptInput, optFns ...func(*kms.Options)) (*kms.EncryptOutput, error)
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

type Repository struct {
	encryptionKeyName string
	client            client
	encoder           *base64.Encoding
}

func NewRepository(c client, encryptionKey string) *Repository {
	r := &Repository{
		encoder:           getBase64Encoder(),
		client:            c,
		encryptionKeyName: encryptionKey,
	}

	return r
}

func getBase64Encoder() *base64.Encoding {
	return base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
}

package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"testing"
)

func (c clientMock) Encrypt(_ context.Context, params *kms.EncryptInput, _ ...func(*kms.Options)) (*kms.EncryptOutput, error) {
	output := &kms.EncryptOutput{
		CiphertextBlob: params.Plaintext,
	}

	return output, nil
}

func TestRepository_Encrypt(t *testing.T) {
	c := clientMock{}
	r := NewRepository(c, "")

	result, err := r.Encrypt(context.TODO(), "test")
	if err != nil {
		t.Fatalf("error encrypting: %v", err)
	}

	expected := "dGVzdA=="

	if result != expected {
		t.Errorf("expected '%v', but got '%v'", expected, result)
	}
}

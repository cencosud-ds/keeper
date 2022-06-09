package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"testing"
)

func (clientMock) Decrypt(_ context.Context, params *kms.DecryptInput, _ ...func(*kms.Options)) (*kms.DecryptOutput, error) {
	output := &kms.DecryptOutput{
		Plaintext: params.CiphertextBlob,
	}

	return output, nil
}

func TestRepository_Decrypt(t *testing.T) {
	c := clientMock{}
	r := NewRepository(c, "")

	result, err := r.Decrypt("dGVzdA==")
	if err != nil {
		t.Fatalf("error decrypting: %v", err)
	}

	expected := "test"

	if result != expected {
		t.Errorf("expected '%v', but got '%v'", expected, result)
	}
}

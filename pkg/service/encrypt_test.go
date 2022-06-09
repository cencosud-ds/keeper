package service

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"
)

type mockedEncryptor struct{}

func (mockedEncryptor) Encrypt(_ context.Context, data string) (string, error) {
	encoder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")

	if data == "corrupted input" {
		return "", errors.New("encryption error")
	}

	return encoder.EncodeToString([]byte(data)), nil
}

func TestService_Encrypt_ValidSingleLineInput(t *testing.T) {
	e := mockedEncryptor{}
	service := NewService(e, nil)

	result, err := service.Encrypt(context.TODO(), `db_password = "fake_db"`)
	if err != nil {
		t.Errorf("error encrypting input: %v", err)
	}

	expected := `db_password = ImZha2VfZGIi`
	if result != expected {
		t.Errorf("expected '%v', but got '%v'", expected, result)
	}
}

func TestService_Encrypt_ValidMultiLineInput(t *testing.T) {
	e := mockedEncryptor{}
	service := NewService(e, nil)

	input := `db_password = "fake_db"
db_user = "fake_user"`

	result, err := service.Encrypt(context.TODO(), input)
	if err != nil {
		t.Errorf("error encrypting input: %v", err)
		return
	}

	expected := `db_password = ImZha2VfZGIi
db_user = ImZha2VfdXNlciI=`

	if result != expected {
		t.Errorf("expected '%v', but got '%v'", expected, result)
	}
}

func TestService_Encrypt_EncryptionError(t *testing.T) {

	e := mockedEncryptor{}
	service := NewService(e, nil)

	_, err := service.Encrypt(context.TODO(), "db_pass = corrupted input")

	if err == nil {
		t.Fatalf("expected an error")
	}

	if err.Error() != "error in encrypt service encrypting terraform string: encryption error" {
		t.Errorf("expected encryption error but got: %v", err)
	}
}

func TestService_Encrypt_InvalidFormat(t *testing.T) {

	e := mockedEncryptor{}
	service := NewService(e, nil)

	_, err := service.Encrypt(context.TODO(), "fake_db")
	if err == nil {
		t.Fatalf("expected an error")
	}

	if err.Error() != "error in encrypt service transforming terraform input: invalid format" {
		t.Errorf("expected invalid format error but got: %v", err)
	}

	_, err = service.Encrypt(context.TODO(), "db_name: fake_db")
	if err == nil {
		t.Fatalf("expected an error")
	}

	if err.Error() != "error in encrypt service transforming terraform input: invalid format" {
		t.Errorf("expected invalid format error but got: %v", err)
	}
}

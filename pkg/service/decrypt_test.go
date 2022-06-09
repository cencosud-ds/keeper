package service

import (
	"encoding/base64"
	"errors"
	"testing"
)

type mockedDecryptor struct{}

func (mockedDecryptor) Decrypt(encodedData string) (string, error) {
	encoder := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	data, err := encoder.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func TestService_Decrypt_ValidInput(t *testing.T) {
	d := mockedDecryptor{}
	service := NewService(nil, d)

	singleLineInput := "db_password = ImZha2VfZGIi"

	result, err := service.DecryptTerraform(singleLineInput)
	if err != nil {
		t.Fatalf("error decrypting %v: %v", singleLineInput, err)
	}

	expected := `db_password = "fake_db"`

	if result != expected {
		t.Fatalf("expected '%v', but got '%v'", expected, result)
	}

	multiLineInput := `db_password = ImZha2VfZGIi
db_user = ImZha2VfdXNlciI=`

	result, err = service.DecryptTerraform(multiLineInput)
	if err != nil {
		t.Fatalf("error decrypting %v: %v", multiLineInput, err)
	}

	expected = `db_password = "fake_db"
db_user = "fake_user"`

	if result != expected {
		t.Errorf("expected '%v', but got '%v'", expected, result)
	}
}

func TestService_Decrypt_InvalidInput(t *testing.T) {
	d := mockedDecryptor{}
	service := NewService(nil, d)

	input := "db_password: ImZha2VfZGIi"

	_, err := service.DecryptTerraform(input)
	if err == nil {
		t.Fatalf("expected an error")
	}

	expected := `error in service transforming terraform map: invalid format`
	if err.Error() != expected {
		t.Fatalf("expected '%v', but got '%v'", expected, err.Error())
	}
}

func TestService_Decrypt_CommentsAndLineBreaks(t *testing.T) {
	d := mockedDecryptor{}
	service := NewService(nil, d)

	input := `db_password = ImZha2VfZGIi
# This is a comment

db_user = ImZha2VfdXNlciI= # this is an inline comment
db_user_2 = ImZha2VfdXNlciI= # this is another inline comment

`

	result, err := service.DecryptTerraform(input)
	if err != nil {
		t.Fatalf("error decrypting: %v", err)
	}

	expected := `db_password = "fake_db"
db_user = "fake_user"
db_user_2 = "fake_user"`

	if result != expected {
		t.Errorf("expected '%v', but got '%v'", expected, result)
	}
}

func TestService_Decrypt_InvalidValues(t *testing.T) {
	d := mockedDecryptor{}
	service := NewService(nil, d)

	_, err := service.DecryptTerraform("db_password = whoops, not encrypted!")

	if err == nil {
		t.Fatalf("expected an error")
	}

	var expectedError base64.CorruptInputError
	isCorruptInputError := errors.As(err, &expectedError)

	if !isCorruptInputError {
		t.Errorf("error was not a CorruptInputError, got: %T", errors.Unwrap(err))
	}
}

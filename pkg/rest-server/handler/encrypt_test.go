package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func defaultRequest(body string) *http.Request {
	req := httptest.NewRequest("POST", "/encrypt", strings.NewReader(body))
	req.Header.Set("Content-Type", "text/plain")

	return req
}

type mockedEncryptor struct{}

func (mockedEncryptor) Encrypt(_ context.Context, data string) (string, error) {
	if data == "variable: wrong body format" {
		return "", errors.New("cannot process data")
	}

	if data == `db_name = "fake_db"` {
		return `db_name = ImZha2VfZGIi`, nil
	}

	return "", errors.New("unexpected finish for encrypt function")
}

func TestHandler_Encrypt_CorrectBody(t *testing.T) {
	handler := NewHandler(mockedEncryptor{})

	body := `db_name = "fake_db"`

	res := httptest.NewRecorder()
	req := defaultRequest(body)

	handler.Encrypt(res, req)

	result := res.Body.String()
	expected := `db_name = ImZha2VfZGIi`

	if res.Code != http.StatusOK {
		t.Errorf("Expected HTTP response 200 but got %d", res.Code)
	}

	if result != expected {
		t.Errorf("Expected body to contain value %q but got %q", expected, result)
	}
}

func TestHandler_Encrypt_WrongBodyFormat(t *testing.T) {
	handler := NewHandler(mockedEncryptor{})

	body := `variable: wrong body format`

	res := httptest.NewRecorder()
	req := defaultRequest(body)

	handler.Encrypt(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP response 400 but got %d", res.Code)
	}

	if res.Body.String() != encryptionErrorMessage {
		t.Errorf("Expected encryption error but got %v", res.Body.String())
	}
}

func TestHandler_Encrypt_WrongHeader(t *testing.T) {
	handler := NewHandler(mockedEncryptor{})

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/encrypt", nil)
	req.Header.Set("Content-Type", "application/json")

	handler.Encrypt(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP response 400 but got %d", res.Code)
	}

	if res.Body.String() != wrongContentTypeMessage {
		t.Errorf("Expected wrong content type error but got %v", res.Body.String())
	}
}

type errReader int

func (errReader) Read([]byte) (int, error) {
	return 0, errors.New("returning error on purpose")
}

func TestHandler_Encrypt_UnreadableBody(t *testing.T) {
	handler := NewHandler(mockedEncryptor{})

	var reader errReader

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/encrypt", reader)
	req.Header.Set("Content-Type", "text/plain")

	handler.Encrypt(res, req)

	if res.Code != http.StatusBadRequest {
		t.Errorf("Expected HTTP response 400 but got %d", res.Code)
	}

	if res.Body.String() != parsingErrorMessage {
		t.Errorf("Expected parsing error but got %v", res.Body.String())
	}
}

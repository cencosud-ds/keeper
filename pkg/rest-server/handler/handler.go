package handler

import "context"

type encryptor interface {
	Encrypt(ctx context.Context, data string) (string, error)
}

type Handler struct {
	encryptor encryptor
}

func NewHandler(e encryptor) *Handler {
	return &Handler{encryptor: e}
}

package service

import (
	"context"
)

type encryptor interface {
	Encrypt(ctx context.Context, data string) (string, error)
}

type decryptor interface {
	Decrypt(data string) (string, error)
}

type Service struct {
	encryptor encryptor
	decryptor decryptor
}

func NewService(e encryptor, d decryptor) *Service {
	s := &Service{
		encryptor: e,
		decryptor: d,
	}

	return s
}

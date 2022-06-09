package service

import (
	"context"
	"fmt"
	"strings"
)

func (s Service) Encrypt(ctx context.Context, data string) (string, error) {

	splitData := strings.Split(data, "\n")
	dataMap, orderedKeys, err := stringToTerraformMap(splitData)
	if err != nil {
		return "", fmt.Errorf("error in encrypt service transforming terraform input: %w", err)
	}

	encryptedDataMap := make(map[string]string)
	for _, key := range orderedKeys {
		encryptedDataMap[key], err = s.encryptor.Encrypt(ctx, dataMap[key])
		if err != nil {
			return "", fmt.Errorf("error in encrypt service encrypting terraform string: %w", err)
		}
	}

	return terraformMapToString(encryptedDataMap, orderedKeys), nil

}

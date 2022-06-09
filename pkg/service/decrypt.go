package service

import (
	"fmt"
	"strings"
)

func (s Service) DecryptTerraform(data string) (string, error) {

	// remove comments and empty lines and then splits
	sanitizedData := sanitizeData(data)

	splitData := strings.Split(sanitizedData, "\n")

	dataMap, orderedKeys, err := stringToTerraformMap(splitData)
	if err != nil {
		return "", fmt.Errorf("error in service transforming terraform map: %w", err)
	}

	decryptedDataMap := make(map[string]string)
	for _, k := range orderedKeys {
		decryptedDataMap[k], err = s.decryptor.Decrypt(dataMap[k])
		if err != nil {
			return "", fmt.Errorf("error in service decrypting: %w", err)
		}
	}

	return terraformMapToString(decryptedDataMap, orderedKeys), nil
}

// sanitizeData removes empty lines and comments with hash within a string
func sanitizeData(data string) string {
	splitData := strings.Split(data, "\n")
	sanitizedLines := make([]string, 0)

	for _, line := range splitData {

		// Remove comments with hash
		line = strings.Split(line, "#")[0]
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		sanitizedLines = append(sanitizedLines, line)
	}

	return strings.Join(sanitizedLines, "\n")
}

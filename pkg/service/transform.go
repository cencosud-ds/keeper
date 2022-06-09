package service

import (
	"errors"
	"fmt"
	"strings"
)

// stringToMapTerraform parses a string slice and returns a map
// and a string array for navigating the map in an ordered fashion (based on the insertion sequence)
// the function expects a multi-line string with the key value pairs separated by an equals sign, example:
// `key1 = value1
//	key2 = value2`
func stringToTerraformMap(splitData []string) (map[string]string, []string, error) {
	m := make(map[string]string)
	insertionOrder := make([]string, 0)

	count := 1
	for _, line := range splitData {

		// If we are at the last value of the array, and we have an empty line we just exit the loop
		if len(splitData) == count && line == "" {
			break
		}

		// SplitN is used instead of Split to avoid splitting values that end with an equals sign ('ImZha2VfdXNlciI=' for example)
		// With SplitN we can make sure we only split the first equals sign in the string
		separatedLine := strings.SplitN(line, "=", 2)

		// If the separated line has 1 value it means it doesn't have the correct format because nothing got separated by an equals sign
		if len(separatedLine) == 1 {
			return nil, nil, errors.New("invalid format")
		}

		insertionOrder = append(insertionOrder, separatedLine[0])
		m[separatedLine[0]] = strings.TrimSpace(separatedLine[1])
		count++
	}

	return m, insertionOrder, nil
}

// terraformMapToString parses a map and builds a string with all the key value pairs and returns it.
func terraformMapToString(data map[string]string, insertionOrder []string) string {
	var s string
	for _, k := range insertionOrder {
		v := data[k]
		s += fmt.Sprintf("%v= %v\n", k, v)
	}

	// Removes last new line before returning the string
	return strings.TrimSuffix(s, "\n")
}

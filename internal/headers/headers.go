package headers

import (
	"errors"
	"fmt"
	"http_protocol/internal/utils"
	"strings"
)


type Headers map[string]string

func NewHeaders() Headers {
    return make(Headers)
}

// parser should only check states and decide where to redirect data to
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	parts := strings.Split(string(data), "\r\n")
	l := len(parts)

	fmt.Printf("parts: %v", parts)

	if l == 0 {
		return
	}

	for _, part := range parts[:l] {
		fieldName, fieldValue, parseErr := parseFieldLine(part)

		if parseErr != nil {
			err = parseErr
			return
		}

		if fieldName != "" && fieldValue != "" {
			h[fieldName] = fieldValue 
		}
	}

	// keep track of number of bytes parsed
	if len(parts) >= 2 {
		done = parts[l-1] == "" && parts[l-2] == "" // means line ended with \r\n\r\n
	}
	
	if !done && len(parts[l-1]) > 0{
		n = l - len(parts[l-1])
	}

	return
}

// Returns Key, Value pair
func parseFieldLine(data string) (string, string , error){
	parts := strings.SplitN(data, ":", 2)
	
	if len(parts) < 2 {
		return "", "", nil
	}

	fieldName, fieldValue := parts[0], parts[1]
	
	if !isValidFieldName(fieldName) {
		return "", "", errors.New("Invalid field name given")
	}

	fieldValue = strings.TrimSpace(fieldValue)

	return fieldName, fieldValue, nil
}

func isValidFieldName(header string) bool {
	return utils.IsAlphabetic(header)
}



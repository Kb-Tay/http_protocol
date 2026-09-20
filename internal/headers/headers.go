package headers

import (
	"errors"
	"http_protocol/internal/utils"
	"strings"
)


type Headers map[string]string

func NewHeaders() Headers {
    return make(Headers)
}

// parser should only check states and decide where to redirect data to
func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	s := string(data)
	
	if s == "\r\n" {
		done = true
		return
	}
	
	parts := strings.SplitN(s, "\r\n", 2)

	if len(parts) < 2 {
		return	
	}

	fieldLine := parts[0]
	fieldName, fieldValue, parseErr := parseFieldLine(fieldLine)

	if parseErr != nil {
		err = parseErr
		return
	}

	if fieldName != "" && fieldValue != "" {
		h[fieldName] = fieldValue
		n += len(fieldLine) + 2 // account for \r\n
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



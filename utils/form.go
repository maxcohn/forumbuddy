package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

//TODO: document
//TODO: unit test all of these
func FormValueIntGtZero(urlValues url.Values, param string) (int, error) {
	// Attempt to parse there string. If it fails to parse or isn't greater than zero, return an error
	// This covers the case where the param doesn't exist, because it would return an empty string and fail to parse
	intVal, err := strconv.Atoi(strings.TrimSpace(urlValues.Get(param)))
	if err != nil || intVal <= 0 {
		return 0, errors.New(fmt.Sprintf("Parameter '%s' wasn't provided or wasn't an integer greater than zero", param))
	}

	return intVal, nil
}

func FormValueStringNonEmpty(urlValues url.Values, param string) (string, error) {
	if !urlValues.Has("text") || strings.TrimSpace(urlValues.Get("text")) != "" {
		return "", errors.New(fmt.Sprintf("Missing parameter or empty parameter '%s'", param))
	}

	return urlValues.Get("text"), nil
}

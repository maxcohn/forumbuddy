package utils

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// Checks the given url.Values for the given param and validates that it is an
// integer greater than zero. Otherwise, returns an error.
func FormValueIntGtZero(urlValues url.Values, param string) (int, error) {
	// Attempt to parse there string. If it fails to parse or isn't greater than zero, return an error
	// This covers the case where the param doesn't exist, because it would return an empty string and fail to parse
	intVal, err := strconv.Atoi(strings.TrimSpace(urlValues.Get(param)))
	if err != nil || intVal <= 0 {
		return 0, errors.New(fmt.Sprintf("Parameter '%s' wasn't provided or wasn't an integer greater than zero", param))
	}

	return intVal, nil
}

// Checks the given ur.Values for the given param and validates that it is a non-empty
// string. Otherwise, returns an error.
func FormValueStringNonEmpty(urlValues url.Values, param string) (string, error) {
	if !urlValues.Has(param) || strings.TrimSpace(urlValues.Get(param)) == "" {
		return "", errors.New(fmt.Sprintf("Missing parameter or empty parameter '%s'", param))
	}

	return urlValues.Get(param), nil
}

var validate = validator.New()
var formDecoder = form.NewDecoder()

//TODO: test
func DecodeAndValidateForm(value any, formData url.Values) error {
	err := formDecoder.Decode(value, formData)
	if err != nil {
		return err
	}

	log.Printf("Value: %v\n", value)

	err = validate.Struct(value)
	if err != nil {
		return err
	}

	return nil
}

package utils_test

import (
	"forumbuddy/utils"
	"net/url"
	"testing"
)

// Test FormValueIntGtZero with valid inputs
func TestFormValueIntGtZeroValid(t *testing.T) {
	form := []struct {
		param    string
		expected int
		input    url.Values
	}{
		{
			"one",
			1,
			url.Values{"one": []string{"1"}},
		},
		{
			"two",
			2,
			url.Values{"two": []string{"2"}},
		},
		{
			"eleven",
			11,
			url.Values{"eleven": []string{"11"}},
		},
		{
			"onemillion",
			1000000,
			url.Values{"onemillion": []string{"1000000"}},
		},
		{
			"twowithspaces",
			2,
			url.Values{"twowithspaces": []string{"   2   "}},
		},
	}

	for _, formVal := range form {
		output, err := utils.FormValueIntGtZero(formVal.input, formVal.param)
		if err != nil {
			t.Errorf("Error in validating form value: %s", err.Error())
		}

		if output != formVal.expected {
			t.Errorf("Output value doesn't match expected: Expected: %d, received: %d.", formVal.expected, output)
		}
	}
}

// Test FormValueIntGtZero with invalid inputs
func TestFormValueIntGtZeroInvalid(t *testing.T) {
	form := []struct {
		param string
		input url.Values
	}{
		{
			"emptystring",
			url.Values{"emptystring": []string{""}},
		},
		{
			"negative",
			url.Values{"negative": []string{"-2"}},
		},
		{
			"zero",
			url.Values{"zero": []string{"0"}},
		},
		{
			"random",
			url.Values{"random": []string{"asdnjkasdnjk"}},
		},
	}

	for _, formVal := range form {
		output, err := utils.FormValueIntGtZero(formVal.input, formVal.param)
		if err == nil {
			t.Errorf("Invalid values passed test. Input: %s. Output: %d ", formVal.input.Get(formVal.param), output)
		}
	}
}

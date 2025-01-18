package chat

import (
	"reflect"
	"testing"
)

func TestParseResults(t *testing.T) {
	// Test case 1: Sample input with multiple tags
	input := "Specific Tags: [Rental Law, Eviction Litigation, Fair Housing Policy, Minority Rights Advocacy, Non-Profit Sector]"

	expected := []string{
		"Rental Law",
		"Eviction Litigation",
		"Fair Housing Policy",
		"Minority Rights Advocacy",
		"Non-Profit Sector",
	}

	result := parseResults(input)

	// Compare the results
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("parseResults(%s) = %v; want %v", input, result, expected)
	}

	// Test case 2: Empty input
	emptyInput := ""
	emptyResult := parseResults(emptyInput)
	if len(emptyResult) != 0 {
		t.Errorf("parseResults empty input = %v; want empty slice", emptyResult)
	}

	// Test case 3: Input without specific tags prefix
	irrelevantInput := "Some other text without specific tags"
	irrelevantResult := parseResults(irrelevantInput)
	if len(irrelevantResult) != 0 {
		t.Errorf("parseResults irrelevant input = %v; want empty slice", irrelevantResult)
	}
}

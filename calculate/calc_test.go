package calc

import (
	"strconv"
	"testing"
)

func TestCalc(t *testing.T) {
	testCases := []struct {
		expression     string
		expectedResult string
		expectError    bool
	}{
		{
			expression:     "43 + (10 - 25) * (2 - 8)",
			expectedResult: "133",
			expectError:    false,
		},
		{
			expression:     "100 / 4",
			expectedResult: "25",
			expectError:    false,
		},
		{
			expression:     "((2+2)",
			expectedResult: "",
			expectError:    true,
		},
		{
			expression:     "5 + g",
			expectedResult: "",
			expectError:    true,
		},
		{
			expression:     "-",
			expectedResult: "",
			expectError:    true,
		},
		{
			expression:     "5 * )",
			expectedResult: "",
			expectError:    true,
		},
		{
			expression:     "100 / 0",
			expectedResult: "",
			expectError:    true,
		},
		{
			expression:     "2**2",
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		result, err := Calc(tc.expression)
		if (err != nil) != tc.expectError {
			t.Fatalf("For expression %s, expected error: %v, got: %v", tc.expression, tc.expectError, err)
		}
		cur, _ := strconv.ParseFloat(tc.expectedResult, 64)
		if !tc.expectError && result != cur {
			t.Fatalf("For expression %s, expected result: %s, got: %f", tc.expression, tc.expectedResult, result)
		}
	}
}

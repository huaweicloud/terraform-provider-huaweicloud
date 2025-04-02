package utils_test

import (
	"reflect"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestMathFunc_Round(t *testing.T) {
	var floatVal = 3.1415926535

	testOutput := utils.Round(floatVal, 2)
	// Test whether the mantissa of positive floating-point numbers less than 5 can be discarded.
	if !reflect.DeepEqual(testOutput, float64(3.14)) {
		t.Fatalf("[Check 1] The processing result of the Round method is not as expected, want %s, but got %s",
			utils.Green(float64(3.14)), utils.Yellow(testOutput))
	}

	testOutput = utils.Round(floatVal, 3)
	// Test whether the mantissa of a positive floating point number greater than 5 can be carried.
	if !reflect.DeepEqual(testOutput, float64(3.142)) {
		t.Fatalf("[Check 2] The processing result of the Round method is not as expected, want %s, but got %s",
			utils.Green(float64(3.142)), utils.Yellow(testOutput))
	}

	// Test whether the mantissa of negative floating-point numbers less than 5 can be discarded.
	testOutput = utils.Round(-floatVal, 2)
	if !reflect.DeepEqual(testOutput, float64(-3.14)) {
		t.Fatalf("[Check 3] The processing result of the Round method is not as expected, want %s, but got %s",
			utils.Green(float64(-3.14)), utils.Yellow(testOutput))
	}

	// Test whether the mantissa of a negative floating point number greater than 5 can be carried.
	testOutput = utils.Round(-floatVal, 3)
	if !reflect.DeepEqual(testOutput, float64(-3.142)) {
		t.Fatalf("[Check 4] The processing result of the Round method is not as expected, want %s, but got %s",
			utils.Green(float64(-3.142)), utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the Round method meets expectation")
}

func TestMathFunc_EqualFloat(t *testing.T) {
	var floatVal = 3.1415

	// Tests whether two positive floating-point numbers to two decimal places are equal.
	testOutput := utils.EqualFloat(utils.Round(floatVal, 2), utils.Round(floatVal, 2))
	if !reflect.DeepEqual(testOutput, true) {
		t.Fatalf("[Check 1] The processing result of the EqualFloat method is not as expected, want %s, but got %s",
			utils.Green(true), utils.Yellow(testOutput))
	}

	// Tests whether two negative floating-point numbers are equal to two decimal places.
	testOutput = utils.EqualFloat(-utils.Round(floatVal, 2), -utils.Round(floatVal, 2))
	if !reflect.DeepEqual(testOutput, true) {
		t.Fatalf("[Check 2] The processing result of the EqualFloat method is not as expected, want %s, but got %s",
			utils.Green(true), utils.Yellow(testOutput))
	}

	// Tests whether two floating-point numbers (one is positive number and anthor is negative number) with the same
	// absolute value to two decimal places are not equal.
	testOutput = utils.EqualFloat(utils.Round(floatVal, 2), -utils.Round(floatVal, 2))
	if !reflect.DeepEqual(testOutput, false) {
		t.Fatalf("[Check 3] The processing result of the EqualFloat method is not as expected, want %s, but got %s",
			utils.Green(false), utils.Yellow(testOutput))
	}

	// Tests whether two positive floating-point numbers with different reserved digits are unequal.
	testOutput = utils.EqualFloat(utils.Round(floatVal, 3), utils.Round(floatVal, 2))
	if !reflect.DeepEqual(testOutput, false) {
		t.Fatalf("[Check 4] The processing result of the EqualFloat method is not as expected, want %s, but got %s",
			utils.Green(false), utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the EqualFloat method meets expectation")
}

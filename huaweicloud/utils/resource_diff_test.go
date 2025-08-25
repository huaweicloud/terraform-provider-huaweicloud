package utils_test

import (
	"reflect"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestResourceDiffunc_ContainsAllKeyValues(t *testing.T) {
	var (
		compareObject = map[string]interface{}{
			"A": map[string]interface{}{
				"Aa": "aa_aa",
				"Ab": "aa_bb",
			},
			"B": map[string]interface{}{
				"Ba": true,
				"Bb": 123,
			},
			"C": "cc",
			"D": 456,
		}
		objects = []map[string]interface{}{
			{
				"A": map[string]interface{}{
					"Aa": "aa_aa",
				},
				"B": map[string]interface{}{
					"Bb": 123,
				},
				"C": "cc",
			},
			{
				"A": map[string]interface{}{
					"Aa": "aa_aaa", // Value changed.
				},
			},
			{
				"B": map[string]interface{}{
					"Ba": false,
				},
			},
			{
				"B": map[string]interface{}{
					"Ba": nil,
				},
			},
			{
				"B": map[string]interface{}{
					"Bb": true, // Change the value type from integer to boolean.
				},
			},
		}
	)

	testOutput := utils.ContainsAllKeyValues(compareObject, objects[0])
	if !reflect.DeepEqual(testOutput, true) {
		t.Fatalf("[Check 1] The processing result of the ContainsAllKeyValues method is not as expected, want %s, but got %s",
			utils.Green(true), utils.Yellow(testOutput))
	}

	testOutput = utils.ContainsAllKeyValues(compareObject, objects[1])
	if !reflect.DeepEqual(testOutput, false) {
		t.Fatalf("[Check 2] The processing result of the ContainsAllKeyValues method is not as expected, want %s, but got %s",
			utils.Green(false), utils.Yellow(testOutput))
	}

	testOutput = utils.ContainsAllKeyValues(compareObject, objects[2])
	if !reflect.DeepEqual(testOutput, false) {
		t.Fatalf("[Check 3] The processing result of the ContainsAllKeyValues method is not as expected, want %s, but got %s",
			utils.Green(false), utils.Yellow(testOutput))
	}

	testOutput = utils.ContainsAllKeyValues(compareObject, objects[3])
	if !reflect.DeepEqual(testOutput, false) {
		t.Fatalf("[Check 4] The processing result of the ContainsAllKeyValues method is not as expected, want %s, but got %s",
			utils.Green(false), utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the ContainsAllKeyValues method meets expectation")
}

func TestResourceDiffunc_FindDecreaseKeys(t *testing.T) {
	var (
		decreaseCacls = []map[string]interface{}{
			{
				"A": map[string]interface{}{
					"Aa": "aa_aa",
					"Ab": "aa_bb",
				},
				"B": map[string]interface{}{
					"Ba": true,
					"Bb": 123,
				},
				"C": "cc",
			},
			{
				"A": map[string]interface{}{
					"Aa": "aa_aaa", // The key is exist, but the value changed.
					"Ac": "aa_cc",
				},
				"B": map[string]interface{}{
					"Ba": 123, // The key is exist, but the value changed.
					"Bb": nil, // The key is exist, but the value changed.
				},
				// The key 'C' is removed.
			},
			{
				"A": map[string]interface{}{
					"Ab": "aa_bb",
				},
				"C": "cc",
			},
		}
	)

	testOutput := utils.FindDecreaseKeys(decreaseCacls[0], decreaseCacls[1])
	if !reflect.DeepEqual(testOutput, decreaseCacls[2]) {
		t.Fatalf("The processing result of the FindDecreaseKeys method is not as expected, want %s, but got %s",
			utils.Green(decreaseCacls[2]), utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the FindDecreaseKeys method meets expectation")
}

func TestResourceDiffunc_TakeObjectsDifferent(t *testing.T) {
	var (
		diffCacls = []map[string]interface{}{
			{
				"A": map[string]interface{}{
					"Aa": "aa_aa",
					"Ab": "aa_bb",
				},
				"B": map[string]interface{}{
					"Ba": true,
					"Bb": 123,
				},
				"C": "cc",
			},
			{
				"A": map[string]interface{}{
					"Ab": "aa_bb",
				},
				"B": map[string]interface{}{
					"Ba": true,
					"Bb": 123,
					"Bc": "bb_cc",
				},
				"D": "dd",
			},
			{
				"A": map[string]interface{}{
					"Aa": "aa_aa",
				},
				"C": "cc",
			},
		}
	)

	testOutput := utils.TakeObjectsDifferent(diffCacls[0], diffCacls[1])
	if !reflect.DeepEqual(testOutput, diffCacls[2]) {
		t.Fatalf("The processing result of the TakeObjectsDifferent method is not as expected, want %s, but got %s",
			utils.Green(diffCacls[2]), utils.Yellow(testOutput))
	}

	t.Logf("All processing results of the TakeObjectsDifferent method meets expectation")
}

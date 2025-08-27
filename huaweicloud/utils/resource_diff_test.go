package utils_test

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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

func TestResourceDiffunc_SuppressStrSliceDiffs(t *testing.T) {
	// Test 1: Test basic functionality with TypeSet length field
	t.Run("TypeSet length field suppression", func(t *testing.T) {
		// Create a mock ResourceData
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2", "tag3"},
			"tags_origin": []interface{}{"tag1", "tag2"},
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: remote addition (should not suppress diff because new count > old count)
		result := suppressFunc("tags.#", "2", "3", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for count increase, got %v", result)
		}

		// Test case: local removal (should suppress diff if no local changes)
		result = suppressFunc("tags.#", "3", "2", resourceData)
		if !result {
			t.Errorf("Expected diff suppression for count decrease, got %v", result)
		}
	})

	// Test 2: Test TypeSet element field suppression
	t.Run("TypeSet element field suppression", func(t *testing.T) {
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2", "tag3"},
			"tags_origin": []interface{}{"tag1", "tag2"},
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: element addition (should not suppress if element not in origin)
		result := suppressFunc("tags.1234567890", "", "tag3", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for new element not in origin, got %v", result)
		}

		// Test case: element removal (should not suppress if element was in origin)
		result = suppressFunc("tags.0987654321", "tag1", "", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for removing element from origin, got %v", result)
		}
	})

	// Test 3: Test main field suppression
	t.Run("main field suppression", func(t *testing.T) {
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2", "tag3"},
			"tags_origin": []interface{}{"tag1", "tag2"},
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: new script contains elements not in console (should not suppress diff)
		result := suppressFunc("tags", "tag1,tag2", "tag1,tag2,tag3", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for new elements, got %v", result)
		}

		// Test case: local additions (should not suppress diff)
		result = suppressFunc("tags", "tag1,tag2", "tag1,tag2,tag4", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for local additions, got %v", result)
		}

		// Test case: local removals (should suppress diff if no local additions)
		result = suppressFunc("tags", "tag1,tag2,tag3", "tag1,tag2", resourceData)
		if !result {
			t.Errorf("Expected diff suppression for removals with no additions, got %v", result)
		}
	})

	// Test 4: Test with empty origin values
	t.Run("empty origin values", func(t *testing.T) {
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2"},
			"tags_origin": []interface{}{},
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: first time setting value (should not suppress diff)
		result := suppressFunc("tags", "", "tag1,tag2", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for first time setting, got %v", result)
		}
	})

	// Test 5: Test with nil origin values
	t.Run("nil origin values", func(t *testing.T) {
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2"},
			"tags_origin": nil,
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: first time setting value (should not suppress diff)
		result := suppressFunc("tags", "", "tag1,tag2", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for first time setting, got %v", result)
		}
	})

	// Test 6: Test with complex scenarios
	t.Run("complex scenarios", func(t *testing.T) {
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2", "tag3", "tag4"},
			"tags_origin": []interface{}{"tag1", "tag2", "tag3"},
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: mixed changes (should not suppress if any local changes)
		result := suppressFunc("tags", "tag1,tag2,tag3", "tag1,tag2,tag4", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for mixed changes, got %v", result)
		}

		// Test case: only remote additions (should not suppress if new elements not in console)
		result = suppressFunc("tags", "tag1,tag2,tag3", "tag1,tag2,tag3,tag4", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for new elements not in console, got %v", result)
		}
	})

	// Test 7: Test edge cases
	t.Run("edge cases", func(t *testing.T) {
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2"},
			"tags_origin": []interface{}{"tag1", "tag2"},
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: empty strings (should not suppress if new elements added)
		result := suppressFunc("tags", "tag1,", "tag1,tag2", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for new elements, got %v", result)
		}

		// Test case: whitespace handling (should not suppress if new elements added)
		result = suppressFunc("tags", " tag1 , tag2 ", "tag1,tag2,tag3", resourceData)
		if result {
			t.Errorf("Expected no diff suppression for new elements, got %v", result)
		}
	})

	// Test 8: Test successful suppression scenarios
	t.Run("successful suppression scenarios", func(t *testing.T) {
		resourceData := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2"},
			"tags_origin": []interface{}{"tag1", "tag2"},
		})

		suppressFunc := utils.SuppressStrSliceDiffs()

		// Test case: same elements, no changes (should suppress diff)
		result := suppressFunc("tags", "tag1,tag2", "tag1,tag2", resourceData)
		if !result {
			t.Errorf("Expected diff suppression for no changes, got %v", result)
		}

		// Test case: only remote additions (should suppress diff if no local changes)
		// Create a scenario where console and new script are the same, but origin is different
		resourceData2 := createMockResourceDataForDiffTest(t, map[string]interface{}{
			"tags":        []interface{}{"tag1", "tag2", "tag3"},
			"tags_origin": []interface{}{"tag1", "tag2"},
		})

		// When console and new script are the same, but origin is different
		// This should suppress diff because there are no local changes
		result = suppressFunc("tags", "tag1,tag2,tag3", "tag1,tag2,tag3", resourceData2)
		if !result {
			t.Errorf("Expected diff suppression for no local changes, got %v", result)
		}
	})

	t.Logf("All processing results of the SuppressStrSliceDiffs method meets expectation")
}

// createMockResourceDataForDiffTest creates a mock ResourceData for diff testing
func createMockResourceDataForDiffTest(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	// Create a simple schema for testing
	schemaMap := make(map[string]*schema.Schema)
	for key, value := range data {
		switch value.(type) {
		case []interface{}:
			schemaMap[key] = &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
			}
		case nil:
			schemaMap[key] = &schema.Schema{
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
			}
		}
	}

	// Create ResourceData with the schema
	resource := &schema.Resource{Schema: schemaMap}
	resourceData := resource.Data(nil)

	// Set the initial values
	for key, value := range data {
		// lintignore:R001
		if err := resourceData.Set(key, value); err != nil {
			t.Fatalf("Failed to set %s: %v", key, err)
		}
	}

	return resourceData
}

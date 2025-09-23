package utils

import (
	"reflect"
	"testing"
)

func TestAccFunction_RemoveNil(t *testing.T) {
	var (
		testInput = map[string]interface{}{
			"level_one_index_zero": map[string]interface{}{},
			"level_one_index_one": []map[string]interface{}{
				{
					"level_two_index_zero": nil,
				},
				{
					"level_two_index_one": "192.168.0.1",
					"level_two_index_two": nil,
				},
			},
			"level_one_index_two": []map[string]interface{}{
				{
					"level_two_index_zero": nil,
				},
			},
			"level_one_index_three": "172.16.0.237",
		}

		expected = map[string]interface{}{
			"level_one_index_one": []map[string]interface{}{
				{
					"level_two_index_one": "192.168.0.1",
				},
			},
			"level_one_index_three": "172.16.0.237",
		}
	)

	testOutput := RemoveNil(testInput)
	if !reflect.DeepEqual(testOutput, expected) {
		t.Fatalf("The processing result of RemoveNil method is not as expected, want %s, but %s", Green(expected), Yellow(testOutput))
	}
	t.Logf("The processing result of RemoveNil method meets expectation: %s", Green(expected))
}

func TestAccFunction_reverse(t *testing.T) {
	var (
		testInput = "abcdefg"
		expected  = "gfedcba"
	)

	if !reflect.DeepEqual(Reverse(testInput), expected) {
		t.Fatalf("The processing result of the function 'Reverse' is not as expected, want '%s', but got '%s'",
			Green(expected), Yellow(testInput))
	}
	t.Logf("The processing result of function 'Reverse' meets expectation: %s", Green(expected))
}

func TestAccFunction_jsonStringsEqual(t *testing.T) {
	var (
		jsonStr1 = "{\n\"key1\":\"value1\",\n\"key2\":\"value2\"\n}"
		jsonStr2 = "{\"key2\":\"value2\",\"key1\":\"value1\"}"
	)

	if !JSONStringsEqual(jsonStr1, jsonStr2) {
		t.Fatalf("The processing result of the function 'JSONStringsEqual' is not as expected, want '%v', "+
			"but got '%v'", Green(true), Yellow(false))
	}
	t.Logf("The processing result of function 'JSONStringsEqual' meets expectation: %s", Green(true))
}

func TestAccFunction_PasswordEncrypt(t *testing.T) {
	var password = "Test@123!"
	encrypted, err := PasswordEncrypt(password)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("The encrypted string of %s is %s", password, Green(encrypted))

	newEncrypted, err := TryPasswordEncrypt(encrypted)
	if err != nil {
		t.Fatal(err)
	}
	if newEncrypted != encrypted {
		t.Fatalf("The encrypted string is not as expected, want %s, but %s", Green(encrypted), Yellow(newEncrypted))
	}
}

func TestAccFunction_ConvertMemoryUnit(t *testing.T) {
	var (
		memories   = []interface{}{2097152, 2048, 2, "2097152", "2048", "2", 2.048, "2.048", 0}
		diffLevels = []int{2, 0, -2, 2, 0, -2, -2, -2, 1}
		expected   = []int{2, 2048, 2097152, 2, 2048, 2097152, -1, -1, 0}
	)

	for i, memory := range memories {
		testOutput := ConvertMemoryUnit(memory, diffLevels[i])
		if !reflect.DeepEqual(testOutput, expected[i]) {
			t.Fatalf("The processing result of ConvertMemoryUnit method is not as expected, want %s, but %s", Green(expected[i]), Yellow(testOutput))
		}
		t.Logf("The processing result of ConvertMemoryUnit method meets expectation: %s", Green(expected[i]))
	}
}

func TestAccFunction_IsUUID(t *testing.T) {
	var (
		ids      = []string{"550e8400-e29b-41d4-a716-446655440000", "550e8400e29b41d4a716446655440000", "abc123", ""}
		expected = []bool{true, true, false, false}
	)

	for i, idInput := range ids {
		if isValid := IsUUID(idInput); isValid != expected[i] {
			t.Fatalf("The processing result of IsUUID method is not as expected, want %s, but %s, the ID is %s",
				Green(expected[i]), Yellow(isValid), idInput)
		}
		t.Logf("The processing result of IsUUID method meets expectation: %s", Green(expected[i]))
	}
}
func TestAccFunction_IsUUIDWithHyphens(t *testing.T) {
	var (
		ids      = []string{"550e8400-e29b-41d4-a716-446655440000", "550e8400e29b41d4a716446655440000", "abc123", ""}
		expected = []bool{true, false, false, false}
	)

	for i, idInput := range ids {
		if isValid := IsUUIDWithHyphens(idInput); isValid != expected[i] {
			t.Fatalf("The processing result of IsUUIDWithHyphens method is not as expected, want %s, but %s, the ID is %s",
				Green(expected[i]), Yellow(isValid), idInput)
		}
		t.Logf("The processing result of IsUUIDWithHyphens method meets expectation: %s", Green(expected[i]))
	}
}

func TestAccFunction_FilterMapWithSameKey(t *testing.T) {
	var (
		rawArray = []map[string]interface{}{
			{"a": "b"},
			{"a": "b"},
			{"a": "b", "c": "d"},
			{"a": "b", "c": "d"},
			{"a": "b"},
			{},
			{},
		}

		filterArray = []map[string]interface{}{
			{"a": "d"},
			{"a": "d", "m": "n"},
			{"a": "d", "m": "n"},
			{"a": "d", "c": "a", "m": "n"},
			{},
			{"m": "n"},
			{},
		}

		expectedArray = []map[string]interface{}{
			{"a": "d"},
			{"a": "d"},
			{"a": "d"},
			{"a": "d", "c": "a"},
			{},
			{},
			{},
		}
	)

	for i := 0; i < 7; i++ {
		rawMap := rawArray[i]
		filterMap := filterArray[i]
		expectedMap := expectedArray[i]
		result := FilterMapWithSameKey(rawMap, filterMap)

		if !reflect.DeepEqual(result, expectedMap) {
			t.Fatalf("The processing result of the function 'FilterMapWithSameKey' is not as expected, want '%v', "+
				"but got '%v'", Green(expectedMap), Yellow(result))
		}
		t.Logf("The processing result of `FilterMapWithSameKey` method meets expectation: %s", Green(expectedMap))
	}
}

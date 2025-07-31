package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awspolicy "github.com/jen20/awspolicyequivalence"
)

// ComposeAnySchemaDiffSuppressFunc allows parameters to determine multiple Diff methods.
// When any method (SchemaDiffSuppressFunc) returns true, this compose function will return true.
// It will only return false when all methods (SchemaDiffSuppressFunc) return false.
func ComposeAnySchemaDiffSuppressFunc(fs ...schema.SchemaDiffSuppressFunc) schema.SchemaDiffSuppressFunc {
	return func(k, o, n string, d *schema.ResourceData) bool {
		for _, f := range fs {
			if f(k, o, n, d) {
				return true
			}
		}
		return false
	}
}

func SuppressEquivalentAwsPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := awspolicy.PoliciesAreEquivalent(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

// Suppress all changes
func SuppressDiffAll(k, old, new string, d *schema.ResourceData) bool {
	return true
}

// The SuppressCaseDiffs method compares the old and new values ​​of the current parameter to determine if their
// contents are consistent while ignoring the case format.
func SuppressCaseDiffs() schema.SchemaDiffSuppressFunc {
	return func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
		return strings.EqualFold(oldVal, newVal)
	}
}

// Suppress changes if we get a computed min_disk_gb if value is unspecified (default 0)
func SuppressMinDisk(k, old, new string, d *schema.ResourceData) bool {
	return new == "0" || old == new
}

// Suppress changes if we get a base64 format or plaint text user_data
func SuppressUserData(k, old, new string, d *schema.ResourceData) bool {
	// user_data is in base64 format
	if HashAndHexEncode(old) == new {
		return true
	}

	// user_data is plaint text
	if plaint, err := base64.StdEncoding.DecodeString(old); err == nil {
		if HashAndHexEncode(string(plaint)) == new {
			return true
		}
	}

	return false
}

func SuppressTrimSpace(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func SuppressLBWhitelistDiffs(k, old, new string, d *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	old_array := strings.Split(old, ",")
	new_array := strings.Split(new, ",")
	sort.Strings(old_array)
	sort.Strings(new_array)

	return reflect.DeepEqual(old_array, new_array)
}

func SuppressSnatFiplistDiffs(k, old, new string, d *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	old_array := strings.Split(old, ",")
	new_array := strings.Split(new, ",")
	sort.Strings(old_array)
	sort.Strings(new_array)

	return reflect.DeepEqual(old_array, new_array)
}

// Suppress changes if we get a string with or without new line
func SuppressNewLineDiffs(k, old, new string, d *schema.ResourceData) bool {
	return strings.Trim(old, "\n") == strings.Trim(new, "\n")
}

func SuppressEquivilentTimeDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldTime, err := time.Parse(time.RFC3339, old)
	if err != nil {
		return false
	}

	newTime, err := time.Parse(time.RFC3339, new)
	if err != nil {
		return false
	}

	return oldTime.Equal(newTime)
}

func SuppressVersionDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldArray := regexp.MustCompile(`[\.\-]+`).Split(old, -1)
	newArray := regexp.MustCompile(`[\.\-]+`).Split(new, -1)
	if len(newArray) > len(oldArray) {
		return false
	}
	for i, v := range newArray {
		if v != oldArray[i] {
			return false
		}
	}
	return true
}

func CompareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	var obj2 interface{}
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Equal(canonicalJson1, canonicalJson2)
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}

func SuppressStringSepratedByCommaDiffs(_, old, new string, _ *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	oldArray := strings.Split(old, ",")
	newArray := strings.Split(new, ",")
	sort.Strings(oldArray)
	sort.Strings(newArray)

	return reflect.DeepEqual(oldArray, newArray)
}

// ContainsAllKeyValues ​​checks whether object A (type map[string]interface{}) recursively contains all the keys and
// corresponding values ​​of object B (type map[string]interface{}).
// If the key-value pair in object B exists in object A and the values ​​are equal (recursively processing nested maps),
// return true; otherwise return false.
func ContainsAllKeyValues(objA, objB map[string]interface{}) bool {
	for key, bVal := range objB {
		aVal, exists := objA[key]
		if !exists {
			return false // A is missing the key of B.
		}

		// Check if the values ​​are both nested maps, if so, recursively compare.
		aMap, aIsMap := aVal.(map[string]interface{})
		bMap, bIsMap := bVal.(map[string]interface{})
		if aIsMap && bIsMap {
			if !ContainsAllKeyValues(aMap, bMap) {
				return false
			}
		} else {
			// Non-map types are compared directly via DeepEqual().
			if !reflect.DeepEqual(bVal, aVal) {
				return false
			}
		}
	}
	return true
}

// FindDecreaseKeys is a method that used to find out the key that objB is missing compared to objA.
// Will ignore the increase parts.
func FindDecreaseKeys(objA, objB map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, valA := range objA {
		if valB, exists := objB[key]; !exists {
			// If the key does not exist in objB, it's considered as a decrease key and is added directly to the result.
			result[key] = valA
		} else {
			// Check if the current values (valA and valB) are both type map for recursive processing.
			mapA, okA := valA.(map[string]interface{})
			mapB, okB := valB.(map[string]interface{})
			// If either valA or valB is not of type map, the subsequent recursive comparison is performed.
			if okA && okB {
				subResult := FindDecreaseKeys(mapA, mapB)
				if len(subResult) > 0 {
					result[key] = subResult
				}
			}
		}
	}
	return result
}

// SuppressObjectDiffs is a method that make the JSON string type parameter ignore the changes made on the console and
// only allow the local script to take effect.
func SuppressObjectDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		if strings.HasSuffix(paramKey, ".%") || strings.HasSuffix(paramKey, ".#") {
			log.Printf("[DEBUG] The current change object is not of type object.")
			return false
		}
		return diffObjectParam(paramKey, o, n, d)
	}
}

// diffObjectParam is used to check whether the parameters of the current object or JSON object type have been modified
// other than those changed in the console.
// The following three scenarios will determine whether the parameter has changed (method return false):
//  1. The new value of the script adds some keys compared to the server return value (which must include keys that do
//     not exist in the value returned by the server).
//  2. The new value of the script modifies some (or all) key/value compared to the server return value.
//  3. The new value of the script removes some (or all) key/value compared to the old value of the script (the key can
//     be a nested structure).
//
// The following are examples of related scenarios:
//
// Service result:
//
//	{
//		"A": {
//			"Aa": "aa_aa",
//			"Ab": "aa_bb"
//		},
//		"B": "bb",
//		"C": "cc",
//		"D": "dd"
//	}
//
// Example 1 (Key 'D' add but the value is the same as the service result, so return true):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc"	->		"C": "cc",
//	}						"D": "dd"
//						}
//
// Example 2 (New key 'D' addreturn false):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc",	->		"C": "cc",
//	}						"E": "ee"
//						}
//
// Example 3 (The value of key 'C' changed, return false):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc",	->		"C": "ccc"
//	}					}
//
// Example 4 (The value of key 'A.Aa' changed, return false):
//
//	{							{
//		"A": {						"A": {
//			"Aa": "aa_aa"				"Aa": "aa_aaa"
//		},					->		},
//		"B": "bb"					"B": "bb"
//	}							}
//
// Example 5 (Key 'D' removed, even it is exist in the service result, return false):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc",	->		"C": "cc"
//		"D": "dd"		}
//	}
func diffObjectParam(paramKey, _, _ string, d *schema.ResourceData) bool {
	var (
		consoleVal, newScriptVal, originVal map[string]interface{}

		originParamKey           = fmt.Sprintf("%s_origin", paramKey)
		oldParamVal, newParamVal = d.GetChange(paramKey)
	)

	// After refresh phase, the value from the service side will be stored in the tfstate, and as old value in the
	// next d.GetChange() method returns.
	consoleVal = TryMapValueAnalysis(oldParamVal)
	newScriptVal = TryMapValueAnalysis(newParamVal)
	// The script value of the last update (if it has) is used as a reference for the historical value of this
	// change.
	originVal = TryMapValueAnalysis(d.Get(originParamKey))

	return ContainsAllKeyValues(consoleVal, newScriptVal) && len(FindDecreaseKeys(originVal, newScriptVal)) < 1
}

func SuppressMapDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		log.Printf("[DEBUG][SuppressMapDiffs] Called with paramKey='%s', old='%s', new='%s'", paramKey, o, n)

		// Ignore length change judgment, because this method will judge each changed key one by one
		if strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressMapDiffs] Ignoring length change for '%s'", paramKey)
			return true
		}

		// Handle the case where the entire map is being changed
		if !strings.Contains(paramKey, ".") {
			return suppressEntireMapDiff(paramKey, d)
		}

		// Handle single key changes
		return suppressSingleKeyDiff(paramKey, d)
	}
}

// suppressEntireMapDiff handles changes to the entire map
func suppressEntireMapDiff(paramKey string, d *schema.ResourceData) bool {
	originMapCategory := fmt.Sprintf("%s_origin", paramKey)
	log.Printf("[DEBUG][EntireMapDiff] Handling entire map change for '%s', origin map category: '%s'",
		paramKey, originMapCategory)

	originMapVal := d.Get(originMapCategory)
	if originMapVal == nil {
		log.Printf("[DEBUG][EntireMapDiff] Origin map '%s' is nil, suppressing diff for entire map '%s'",
			originMapCategory, paramKey)
		return true
	}

	originMap := TryMapValueAnalysis(originMapVal)
	if len(originMap) == 0 {
		log.Printf("[DEBUG][EntireMapDiff] Origin map '%s' is empty, suppressing diff for entire map '%s'",
			originMapCategory, paramKey)
		return true
	}

	// For entire map changes, we need to check if all keys in the new value exist in origin
	// This is a simplified approach - in practice, you might want more sophisticated comparison
	log.Printf("[DEBUG][EntireMapDiff] Entire map '%s' change detected, checking against origin", paramKey)
	return false // For now, report the change and let individual key suppression handle it
}

// suppressSingleKeyDiff handles changes to a single key
func suppressSingleKeyDiff(paramKey string, d *schema.ResourceData) bool {
	categories := strings.Split(paramKey, ".")
	mapCategory := strings.Join(categories[:len(categories)-1], ".")
	originMapCategory := fmt.Sprintf("%s_origin", mapCategory)
	keyName := categories[len(categories)-1]

	log.Printf("[DEBUG][SingleKeyDiff] Processing key '%s', mapCategory='%s', originMapCategory='%s', keyName='%s'",
		paramKey, mapCategory, originMapCategory, keyName)

	// Get origin map (last local configuration)
	originMapVal := d.Get(originMapCategory)
	originMap := TryMapValueAnalysis(originMapVal)
	log.Printf("[DEBUG][SingleKeyDiff] Origin map '%s' content: %+v", originMapCategory, originMap)

	// Get current configuration map
	currentMapVal := GetNestedObjectFromRawConfig(d.GetRawConfig(), mapCategory)
	if currentMapVal == nil {
		log.Printf("[DEBUG][SingleKeyDiff] Current map '%s' is nil, suppressing diff for key '%s'",
			mapCategory, keyName)
		return true
	}

	currentMap := TryMapValueAnalysis(currentMapVal)
	log.Printf("[DEBUG][SingleKeyDiff] Current map '%s' content: %+v", mapCategory, currentMap)

	// Check if the key exists in current configuration
	existsInCurrent := keyExists(currentMap, keyName)
	existsInOrigin := keyExists(originMap, keyName)
	isOriginEmpty := originMapVal == nil || len(originMap) == 0

	// Determine suppression based on key existence
	return determineSuppression(existsInCurrent, existsInOrigin, isOriginEmpty, keyName)
}

// keyExists checks if a key exists in the map
func keyExists(m map[string]interface{}, key string) bool {
	_, exists := m[key]
	return exists
}

// determineSuppression determines whether to suppress the diff based on key existence
func determineSuppression(existsInCurrent, existsInOrigin, isOriginEmpty bool, keyName string) bool {
	if existsInCurrent {
		// The key exists in current configuration
		if isOriginEmpty {
			// Origin is empty or nil, report the change (locally added)
			log.Printf("[DEBUG] The key '%s' found in current config but origin is empty", keyName)
			return false
		}

		if existsInOrigin {
			// The key exists in both current config and origin, report this change
			log.Printf("[DEBUG] The key '%s' found in both current config and origin", keyName)
			return false
		}

		// The key exists in current config but not in origin (locally added)
		log.Printf("[DEBUG] The key '%s' found in current config but not in origin", keyName)
		return false
	}

	// The key does not exist in current configuration
	if isOriginEmpty {
		// Origin is empty or nil, suppress the diff (remotely added)
		log.Printf("[DEBUG] The key '%s' not found in current config and origin is empty, suppressing diff",
			keyName)
		return true
	}

	if existsInOrigin {
		// The key exists in origin but not in current config (locally removed)
		log.Printf("[DEBUG] The key '%s' found in origin but not in current config (locally removed)",
			keyName)
		return false
	}

	// The key does not exist in either current config or origin (remotely added)
	log.Printf("[DEBUG] The key '%s' not found in either current config or origin (remotely added), suppressing diff",
		keyName)
	return true
}

// TakeObjectsDifferent is a method that used to recursively get the complement of object A (objA) compared to
// object B (objB) (including nested differences).
// In Math, it means A-B, also A\B or {x | x∈A，x∉B}.
func TakeObjectsDifferent(objA, objB map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, valueA := range objA {
		valueB, exists := objB[key]

		// The key in objA does not exist in objB
		if !exists {
			result[key] = valueA
			continue
		}

		// Try recursively processing nested map.
		if subMapA, okA := valueA.(map[string]interface{}); okA {
			if subMapB, okB := valueB.(map[string]interface{}); okB {
				// Both sides are maps, recursive comparison.
				subDiff := TakeObjectsDifferent(subMapA, subMapB)
				if len(subDiff) > 0 {
					result[key] = subDiff
				}
			} else {
				// The value of objA is a map but the value of objB is not (type inconsistency).
				result[key] = valueA
			}
			continue
		}

		// Handling non-map types or inconsistent types.
		if !reflect.DeepEqual(valueA, valueB) {
			result[key] = valueA
		}
	}

	return result
}

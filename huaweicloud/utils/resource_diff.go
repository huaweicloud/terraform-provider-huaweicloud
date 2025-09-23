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
	"strconv"
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

// SuppressStrSliceDiffs is a method that makes the string slice type parameter ignore the changes made on the console and
// only allow the local script to take effect. It identifies elements that are decreased compared to origin and
// elements that are newly added remotely.
func SuppressStrSliceDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		log.Printf("[DEBUG][SuppressStrSliceDiffs] Called with paramKey='%s', oldVal='%s', newVal='%s'", paramKey, o, n)

		// Handle TypeSet length field
		if strings.HasSuffix(paramKey, ".#") {
			log.Printf("[DEBUG][SuppressStrSliceDiffs] Processing TypeSet length field: %s", paramKey)
			return diffStrSliceLength(paramKey, o, n, d)
		}

		// Handle TypeSet element fields (e.g., {set_param_key}.1234567890)
		if strings.Contains(paramKey, ".") && !strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressStrSliceDiffs] Processing TypeSet element field: %s", paramKey)
			return diffStrSliceElement(paramKey, o, n, d)
		}

		if strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressStrSliceDiffs] The current change object is not of type slice.")
			return false
		}

		log.Printf("[DEBUG][SuppressStrSliceDiffs] Processing main field: %s", paramKey)
		result := diffStrSliceParam(paramKey, o, n, d)
		log.Printf("[DEBUG][SuppressStrSliceDiffs] Final result: %v", result)
		return result
	}
}

// diffStrSliceLength handles the length field of TypeList or TypeSet
func diffStrSliceLength(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	baseField := strings.TrimSuffix(paramKey, ".#")

	// Get the origin value
	originParamKey := fmt.Sprintf("%s_origin", baseField)
	originVal := d.Get(originParamKey)

	// Get current values
	oldCount, _ := strconv.Atoi(oldVal)
	newCount, _ := strconv.Atoi(newVal)

	// If origin is empty or nil, this is the first time setting the value
	if originVal == nil {
		return false
	}

	// Check if origin is effectively empty
	var originCount int
	var isEmpty bool
	switch v := originVal.(type) {
	case []interface{}:
		originCount = len(v)
		isEmpty = len(v) == 0
	case *schema.Set:
		originCount = v.Len()
		isEmpty = v.Len() == 0
	default:
		originCount = 0
		isEmpty = true
	}

	// If origin is empty, check if this is a remote-only change that should be suppressed
	if isEmpty {
		// Get current remote state to check if this is a remote-only change
		currentVal := d.Get(baseField)
		if currentVal != nil {
			var currentCount int
			switch v := currentVal.(type) {
			case []interface{}:
				currentCount = len(v)
			case *schema.Set:
				currentCount = v.Len()
			default:
				currentCount = 0
			}

			// If new count is less than current count, this might be a remote removal
			// that should be suppressed (unless it's a local removal)
			if newCount < currentCount {
				return true
			}
		}

		return false
	}

	// Check if there are actual changes that require length difference to be shown
	hasLocalAdditions := newCount > oldCount
	hasLocalRemovals := newCount < originCount

	// If there are actual local changes, don't suppress length difference
	if hasLocalAdditions || hasLocalRemovals {
		log.Printf("[DEBUG][diffStrSliceLength] Is local additions happened? %v", hasLocalAdditions)
		log.Printf("[DEBUG][diffStrSliceLength] Is local removals happened? %v", hasLocalRemovals)
		return false
	}

	// If no actual changes, suppress length difference (e.g., remote-only additions/removals)
	return true
}

// diffStrSliceElement handles individual slice elements. And for the TypeSet, there are indexes of each element,
// so we need to handle them separately.
func diffStrSliceElement(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	parts := strings.Split(paramKey, ".")
	if len(parts) != 2 {
		log.Printf("[DEBUG][diffStrSliceElement] Invalid paramKey format: %s", paramKey)
		return false
	}
	baseField := parts[0]

	originParamKey := fmt.Sprintf("%s_origin", baseField)
	originVal := d.Get(originParamKey)

	log.Printf("[DEBUG][diffStrSliceElement] baseField='%s', oldVal='%s', newVal='%s', originVal=%v",
		baseField, oldVal, newVal, originVal)

	// Handle element removal case
	if newVal == "" {
		return handleElementRemoval(oldVal, originVal, baseField, d)
	}

	// Handle element addition/modification case
	return handleElementAddition(oldVal, newVal, originVal, baseField, d)
}

// handleElementRemoval handles the case when an element is being removed
func handleElementRemoval(oldVal string, originVal interface{}, baseField string, d *schema.ResourceData) bool {
	log.Printf("[DEBUG][handleElementRemoval] Element '%s' is being removed, checking if should suppress diff", oldVal)

	// Check if this element was in origin
	if isElementInOrigin(oldVal, originVal) {
		log.Printf("[DEBUG][handleElementRemoval] Element '%s' was in origin, NOT suppressing diff (allow removal)", oldVal)
		return false // NOT suppressing - allow removal of origin elements
	}

	// If origin is empty or nil, check if this element exists in remote state
	if checkElementInRemoteState(baseField, oldVal, d) {
		log.Printf("[DEBUG][handleElementRemoval] Element '%s' exists in remote state but not in origin, suppressing diff (ignore remote removal)",
			oldVal)
		return true // Suppress diff - ignore removal of remote-only elements
	}

	// If element was not in origin and not in remote state, suppress the diff
	log.Printf("[DEBUG][handleElementRemoval] Element '%s' was not in origin or remote state, suppressing diff",
		oldVal)
	return true
}

// handleElementAddition handles the case when an element is being added or modified
func handleElementAddition(oldVal, newVal string, originVal interface{}, baseField string, d *schema.ResourceData) bool {
	// If origin is nil or empty, this is the first time setting the value
	if isOriginEmpty(originVal) {
		return handleFirstTimeSetting(oldVal, newVal, baseField, d)
	}

	// Check if this element is in origin
	if isElementInOrigin(newVal, originVal) {
		// If the element is unchanged (oldVal == newVal), don't suppress diff
		// This ensures Terraform knows the config value still exists
		if oldVal == newVal {
			log.Printf("[DEBUG][handleElementAddition] Element '%s' unchanged and in origin, not suppressing diff to preserve config value", newVal)
			return false
		}
		log.Printf("[DEBUG][handleElementAddition] Element '%s' was in origin, suppressing diff", newVal)
		return true
	}

	// If element was not in origin, don't suppress (this is a local addition)
	log.Printf("[DEBUG][handleElementAddition] Element '%s' was not in origin, not suppressing diff (local addition)", newVal)
	return false
}

// handleFirstTimeSetting handles the case when origin is empty or nil
func handleFirstTimeSetting(oldVal, newVal, baseField string, d *schema.ResourceData) bool {
	// If oldVal is empty, this is a CREATE scenario - use main logic
	if oldVal == "" {
		log.Printf("[DEBUG][handleFirstTimeSetting] This is a CREATE scenario (oldVal=''), using main diffStrSliceParam logic")
		return false // Let the main logic handle this
	}
	// If oldVal is not empty, this is an UPDATE scenario - check remote state
	log.Printf("[DEBUG][handleFirstTimeSetting] This is an UPDATE scenario (oldVal='%s'), checking if new value exists in remote state", oldVal)
	return checkElementInRemoteState(baseField, newVal, d)
}

// isOriginEmpty checks if origin value is effectively empty
func isOriginEmpty(originVal interface{}) bool {
	if originVal == nil {
		return true
	}

	switch v := originVal.(type) {
	case []interface{}:
		return len(v) == 0
	case *schema.Set:
		return v.Len() == 0
	default:
		log.Printf("[DEBUG][isOriginEmpty] Unexpected originVal type: %T", originVal)
		return true
	}
}

// isElementInOrigin checks if an element exists in origin value
func isElementInOrigin(element string, originVal interface{}) bool {
	if originVal == nil {
		return false
	}

	switch v := originVal.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str == element {
				return true
			}
		}
	case *schema.Set:
		return v.Contains(element)
	}

	return false
}

// diffStrSliceParam is used to check whether the parameters of the current string slice type have been modified
// other than those changed in the console.
// The following scenarios will determine whether the parameter has changed (method return false):
//  1. The new value of the script adds new elements compared to the console value (locally added elements).
//  2. The new value of the script has elements decreased compared to the origin value (locally removed elements).
//
// The following scenarios will suppress the diff (method return true):
//  1. The new value of the script is a subset of the console value AND
//     the new value of the script has no elements decreased compared to the origin value.
//
// Examples:
//
// Origin value: ["a", "b", "c"]
// Console value: ["a", "b", "c", "d"] (remotely added "d")
//
// Example 1 (Return false - locally added new element):
//
//	Script value: ["a", "b", "c", "e"] -> Contains "e" not in console (locally added)
//
// Example 2 (Return false - locally removed element):
//
//	Script value: ["a", "b"] -> Removed "c" from origin (locally removed)
//
// Example 3 (Return true - subset of console and no decrease from origin):
//
//	Script value: ["a", "b", "c"] -> Subset of console, same as origin
//
// Example 4 (Return true - subset of console and no decrease from origin):
//
//	Script value: ["a", "b"] -> Subset of console, subset of origin (allowed decrease)
//
// Example 5 (Return false - locally removed element, even if console has new elements):
//
//	Origin: ["a", "b", "c"]
//	Console: ["a", "b", "c", "d", "e"] (remotely added "d", "e")
//	Script: ["a", "b"] -> Removed "c" from origin (locally removed)
//
// Example 6 (Return true - ignore remotely added elements, no local changes):
//
//	Origin: ["a", "b", "c"]
//	Console: ["a", "b", "c", "d", "e"] (remotely added "d", "e")
//	Script: ["a", "b", "c"] -> Same as origin, ignore remote additions
func diffStrSliceParam(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	var (
		originSlice, consoleSlice, newScriptSlice []string
		originParamKey                            = fmt.Sprintf("%s_origin", paramKey)
	)

	// Get the origin value (last local configuration) - this is the key fix
	originVal := d.Get(originParamKey)
	if originVal != nil {
		// Handle different types that originVal might be
		switch v := originVal.(type) {
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok && str != "" {
					originSlice = append(originSlice, str)
				}
			}
		case *schema.Set:
			for _, item := range v.List() {
				if str, ok := item.(string); ok && str != "" {
					originSlice = append(originSlice, str)
				}
			}
		case string:
			if v != "" {
				originSlice = strings.Split(v, ",")
				originSlice = removeEmptyStrings(originSlice)
			}
		default:
			log.Printf("[DEBUG][diffStrSliceParam] Unexpected originVal type: %T", originVal)
		}
	}

	// If origin is empty, this is the first time setting the value
	// In this case, we should allow the change (not suppress diff)
	if len(originSlice) == 0 {
		log.Printf("[DEBUG][diffStrSliceParam] Origin is empty, allowing change (first time setting)")
		return false
	}

	// Parse the old and new values from GetChange
	// oldVal and newVal are already strings from Terraform's diff suppression
	// They represent the serialized form of the lists
	if oldVal != "" {
		consoleSlice = strings.Split(oldVal, ",")
		consoleSlice = removeEmptyStrings(consoleSlice)
	}
	if newVal != "" {
		newScriptSlice = strings.Split(newVal, ",")
		newScriptSlice = removeEmptyStrings(newScriptSlice)
	}

	log.Printf("[DEBUG][diffStrSliceParam] paramKey='%s', originSlice=%v, consoleSlice=%v, newScriptSlice=%v",
		paramKey, originSlice, consoleSlice, newScriptSlice)

	// Check if only care about elements that are in new script but NOT in console (locally added)
	// This means we ignore elements that are in console but NOT in new script (remotely added)
	localAdditions := FindStrSliceElementsNotInAnother(newScriptSlice, consoleSlice)
	if len(localAdditions) > 0 {
		log.Printf("[DEBUG][diffStrSliceParam] New script contains elements not in console (locally added): %v, not suppressing diff",
			localAdditions)
		return false
	}

	// Check if new script has elements decreased compared to origin (locally removed)
	// These are elements that will be deleted from remote
	log.Printf("[DEBUG][diffStrSliceParam] comparing newScriptSlice=%v with originSlice=%v", newScriptSlice, originSlice)
	localRemovals := FindStrSliceElementsNotInAnother(originSlice, newScriptSlice)
	if len(localRemovals) > 0 {
		log.Printf("[DEBUG][diffStrSliceParam] New script has elements decreased compared to origin (locally removed), not suppressing diff")
		return false
	}

	// Both conditions are met, suppress the diff
	log.Printf("[DEBUG][diffStrSliceParam] No local additions and no local removals, suppressing diff")
	return true
}

// removeEmptyStrings removes empty strings from a slice
func removeEmptyStrings(slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, strings.TrimSpace(s))
		}
	}
	return result
}

// checkElementInRemoteState checks if an element exists in remote state
func checkElementInRemoteState(baseField, elementValue string, d *schema.ResourceData) bool {
	// Get the current remote state value
	currentVal := d.Get(baseField)
	if currentVal != nil {
		switch v := currentVal.(type) {
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok && str == elementValue {
					log.Printf("[DEBUG][checkElementInRemoteState] Element '%s' already exists in remote state, suppressing diff",
						elementValue)
					return true
				}
			}
		case *schema.Set:
			if v.Contains(elementValue) {
				log.Printf("[DEBUG][checkElementInRemoteState] Element '%s' already exists in remote state, suppressing diff",
					elementValue)
				return true
			}
		}
	}

	log.Printf("[DEBUG][checkElementInRemoteState] Element '%s' does not exist in remote state, allowing change",
		elementValue)
	return false
}

package utils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetNestedObjectFromRawConfig extracts a nested object using a dot-separated path
// and returns it as a Go interface{} type. It supports deep nesting with both
// object properties and list indices.
func GetNestedObjectFromRawConfig(rawConfig cty.Value, path string) interface{} {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}

	// If path is empty, return the entire object
	if path == "" {
		return getObjectFromRawConfig(rawConfig)
	}

	paths := strings.Split(path, ".")
	return getNestedObject(rawConfig, paths)
}

// getObjectFromRawConfig recursively extracts the entire object from rawConfig
// and converts it to a Go interface{} type that can be used in Terraform operations.
// It supports all cty types including nested structures, lists, maps, and primitive types.
func getObjectFromRawConfig(rawConfig cty.Value) interface{} {
	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil
	}

	ty := rawConfig.Type()
	switch {
	case ty == cty.String:
		return rawConfig.AsString()
	case ty == cty.Number:
		// For numbers, convert to float64 for consistency
		bigFloat := rawConfig.AsBigFloat()
		f64, _ := bigFloat.Float64()
		return f64
	case ty == cty.Bool:
		return rawConfig.True()
	case ty.IsListType() || ty.IsSetType() || ty.IsTupleType():
		return getListFromRawConfig(rawConfig)
	case ty.IsMapType() || ty.IsObjectType():
		return getMapFromRawConfig(rawConfig)
	default:
		// For unknown types, return the raw value
		return rawConfig
	}
}

// getListFromRawConfig handles list, set, and tuple types
func getListFromRawConfig(listValue cty.Value) []interface{} {
	if listValue.IsNull() || !listValue.IsKnown() {
		return nil
	}

	if !listValue.CanIterateElements() {
		return nil
	}

	var result []interface{}
	it := listValue.ElementIterator()
	for it.Next() {
		_, val := it.Element()
		converted := getObjectFromRawConfig(val)
		result = append(result, converted)
	}

	return result
}

// getMapFromRawConfig handles map and object types
func getMapFromRawConfig(mapValue cty.Value) map[string]interface{} {
	if mapValue.IsNull() || !mapValue.IsKnown() {
		return nil
	}

	if !mapValue.CanIterateElements() {
		return nil
	}

	result := make(map[string]interface{})
	it := mapValue.ElementIterator()
	for it.Next() {
		key, val := it.Element()
		keyStr := key.AsString()
		converted := getObjectFromRawConfig(val)
		result[keyStr] = converted
	}

	return result
}

// getNestedObject recursively navigates through the object structure
// using the provided path segments
func getNestedObject(obj cty.Value, paths []string) interface{} {
	if len(paths) == 0 {
		return getObjectFromRawConfig(obj)
	}

	if obj.IsNull() || !obj.IsKnown() {
		return nil
	}

	currentPath := paths[0]
	remainingPaths := paths[1:]

	ty := obj.Type()
	switch {
	case ty.IsObjectType():
		if !obj.Type().HasAttribute(currentPath) {
			return nil
		}
		nextObj := obj.GetAttr(currentPath)
		return getNestedObject(nextObj, remainingPaths)

	case ty.IsListType() || ty.IsSetType() || ty.IsTupleType():
		// Handle list indexing
		if index, err := strconv.Atoi(currentPath); err == nil {
			if !obj.CanIterateElements() {
				return nil
			}
			it := obj.ElementIterator()
			var targetValue cty.Value
			currentIndex := 0
			for it.Next() {
				if currentIndex == index {
					_, targetValue = it.Element()
					break
				}
				currentIndex++
			}
			if targetValue.IsNull() {
				return nil
			}
			return getNestedObject(targetValue, remainingPaths)
		}
		// If not a valid index, treat as a property (might be for tuple types)
		return nil

	case ty.IsMapType():
		if !obj.CanIterateElements() {
			return nil
		}
		it := obj.ElementIterator()
		for it.Next() {
			key, val := it.Element()
			if key.AsString() == currentPath {
				return getNestedObject(val, remainingPaths)
			}
		}
		return nil

	default:
		// For primitive types, if there are remaining paths, return nil
		// as we can't navigate further
		if len(remainingPaths) > 0 {
			return nil
		}
		return getObjectFromRawConfig(obj)
	}
}

// RefreshObjectParamOriginValues updates origin values after all diff calculations are completed.
// This function captures the final configuration values that will be used for comparison in DiffSuppressFunc.
// It handles both direct field changes and length changes (e.g., lts_custom_tag_origin.%).
// Origin values are used to store the current configuration state for subsequent diff suppression.
func RefreshObjectParamOriginValues(d *schema.ResourceData, objectParamKeys []string) error {
	log.Printf("[DEBUG][RefreshObjectParamOriginValues] Starting with %d object param keys: %v",
		len(objectParamKeys), objectParamKeys)

	rawConfig := d.GetRawConfig()
	for _, absParamKeyPath := range objectParamKeys {
		// Construct the corresponding _origin path.
		absOriginParamKeyPath := fmt.Sprintf("%s_origin", absParamKeyPath)
		log.Printf("[DEBUG][RefreshObjectParamOriginValues] Processing '%s' -> '%s'",
			absParamKeyPath, absOriginParamKeyPath)

		// Get current configuration value from rawConfig
		rawVal := GetNestedObjectFromRawConfig(rawConfig, absParamKeyPath)

		if rawVal == nil {
			log.Printf("[DEBUG] Failed to get origin value for the parameter '%s'", absParamKeyPath)
			// If the acquisition fails, the subsequent operation of the current parameter is skipped because this
			// parameter may not be configured.
			continue
		}

		// Set the origin value to match the configuration
		log.Printf("[DEBUG][RefreshObjectParamOriginValues] Setting origin value for '%s'", absOriginParamKeyPath)

		// Set the actual value using setNestedValueSafelyForResourceData to ensure nested safety
		if err := setNestedValueSafelyForResourceData(d, absOriginParamKeyPath, rawVal); err != nil {
			log.Printf("[ERROR][RefreshObjectParamOriginValues] Failed to set origin value for '%s': %v",
				absOriginParamKeyPath, err)
			return fmt.Errorf("failed to set origin value for '%s': %v", absOriginParamKeyPath, err)
		}

		log.Printf("[DEBUG][RefreshObjectParamOriginValues] Successfully set origin value for '%s'",
			absOriginParamKeyPath)
	}

	return nil
}

// setNestedValueSafelyForResourceData safely sets nested values in a ResourceData
func setNestedValueSafelyForResourceData(d *schema.ResourceData, absOriginParamKeyPath string, value interface{}) error {
	parts := strings.Split(absOriginParamKeyPath, ".")
	rootKey := parts[0]

	// Get current value and create a deep copy to avoid affecting other references
	current := d.Get(rootKey)
	currentCopy := deepCopyInterface(current)

	// Update the copy
	updated, err := updateNestedStructureSafely(currentCopy, parts[1:], value)
	if err != nil {
		return err
	}

	// Set the updated value
	// lintignore:R001
	err = d.Set(rootKey, updated)
	if err != nil {
		return err
	}
	return nil
}

// deepCopyInterface creates a deep copy of an interface{} value
func deepCopyInterface(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, val := range v {
			result[key] = deepCopyInterface(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, val := range v {
			result[i] = deepCopyInterface(val)
		}
		return result
	default:
		// For primitive types, return as is
		return value
	}
}

// updateNestedStructureSafely safely updates nested structure without affecting other parts
func updateNestedStructureSafely(current interface{}, parts []string, value interface{}) (interface{}, error) {
	if len(parts) == 0 {
		return value, nil
	}

	part := parts[0]
	switch cv := current.(type) {
	case []interface{}:
		if len(cv) == 0 {
			return nil, errors.New("cannot update empty list")
		}

		// Handle list indexing
		if index, err := strconv.Atoi(part); err == nil {
			if index < 0 || index >= len(cv) {
				return nil, fmt.Errorf("index %d out of range for list of length %d", index, len(cv))
			}

			// Create a copy of the slice
			result := make([]interface{}, len(cv))
			copy(result, cv)

			// Update the specific index
			updatedElem, err := updateNestedStructureSafely(result[index], parts[1:], value)
			if err != nil {
				return nil, err
			}
			result[index] = updatedElem
			return result, nil
		}

		// If not a valid index, treat as property access (for tuple types)
		return nil, fmt.Errorf("invalid list index: %s", part)

	case map[string]interface{}:
		// Create a copy of the map
		result := make(map[string]interface{})
		for key, val := range cv {
			result[key] = deepCopyInterface(val)
		}

		// Check if the key exists
		subVal, ok := result[part]
		if !ok {
			return nil, fmt.Errorf("the parameter key '%s' not found", part)
		}

		// Update the specific key
		updatedSubVal, err := updateNestedStructureSafely(subVal, parts[1:], value)
		if err != nil {
			return nil, err
		}
		result[part] = updatedSubVal
		return result, nil

	default:
		return nil, fmt.Errorf("unsupported type at '%s': %T", part, current)
	}
}

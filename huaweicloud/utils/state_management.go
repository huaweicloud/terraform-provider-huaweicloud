package utils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-multierror"
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

func RefreshObjectParamOriginValues(d *schema.ResourceData, objectParamKeys []string) error {
	var mErr *multierror.Error

	for _, key := range objectParamKeys {
		parts := strings.Split(key, ".")
		// Construct the corresponding _origin path.
		originParts := make([]string, len(parts))
		copy(originParts, parts)
		lastIdx := len(originParts) - 1
		originParts[lastIdx] += "_origin"

		// Obtain the origin value
		rawVal, err := getNestedValue(d, parts)
		if err != nil {
			log.Printf("[DEBUG] failed to get origin value for the parameter '%s': %v", key, err)
			// If the acquisition fails, the subsequent operation of the current parameter is skipped because this
			// parameter may not be configured.
			continue
		}

		// Setting the origin value
		if err := setNestedValue(d, originParts, rawVal); err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("failed to set origin value for '%s': %v", key, err))
		}
	}

	return mErr.ErrorOrNil()
}

// getNestedValue method that used to obtain nested values ​​based on the path recursively, because the nested parameter
// must ensure that the complete structure nesting of its corresponding subscript is obtained (only the corresponding
// index is covered)
func getNestedValue(d *schema.ResourceData, parts []string) (interface{}, error) {
	var current interface{}
	current = d.Get(parts[0])

	for i := 1; i < len(parts); i++ {
		part := parts[i]
		switch cv := current.(type) {
		case []interface{}:
			if len(cv) == 0 {
				return nil, fmt.Errorf("empty list at '%s'", strings.Join(parts[:i+1], "."))
			}
			// Processing lists/collections (automatically taking the first element if the index number is missing).
			current = cv[0]
			if index, err := strconv.Atoi(part); err == nil {
				if index >= len(cv) {
					return nil, fmt.Errorf("index %d out of range", index)
				}
				current = cv[index]
			} else {
				elem, ok := current.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("invalid nested path at '%s'", strings.Join(parts[:i+1], "."))
				}
				current = elem[part]
			}
		case map[string]interface{}:
			var ok bool
			current, ok = cv[part]
			if !ok {
				return nil, fmt.Errorf("missing key '%s'", part)
			}
		default:
			return nil, fmt.Errorf("unsupported type at '%s'", strings.Join(parts[:i+1], "."))
		}
	}
	return current, nil
}

// setNestedValue method that used to set nested value recursively, because nested parameters must set their full
// structure nesting according to their index (only overwrite the corresponding index).
func setNestedValue(d *schema.ResourceData, parts []string, value interface{}) error {
	rootKey := parts[0]
	current := d.Get(rootKey)

	updated, err := updateNestedStructure(current, parts[1:], value)
	if err != nil {
		return err
	}

	// lintignore:R001
	return d.Set(rootKey, updated)
}

func updateNestedStructure(current interface{}, parts []string, value interface{}) (interface{}, error) {
	if len(parts) == 0 {
		return value, nil
	}

	part := parts[0]
	switch cv := current.(type) {
	case []interface{}:
		if len(cv) == 0 {
			return nil, errors.New("cannot update empty list")
		}
		// Considering that the index of the Set type is inconsistent during the change before and after, currently only
		// the first element of the List type is automatically processed (applicable to the MaxItems=1 scenario).
		updatedElem, err := updateNestedStructure(cv[0], parts[1:], value)
		if err != nil {
			return nil, err
		}
		cv[0] = updatedElem
		return cv, nil
	case map[string]interface{}:
		subVal, ok := cv[part]
		if !ok {
			return nil, fmt.Errorf("the parameter key '%s' not found", part)
		}
		updatedSubVal, err := updateNestedStructure(subVal, parts[1:], value)
		if err != nil {
			return nil, err
		}
		cv[part] = updatedSubVal
		return cv, nil
	default:
		return nil, fmt.Errorf("unsupported type at '%s'", part)
	}
}

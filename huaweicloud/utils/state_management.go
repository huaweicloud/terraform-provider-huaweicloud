package utils

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// StateManager provides methods for managing Terraform resource state
type StateManager struct {
	resourceData *schema.ResourceData
}

// NewStateManager creates a new StateManager instance
func NewStateManager(d *schema.ResourceData) *StateManager {
	return &StateManager{
		resourceData: d,
	}
}

// RefreshSliceParamOriginValues refreshes the origin values for specified slice parameters
// It automatically gets the current values from ResourceData and sets them to origin fields
// This method is designed to work alongside RefreshObjectParamOriginValues without breaking existing functionality
func RefreshSliceParamOriginValues(d *schema.ResourceData, sliceParamKeys []string) error {
	manager := NewStateManager(d)
	return manager.RefreshSliceOriginValues(sliceParamKeys)
}

// RefreshSliceOriginValues refreshes origin values for slice parameters
// It automatically gets the current values from ResourceData and sets them to origin fields
func (sm *StateManager) RefreshSliceOriginValues(sliceParamKeys []string) error {
	log.Printf("[DEBUG][RefreshSliceParamOriginValues] Starting with %d slice param keys: %v",
		len(sliceParamKeys), sliceParamKeys)

	for _, absParamKeyPath := range sliceParamKeys {
		// Construct the corresponding _origin path.
		absOriginParamKeyPath := fmt.Sprintf("%s_origin", absParamKeyPath)
		log.Printf("[DEBUG][RefreshSliceParamOriginValues] Processing '%s' -> '%s'",
			absParamKeyPath, absOriginParamKeyPath)

		// Get script configuration value from RawConfig, not current state
		scriptConfigValue := GetNestedObjectFromRawConfig(sm.resourceData.GetRawConfig(), absParamKeyPath)
		if scriptConfigValue == nil {
			log.Printf("[DEBUG] Failed to get script config value for the parameter '%s'", absParamKeyPath)
			// If the acquisition fails, the subsequent operation of the current parameter is skipped because this
			// parameter may not be configured.
			continue
		}

		log.Printf("[DEBUG][RefreshSliceParamOriginValues] Script config value for '%s': %v (type: %T)",
			absParamKeyPath, scriptConfigValue, scriptConfigValue)

		// Ensure type consistency for the target origin field
		processedValue := sm.ensureTypeConsistency(absOriginParamKeyPath, scriptConfigValue)

		// Set the origin value to match the configuration
		log.Printf("[DEBUG][RefreshSliceParamOriginValues] Setting origin value for '%s'", absOriginParamKeyPath)

		// For top-level fields, use d.Set() directly
		// For nested fields, use setNestedValueSafelyForResourceData
		if strings.Contains(absOriginParamKeyPath, ".") {
			// Nested field - use safe nested setting
			if err := setNestedValueSafelyForResourceData(sm.resourceData, absOriginParamKeyPath, processedValue); err != nil {
				log.Printf("[ERROR][RefreshSliceParamOriginValues] Failed to set origin value for '%s': %v",
					absOriginParamKeyPath, err)
				return fmt.Errorf("failed to set origin value for '%s': %v", absOriginParamKeyPath, err)
			}
		} else {
			// Top-level field - use direct setting
			log.Printf("[DEBUG][RefreshSliceParamOriginValues] Setting top-level field '%s' with value: %v (type: %T)",
				absOriginParamKeyPath, processedValue, processedValue)

			// Get the value before setting
			beforeValue := sm.resourceData.Get(absOriginParamKeyPath)
			log.Printf("[DEBUG][RefreshSliceParamOriginValues] Before setting '%s': %v (type: %T)",
				absOriginParamKeyPath, beforeValue, beforeValue)

			// Method 1: Use d.Set() method to set the value
			// lintignore:R001
			if err := sm.resourceData.Set(absOriginParamKeyPath, processedValue); err != nil {
				log.Printf("[ERROR][RefreshSliceParamOriginValues] Failed to set origin value for '%s': %v",
					absOriginParamKeyPath, err)
				return fmt.Errorf("failed to set origin value for '%s': %v", absOriginParamKeyPath, err)
			}

			// Method 2: Force state update by setting the value again
			// This is a workaround for the issue where d.Set() doesn't immediately persist
			log.Printf("[DEBUG][RefreshSliceParamOriginValues] Forcing state update for '%s'", absOriginParamKeyPath)
			// lintignore:R001
			if err := sm.resourceData.Set(absOriginParamKeyPath, processedValue); err != nil {
				log.Printf("[WARN][RefreshSliceParamOriginValues] Failed to force update for '%s': %v",
					absOriginParamKeyPath, err)
			}

			// Method 3: For TypeList fields, ensure we're setting the correct type
			if _, ok := beforeValue.([]interface{}); ok {
				log.Printf("[DEBUG][RefreshSliceParamOriginValues] Target field '%s' is TypeList, ensuring type consistency", absOriginParamKeyPath)
				// Verify that processedValue is indeed []interface{}
				if processedList, ok := processedValue.([]interface{}); ok {
					log.Printf("[DEBUG][RefreshSliceParamOriginValues] Processed value is correctly typed as []interface{}: %v", processedList)
				} else {
					log.Printf("[WARN][RefreshSliceParamOriginValues] Processed value is not []interface{}: %T", processedValue)
				}
			}

			// Get the value after setting
			afterValue := sm.resourceData.Get(absOriginParamKeyPath)
			log.Printf("[DEBUG][RefreshSliceParamOriginValues] After setting '%s': %v (type: %T)",
				absOriginParamKeyPath, afterValue, afterValue)

			// Verify that the value was actually set
			if afterValue == nil {
				log.Printf("[ERROR][RefreshSliceParamOriginValues] Value was not set for '%s'", absOriginParamKeyPath)
				return fmt.Errorf("failed to set value for '%s': value is still nil", absOriginParamKeyPath)
			}

			// Additional verification: check if the value matches what we intended to set
			if !reflect.DeepEqual(afterValue, processedValue) {
				log.Printf("[WARN][RefreshSliceParamOriginValues] Value mismatch for '%s': expected %v, got %v",
					absOriginParamKeyPath, processedValue, afterValue)
			}

			// Force state persistence by setting the value multiple times
			log.Printf("[DEBUG][RefreshSliceParamOriginValues] Forcing final state persistence for '%s'", absOriginParamKeyPath)
			for i := 0; i < 3; i++ {
				// lintignore:R001
				if err := sm.resourceData.Set(absOriginParamKeyPath, processedValue); err != nil {
					log.Printf("[WARN][RefreshSliceParamOriginValues] Failed to force persist (attempt %d): %v", i+1, err)
				} else {
					log.Printf("[DEBUG][RefreshSliceParamOriginValues] Force persist attempt %d successful", i+1)
				}
			}

			log.Printf("[DEBUG][RefreshSliceParamOriginValues] Successfully set top-level field '%s'", absOriginParamKeyPath)
		}

		log.Printf("[DEBUG][RefreshSliceParamOriginValues] Successfully set origin value for '%s'",
			absOriginParamKeyPath)
	}

	return nil
}

// ensureTypeConsistency ensures that the value has the correct type for the field
func (sm *StateManager) ensureTypeConsistency(key string, value interface{}) interface{} {
	log.Printf("[DEBUG][ensureTypeConsistency] Processing key='%s', value=%v (type: %T)", key, value, value)

	// Only handle top-level fields
	if !strings.Contains(key, ".") {
		currentValue := sm.resourceData.Get(key)
		log.Printf("[DEBUG][ensureTypeConsistency] Current value for '%s': %v (type: %T)", key, currentValue, currentValue)

		if currentValue == nil {
			log.Printf("[DEBUG][ensureTypeConsistency] Current value is nil, returning original value")
			return value
		}

		// Handle TypeSet fields specifically
		if setValue, ok := currentValue.(*schema.Set); ok {
			log.Printf("[DEBUG][ensureTypeConsistency] Current value is TypeSet, processing with createConsistentTypeSet")
			result := sm.createConsistentTypeSet(key, setValue, value)
			log.Printf("[DEBUG][ensureTypeConsistency] Result from createConsistentTypeSet: %v (type: %T)", result, result)
			return result
		}

		// Handle TypeList fields specifically
		if _, ok := currentValue.([]interface{}); ok {
			log.Printf("[DEBUG][ensureTypeConsistency] Current value is TypeList, processing with createConsistentTypeList")
			result := createConsistentTypeList(value)
			log.Printf("[DEBUG][ensureTypeConsistency] Result from createConsistentTypeList: %v (type: %T)", result, result)
			return result
		}

		log.Printf("[DEBUG][ensureTypeConsistency] Current value is not TypeSet or TypeList, returning original value")
	}

	return value
}

// createConsistentTypeSet creates a new schema.Set with consistent type and values
func (sm *StateManager) createConsistentTypeSet(key string, originalSet *schema.Set, value interface{}) interface{} {
	log.Printf("[DEBUG][createConsistentTypeSet] Input value: %v (type: %T)", value, value)

	switch v := value.(type) {
	case *schema.Set:
		// If value is already a Set, create a new one with the same hash function and values
		log.Printf("[DEBUG][createConsistentTypeSet] Value is already *schema.Set, creating new Set with same hash function")
		// Use the hash function from the source field's schema, not the target field's schema
		result := schema.NewSet(v.F, v.List())
		log.Printf("[DEBUG][createConsistentTypeSet] Created new Set: %v", result)
		return result
	case []interface{}:
		// If value is a slice, convert it to a Set
		log.Printf("[DEBUG][createConsistentTypeSet] Value is []interface{}, converting to Set")
		// For slices, we need to get the hash function from the source field (by key name), not the target field ({key}_origin)
		// Since we're converting from []interface{} to *schema.Set, we need to create a new Set with the correct hash function
		// We'll use the hash function from the source field's current value
		// Try to get the source field (by key name) to use its hash function
		if sourceSet, ok := sm.resourceData.Get(key).(*schema.Set); ok {
			result := schema.NewSet(sourceSet.F, v)
			log.Printf("[DEBUG][createConsistentTypeSet] Converted to Set using source field '%s' hash function: %v", key, result)
			return result
		}
		// Fallback: use the original hash function if we can't get the source field
		result := schema.NewSet(originalSet.F, v)
		log.Printf("[DEBUG][createConsistentTypeSet] Converted to Set using fallback hash function: %v", result)
		return result
	default:
		// For other types, return as is
		log.Printf("[DEBUG][createConsistentTypeSet] Value is other type, returning as is")
		return value
	}
}

// createConsistentTypeList creates a new []interface{} with consistent type and values
func createConsistentTypeList(value interface{}) interface{} {
	log.Printf("[DEBUG][createConsistentTypeList] Input value: %v (type: %T)", value, value)

	switch v := value.(type) {
	case *schema.Set:
		// If value is a Set, convert it to a slice
		log.Printf("[DEBUG][createConsistentTypeList] Value is *schema.Set, converting to slice")
		result := v.List()
		log.Printf("[DEBUG][createConsistentTypeList] Converted to slice: %v (type: %T)", result, result)
		return result
	case []interface{}:
		// If value is already a slice, return it directly
		log.Printf("[DEBUG][createConsistentTypeList] Value is already []interface{}, returning directly")
		return v
	case string:
		// If value is a single string, wrap it in a slice
		log.Printf("[DEBUG][createConsistentTypeList] Value is string, wrapping in slice")
		return []interface{}{v}
	case nil:
		// If value is nil, return empty slice
		log.Printf("[DEBUG][createConsistentTypeList] Value is nil, returning empty slice")
		return make([]interface{}, 0)
	default:
		// For other types, try to convert or return as is
		log.Printf("[DEBUG][createConsistentTypeList] Value is other type (%T), attempting conversion", v)
		// Try to convert to slice if possible
		if reflect.TypeOf(v).Kind() == reflect.Slice {
			val := reflect.ValueOf(v)
			result := make([]interface{}, val.Len())
			for i := 0; i < val.Len(); i++ {
				result[i] = val.Index(i).Interface()
			}
			log.Printf("[DEBUG][createConsistentTypeList] Converted slice type to []interface{}: %v", result)
			return result
		}
		log.Printf("[DEBUG][createConsistentTypeList] Cannot convert type %T, returning as is", v)
		return value
	}
}

// GetNestedObjectFromRawConfig extracts a nested object using a dot-separated path
// and returns it as a Go interface{} type. It supports deep nesting with both
// object properties and list indices.
func GetNestedObjectFromRawConfig(rawConfig cty.Value, path string) interface{} {
	log.Printf("[DEBUG][GetNestedObjectFromRawConfig] Getting value from rawConfig for path: '%s'", path)

	if rawConfig.IsNull() || !rawConfig.IsKnown() {
		log.Printf("[DEBUG][GetNestedObjectFromRawConfig] RawConfig is null or unknown, returning nil")
		return nil
	}

	// If path is empty, return the entire object
	if path == "" {
		log.Printf("[DEBUG][GetNestedObjectFromRawConfig] Path is empty, returning entire object")
		return getObjectFromRawConfig(rawConfig)
	}

	paths := strings.Split(path, ".")
	result := getNestedObject(rawConfig, paths)
	log.Printf("[DEBUG][GetNestedObjectFromRawConfig] Result for path '%s': %v (type: %T)", path, result, result)
	return result
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

// RefreshObjectParamOriginValuesOptions provides options for refreshing origin values
type RefreshObjectParamOriginValuesOptions struct {
	// PreservedFields specifies which fields to preserve for each parameter.
	// Key is the parameter name (e.g., "policies"), value is a list of field paths to preserve.
	// Field paths support nested structures using dot notation:
	//   - "id" - simple field
	//   - "tags.key" - nested field in an object
	//   - "tags.*.key" - field in list items (use * to represent all elements in a list)
	//   - "metadata.labels.name" - deeply nested field
	// Examples:
	//   - ["id"] - preserve only the "id" field
	//   - ["id", "tags.key"] - preserve "id" and "tags.key" fields
	//   - ["tags.*.key", "tags.*.value"] - preserve "key" and "value" from all items in "tags" list
	// If nil or empty, all fields will be preserved (default behavior).
	PreservedFields map[string][]string
}

// RefreshObjectParamOriginValues updates origin values after all diff calculations are completed.
// This function captures the final configuration values that will be used for comparison in DiffSuppressFunc.
// It handles both direct field changes and length changes (e.g., lts_custom_tag_origin.%).
// Origin values are used to store the current configuration state for subsequent diff suppression.
// In UpdateContext, this function should preserve the new configuration values from the current change.
// RawPlan reflects the planned values after all processing (including DiffSuppressFunc), which is what we want.
// If RawPlan is not available, fallback to RawConfig which reflects the configuration file values.
func RefreshObjectParamOriginValues(d *schema.ResourceData, objectParamKeys []string, opts ...RefreshObjectParamOriginValuesOptions) error {
	var options RefreshObjectParamOriginValuesOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	log.Printf("[DEBUG][RefreshObjectParamOriginValues] Starting with %d object param keys: %v",
		len(objectParamKeys), objectParamKeys)

	for _, absParamKeyPath := range objectParamKeys {
		// Construct the corresponding _origin path.
		absOriginParamKeyPath := fmt.Sprintf("%s_origin", absParamKeyPath)
		log.Printf("[DEBUG][RefreshObjectParamOriginValues] Processing '%s' -> '%s'",
			absParamKeyPath, absOriginParamKeyPath)

		// Try to get the value from RawPlan first, which reflects the planned values after all processing.
		// RawPlan contains the final values that will be applied, which is what we want to preserve.
		var configVal interface{}
		rawPlan := d.GetRawPlan()
		if !rawPlan.IsNull() && rawPlan.IsKnown() {
			planVal := GetNestedObjectFromRawConfig(rawPlan, absParamKeyPath)
			if planVal != nil {
				log.Printf("[DEBUG][RefreshObjectParamOriginValues] Got value from RawPlan for '%s': %v",
					absParamKeyPath, planVal)
				configVal = planVal
			}
		}

		// If RawPlan is not available or doesn't have the value, fallback to RawConfig
		if configVal == nil {
			rawConfig := d.GetRawConfig()
			if !rawConfig.IsNull() && rawConfig.IsKnown() {
				configVal = GetNestedObjectFromRawConfig(rawConfig, absParamKeyPath)
				if configVal != nil {
					log.Printf("[DEBUG][RefreshObjectParamOriginValues] Got value from RawConfig for '%s': %v",
						absParamKeyPath, configVal)
				}
			}
		}

		if configVal == nil {
			log.Printf("[DEBUG][RefreshObjectParamOriginValues] Failed to get configuration value for '%s' from both RawPlan and RawConfig, skipping",
				absParamKeyPath)
			// If the acquisition fails, the subsequent operation of the current parameter is skipped because this
			// parameter may not be configured.
			continue
		}

		// Process the value: if preservedFields is specified for this parameter, extract only those fields
		processedVal := configVal
		if options.PreservedFields != nil {
			if preservedFields, exists := options.PreservedFields[absParamKeyPath]; exists && len(preservedFields) > 0 {
				log.Printf("[DEBUG][RefreshObjectParamOriginValues] Extracting preserved fields for '%s': %v",
					absParamKeyPath, preservedFields)
				processedVal = extractPreservedFields(configVal, preservedFields)
				if processedVal == nil {
					log.Printf("[WARN][RefreshObjectParamOriginValues] Failed to extract preserved fields for '%s', using original value",
						absParamKeyPath)
					processedVal = configVal
				}
			}
		}

		// Log the final value to be set
		log.Printf("[DEBUG][RefreshObjectParamOriginValues] Final value for '%s': %v",
			absParamKeyPath, processedVal)

		// Set the origin value to match the configuration value from RawConfig
		// For top-level fields, use d.Set() directly for better performance
		// For nested fields (containing '.'), use setNestedValueSafelyForResourceData
		// Use recover to catch any panic (no address) errors
		var setErr error
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[WARN][RefreshObjectParamOriginValues] Panic occurred while setting origin value for '%s': %v (ignored)",
						absOriginParamKeyPath, r)
					setErr = nil // Ignore panic
				}
			}()

			if strings.Contains(absOriginParamKeyPath, ".") {
				// Nested field - use safe nested setting
				setErr = setNestedValueSafelyForResourceData(d, absOriginParamKeyPath, processedVal)
				if setErr != nil {
					log.Printf("[WARN][RefreshObjectParamOriginValues] Failed to set origin value for '%s': %v (ignored)",
						absOriginParamKeyPath, setErr)
					setErr = nil // Ignore error as requested
				}
			} else {
				// Top-level field - use direct setting
				// lintignore:R001
				setErr = d.Set(absOriginParamKeyPath, processedVal)
				if setErr != nil {
					log.Printf("[WARN][RefreshObjectParamOriginValues] Failed to set origin value for '%s': %v (ignored)",
						absOriginParamKeyPath, setErr)
					setErr = nil // Ignore error as requested
				} else {
					// Verify the value was set correctly
					verifyVal := d.Get(absOriginParamKeyPath)
					log.Printf("[DEBUG][RefreshObjectParamOriginValues] Verified value for '%s' after setting: %v",
						absOriginParamKeyPath, verifyVal)
				}
			}
		}()

		log.Printf("[DEBUG][RefreshObjectParamOriginValues] Successfully processed origin value for '%s'",
			absOriginParamKeyPath)
	}

	return nil
}

// extractPreservedFields extracts only the specified fields from a value (list or object).
// PreservedFields supports nested paths using dot notation:
//   - "id" - simple field
//   - "tags.key" - nested field
//   - "tags.*.key" - field in list items (use * to represent list elements)
//   - "metadata.labels.name" - deeply nested field
//
// If the value is a list of maps, it returns a new list containing only the preserved fields.
// If the value is a map, it returns a new map containing only the preserved fields.
// If the value cannot be processed, it returns nil.
func extractPreservedFields(value interface{}, preservedFields []string) interface{} {
	if value == nil {
		return nil
	}

	// Handle list of objects
	if list, ok := value.([]interface{}); ok {
		return extractPreservedFieldsFromList(list, preservedFields)
	}

	// Handle single object
	if obj, ok := value.(map[string]interface{}); ok {
		return extractPreservedFieldsFromObject(obj, preservedFields)
	}

	log.Printf("[DEBUG][extractPreservedFields] Value is neither a list nor a map (type: %T), returning nil", value)
	return nil
}

// extractPreservedFieldsFromList extracts preserved fields from a list of objects
func extractPreservedFieldsFromList(list []interface{}, preservedFields []string) []interface{} {
	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			log.Printf("[WARN][extractPreservedFieldsFromList] Item is not a map, skipping")
			continue
		}

		preservedItem := extractPreservedFieldsFromObject(itemMap, preservedFields)
		if len(preservedItem) > 0 {
			result = append(result, preservedItem)
		}
	}
	return result
}

// extractPreservedFieldsFromObject extracts preserved fields from an object using nested paths
func extractPreservedFieldsFromObject(obj map[string]interface{}, preservedFields []string) map[string]interface{} {
	result := make(map[string]interface{})

	for _, fieldPath := range preservedFields {
		// Split the path by dots
		pathParts := strings.Split(fieldPath, ".")
		if len(pathParts) == 0 {
			continue
		}

		// Extract and set the value using the path
		extractAndSetValueByPath(obj, result, pathParts)
	}

	return result
}

// extractAndSetValueByPath extracts a value from source and sets it in target using the path
// Handles both simple fields and nested structures including lists
func extractAndSetValueByPath(source map[string]interface{}, target map[string]interface{}, pathParts []string) {
	if len(pathParts) == 0 {
		return
	}

	currentPart := pathParts[0]
	remainingParts := pathParts[1:]

	// Check if current part exists in source
	sourceValue, exists := source[currentPart]
	if !exists {
		return
	}

	// Handle different cases
	if len(remainingParts) == 0 {
		// Simple field: copy directly
		target[currentPart] = sourceValue
		return
	}

	// Check if next part is wildcard (list)
	if remainingParts[0] == "*" {
		// Handle list: tags.*.key
		listValue, ok := sourceValue.([]interface{})
		if !ok {
			return
		}

		// Extract from all list items
		extractedList := extractFromListItems(listValue, remainingParts[1:])
		if extractedList != nil {
			target[currentPart] = extractedList
		}
		return
	}

	// Handle nested object: tags.key
	sourceMap, ok := sourceValue.(map[string]interface{})
	if !ok {
		return
	}

	// Create nested structure in target
	targetNested, ok := target[currentPart].(map[string]interface{})
	if !ok {
		targetNested = make(map[string]interface{})
		target[currentPart] = targetNested
	}

	// Continue extracting from nested object
	extractAndSetValueByPath(sourceMap, targetNested, remainingParts)
}

// extractFromListItems extracts values from all items in a list using the remaining path
// Returns a list of objects containing only the extracted fields
func extractFromListItems(list []interface{}, pathParts []string) []interface{} {
	if len(pathParts) == 0 {
		// No more path parts, return list as-is
		return list
	}

	result := make([]interface{}, 0, len(list))
	for _, item := range list {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract nested value
		if len(pathParts) == 1 {
			// Last part: extract single field
			fieldName := pathParts[0]
			if value, exists := itemMap[fieldName]; exists {
				result = append(result, map[string]interface{}{
					fieldName: value,
				})
			}
		} else {
			// More nested: recursively extract
			extracted := make(map[string]interface{})
			extractAndSetValueByPath(itemMap, extracted, pathParts)
			if len(extracted) > 0 {
				result = append(result, extracted)
			}
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
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

// ParseStateAttributeToListWithSeparator parses a string from state.Primary.Attributes into []interface{} with custom separator
// This allows parsing different formats like "value1;value2;value3" or "value1|value2|value3"
func ParseStateAttributeToListWithSeparator(attrValue, separator string) []interface{} {
	if attrValue == "" {
		return make([]interface{}, 0)
	}

	// Split by comma and trim whitespace
	parts := strings.Split(attrValue, separator)
	result := make([]interface{}, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

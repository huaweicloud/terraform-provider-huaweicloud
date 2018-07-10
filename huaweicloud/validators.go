package huaweicloud

import (
	"fmt"
	"time"
)

func ValidateStringList(v interface{}, k string, l []string) (ws []string, errors []error) {
	value := v.(string)
	for i := range l {
		if value == l[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, l))
	return
}

func ValidateIntRange(v interface{}, k string, l int, h int) (ws []string, errors []error) {
	i, ok := v.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("%q must be an integer", k))
		return
	}
	if i < l || i > h {
		errors = append(errors, fmt.Errorf("%q must be between %d and %d", k, l, h))
		return
	}
	return
}

func validateS3BucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}

	return
}

func validateTrueOnly(v interface{}, k string) (ws []string, errors []error) {
	if b, ok := v.(bool); ok && b {
		return
	}
	if v, ok := v.(string); ok && v == "true" {
		return
	}
	errors = append(errors, fmt.Errorf("%q must be true", k))
	return
}

func validateS3BucketLifecycleExpirationDays(v interface{}, k string) (ws []string, errors []error) {
	if v.(int) <= 0 {
		errors = append(errors, fmt.Errorf(
			"%q must be greater than 0", k))
	}

	return
}

func validateS3BucketLifecycleRuleId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 255 {
		errors = append(errors, fmt.Errorf(
			"%q cannot exceed 255 characters", k))
	}
	return
}

func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func validateKmsKeyStatus(v interface{}, k string) (ws []string, errors []error) {
	status := v.(string)
	if status != EnabledState && status != DisabledState && status != PendingDeletionState {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid status, expected %s or %s or %s, got %s.",
			k, EnabledState, DisabledState, PendingDeletionState, status))
	}
	return
}

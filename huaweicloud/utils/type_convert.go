package utils

// returns a pointer to the bool value
func Bool(v bool) *bool {
	return &v
}

// returns a pointer to the string value
func String(v string) *string {
	return &v
}

// Int32 returns a pointer to the int32 value
func Int32(v int32) *int32 {
	return &v
}

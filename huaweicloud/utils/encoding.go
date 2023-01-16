package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

// HashAndHexEncode is one of the implementations of SchemaStateFunc.
// The function gets the hash of v then returns the hexadecimal encoding string.
// If the type of v is not string, just returns an empty string.
func HashAndHexEncode(v interface{}) string {
	switch v := v.(type) {
	case string:
		hash := sha1.Sum([]byte(v))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

// DecodeHashAndHexEncode is one of the implementations of SchemaStateFunc.
// The function tries to decode v if it is in base64 format, then gets the hash of
// decode string and returns the hexadecimal encoding string.
// If the type of v is not string, just returns an empty string.
func DecodeHashAndHexEncode(v interface{}) string {
	switch v := v.(type) {
	case string:
		return installScriptHashSum(v)
	default:
		return ""
	}
}

func installScriptHashSum(script string) string {
	// Check whether the script is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(script)
	if base64DecodeError != nil {
		v = []byte(script)
	}

	hash := sha1.Sum(v)
	return hex.EncodeToString(hash[:])
}

// TryBase64EncodeString will encode a string with base64.
// If the string is already base64 encoded, returns it directly.
func TryBase64EncodeString(str string) string {
	if _, err := base64.StdEncoding.DecodeString(str); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	return str
}

// Base64EncodeString is used to encode a string by base64.
func Base64EncodeString(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

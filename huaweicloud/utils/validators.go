package utils

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ValidateSubnetV2IPv6Mode = validation.StringInSlice([]string{
	"slaac", "dhcpv6-stateful", "dhcpv6-stateless",
}, false)

func ValidateTrueOnly(v interface{}, k string) (ws []string, errors []error) {
	if b, ok := v.(bool); ok && b {
		return
	}
	if v, ok := v.(string); ok && v == "true" {
		return
	}
	errors = append(errors, fmt.Errorf("%q must be true", k))
	return
}

func ValidateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := NormalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

// lintignore:V001
func ValidateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

// lintignore:V001
func ValidateString64WithChinese(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}

	pattern := "^[\\-._A-Za-z0-9\u4e00-\u9fa5]+$"
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}

	return
}

func ValidateCIDR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || strings.ToLower(value) != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, got %q", k, value))
	}

	return
}

func ValidateIP(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ipnet := net.ParseIP(value)

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network IP address, got %q", k, value))
	}

	return
}

// lintignore:V001
func ValidateVBSPolicyName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(strings.ToLower(value), "default") {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with default: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

// lintignore:V001
func ValidateVBSTagKey(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 36 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 36 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

// lintignore:V001
func ValidateVBSTagValue(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if len(value) > 43 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 43 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

// lintignore:V001
func ValidateVBSBackupName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if strings.HasPrefix(strings.ToLower(value), "autobk") {
		errors = append(errors, fmt.Errorf(
			"%q cannot start with autobk: %q", k, value))
	}

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}
	pattern := `^[\.\-_A-Za-z0-9]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

// lintignore:V001
func ValidateVBSBackupDescription(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 {
		errors = append(errors, fmt.Errorf(
			"%q cannot be longer than 64 characters: %q", k, value))
	}
	pattern := `^[^<>]+$`
	if !regexp.MustCompile(pattern).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q doesn't comply with restrictions (%q): %q",
			k, pattern, value))
	}
	return
}

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

func ValidateIPRange(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ipAddresses := strings.Split(value, "-")
	if len(ipAddresses) != 2 {
		errors = append(errors, fmt.Errorf(
			"%q must be a valid network IP address range, such as 0.0.0.0-255.255.255.0, but got %q", k, value))
		return
	}
	for _, address := range ipAddresses {
		ipnet := net.ParseIP(address)
		if ipnet == nil || address != ipnet.String() {
			errors = append(errors, fmt.Errorf("%q must contains valid network IP address, got %q", k, address))
		}
	}
	if len(errors) == 0 {
		if ipAddresses[0] == ipAddresses[1] {
			errors = append(errors, fmt.Errorf("Two network IP address of %q cannot equal, got %q", k, value))
		}
		// Split the IP address into a string array for comparison.
		startAddress := strings.Split(ipAddresses[0], ".")
		endAddress := strings.Split(ipAddresses[1], ".")
		// Verify the correctness of the IP address range: The starting IP address must be less than the ending IP address.
		// The For loop compares the four parts of the IPv4 address in turn.
		for i := 0; i < len(startAddress); i++ {
			if startAddress[i] > endAddress[i] {
				errors = append(errors, fmt.Errorf(
					"%q starting IP address cannot be greater than the ending IP address, got %q", k, value))
				return
			} else if startAddress[i] < endAddress[i] {
				return
			}
		}
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

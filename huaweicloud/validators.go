package huaweicloud

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

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

func validateJsonString(v interface{}, k string) (ws []string, errors []error) {
	if _, err := normalizeJsonString(v); err != nil {
		errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
	}
	return
}

func looksLikeJsonString(s interface{}) bool {
	return regexp.MustCompile(`^\s*{`).MatchString(s.(string))
}

func validateStackTemplate(v interface{}, k string) (ws []string, errors []error) {
	if looksLikeJsonString(v) {
		if _, err := normalizeJsonString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
		}
	} else {
		if _, err := checkYamlString(v); err != nil {
			errors = append(errors, fmt.Errorf("%q contains an invalid YAML: %s", k, err))
		}
	}
	return
}

//lintignore:V001
func validateName(v interface{}, k string) (ws []string, errors []error) {
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

//lintignore:V001
func validateString64WithChinese(v interface{}, k string) (ws []string, errors []error) {
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

func validateCIDR(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network CIDR, got %q", k, value))
	}

	return
}

func validateIPRange(v interface{}, k string) (ws []string, errors []error) {
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

func validateIP(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ipnet := net.ParseIP(value)

	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf(
			"%q must contain a valid network IP address, got %q", k, value))
	}

	return
}

func validateSubnetV2IPv6Mode(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value != "slaac" && value != "dhcpv6-stateful" && value != "dhcpv6-stateless" {
		err := fmt.Errorf("%s must be one of slaac, dhcpv6-stateful or dhcpv6-stateless", k)
		errors = append(errors, err)
	}
	return
}

//lintignore:V001
func validateVBSPolicyName(v interface{}, k string) (ws []string, errors []error) {
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

//lintignore:V001
func validateVBSTagKey(v interface{}, k string) (ws []string, errors []error) {
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

//lintignore:V001
func validateVBSTagValue(v interface{}, k string) (ws []string, errors []error) {
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

//lintignore:V001
func validateVBSBackupName(v interface{}, k string) (ws []string, errors []error) {
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

//lintignore:V001
func validateVBSBackupDescription(v interface{}, k string) (ws []string, errors []error) {
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

//lintignore:V001
func validateECSTagValue(v interface{}, k string) (ws []string, errors []error) {
	tagmap := v.(map[string]interface{})
	vv := regexp.MustCompile(`^[0-9a-zA-Z-_]+$`)
	for k, v := range tagmap {
		value := v.(string)
		if !vv.MatchString(value) {
			errors = append(errors, fmt.Errorf("Tag value must be string only contains digits, letters, underscores(_) and hyphens(-), but got %s=%s", k, value))
			break
		}
	}
	return
}

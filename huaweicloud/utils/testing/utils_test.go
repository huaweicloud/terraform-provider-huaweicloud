package testing

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestHasMapContainsFunc_match(t *testing.T) {
	rawMap := map[string]string{
		"foo": "bar",
	}
	filter1 := map[string]interface{}{
		"foo": "bar",
	}
	filter2 := map[string]interface{}{
		"foo": "bar,dor",
	}

	if !utils.HasMapContains(rawMap, filter1) {
		t.Fail()
	}
	if !utils.HasMapContains(rawMap, filter2) {
		t.Fail()
	}
}

func TestHasMapContainsFunc_mismatch(t *testing.T) {
	rawMap := map[string]string{
		"foo": "bar",
	}
	filter1 := map[string]interface{}{
		"foo": "dor",
	}
	filter2 := map[string]interface{}{
		"foo1": "bar",
	}
	filter3 := map[string]interface{}{
		"": "bar",
	}

	if utils.HasMapContains(rawMap, filter1) {
		t.Fail()
	}
	if utils.HasMapContains(rawMap, filter2) {
		t.Fail()
	}
	if utils.HasMapContains(rawMap, filter3) {
		t.Fail()
	}
}

func TestHasMapContainsFunc_empty(t *testing.T) {
	rawMap1 := map[string]string{
		"foo": "bar",
		"key": "",
		"":    "value",
	}
	rawMap2 := map[string]string{}

	filter1 := map[string]interface{}{
		"foo": "",
	}
	filter2 := map[string]interface{}{
		"key": "",
	}
	filter3 := map[string]interface{}{
		"": "value",
	}
	filter4 := map[string]interface{}{}

	if !utils.HasMapContains(rawMap1, filter1) || !utils.HasMapContains(rawMap1, filter2) ||
		!utils.HasMapContains(rawMap1, filter3) || !utils.HasMapContains(rawMap1, filter4) {
		t.Fail()
	}
	if utils.HasMapContains(rawMap2, filter1) || utils.HasMapContains(rawMap2, filter2) ||
		utils.HasMapContains(rawMap2, filter3) || !utils.HasMapContains(rawMap2, filter4) {
		t.Fail()
	}
}

func TestIsIPv4Address(t *testing.T) {
	var validAddresses = []string{
		"10.250.4.192", "192.168.10.1", "192.168.10.255", "255.255.255.255", "0.0.0.0",
	}
	for _, addr := range validAddresses {
		if !utils.IsIPv4Address(addr) {
			t.Fatalf("%s should be a valid IPv4 address", addr)
		}
	}

	var invalidAddresses = []string{
		"258.255.255.300", "192.168.010.76", "a.b.c.def",
	}
	for _, addr := range invalidAddresses {
		if utils.IsIPv4Address(addr) {
			t.Fatalf("%s should not be a valid IPv4 address", addr)
		}
	}

	var ipv6Addresses = []string{
		"2407:c080:1200:e4e:49ad:b067:d81f:6467", "2407:c080:1200:e4e::1", "2002:0:0:0:0:0:9f8a:20c3",
	}
	for _, addr := range ipv6Addresses {
		if utils.IsIPv4Address(addr) {
			t.Fatalf("%s should be a valid IPv6 address", addr)
		}
	}
}

func TestSuppressVersionDiffs(t *testing.T) {
	var validVersions = []string{
		"v1.19", "v1.19.10", "v1.19.10-r0",
	}
	for _, version := range validVersions {
		if !utils.SuppressVersionDiffs("", "v1.19.10-r0", version, nil) {
			t.Fatalf("%s and v1.19.10-r0 should be considered the same version", version)
		}
	}

	var invalidVersions = []string{
		"v1.21", "v1.19.1", "v1.19.10-r1",
	}
	for _, version := range invalidVersions {
		if utils.SuppressVersionDiffs("", "v1.19.10-r0", version, nil) {
			t.Fatalf("%s and v1.19.10-r0 should be considered different versions", version)
		}
	}
}

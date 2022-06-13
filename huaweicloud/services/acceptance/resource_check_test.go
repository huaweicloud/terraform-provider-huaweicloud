package acceptance

import (
	"strings"
	"testing"
)

// AssertEquals compares two arbitrary values and performs a comparison. If the
// comparison fails, a fatal error is raised that will fail the test
func AssertEquals(t *testing.T, actual, expected interface{}) {
	if expected != actual {
		t.Fatalf("expected %s but got %s", expected, actual)
	}
}

// AssertNoErr is a convenience function for checking whether an error value is
// an actual error
func AssertNoErr(t *testing.T, e error) {
	if e != nil {
		t.Fatalf("unexpected error %s", e)
	}
}

func TestResourceCheck_parseVariable(t *testing.T) {
	// resource and data source within huaweicloud
	resource := "${huaweicloud_vpc.test.id}"
	key, field, err := parseVariableToName(resource)
	AssertNoErr(t, err)
	AssertEquals(t, key, "huaweicloud_vpc.test")
	AssertEquals(t, field, "id")

	data := "${data.huaweicloud_compute_flavors.test.ids.0}"
	key, field, err = parseVariableToName(data)
	AssertNoErr(t, err)
	AssertEquals(t, key, "data.huaweicloud_compute_flavors.test")
	AssertEquals(t, field, "ids.0")

	// resource and data source within other clouds
	resource = "${joincloud_vpc.test.id}"
	key, field, err = parseVariableToName(resource)
	AssertNoErr(t, err)
	AssertEquals(t, key, "joincloud_vpc.test")
	AssertEquals(t, field, "id")

	data = "${data.joincloud_compute_flavors.test.ids.0}"
	key, field, err = parseVariableToName(data)
	AssertNoErr(t, err)
	AssertEquals(t, key, "data.joincloud_compute_flavors.test")
	AssertEquals(t, field, "ids.0")

	// error parsing
	resource = "${huaweicloud_vpc.test}"
	_, _, err = parseVariableToName(resource)
	AssertEquals(t, err != nil, true)
	AssertEquals(t, strings.HasPrefix(err.Error(), "attribute field is missing:"), true)
}

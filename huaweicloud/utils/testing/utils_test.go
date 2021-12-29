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

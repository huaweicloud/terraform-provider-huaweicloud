package huaweicloud

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var successHTTPCodes = []int{200, 201, 202, 203, 204, 205, 206, 207, 208, 226}

func isEmptyValue(v reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0, nil
	case reflect.Bool:
		return !v.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0, nil
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0, nil
	case reflect.Interface, reflect.Ptr:
		return v.IsNil(), nil
	case reflect.Invalid:
		return true, nil
	}
	return false, fmt.Errorf("isEmptyValue:: unknown type")
}

func addQueryParams(rawurl string, params map[string]string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func replaceVars(d *schema.ResourceData, linkTmpl string, kv map[string]string) (string, error) {
	re := regexp.MustCompile("{([[:word:]]+)}")

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if kv != nil {
			if v, ok := kv[m]; ok {
				return v
			}
		}
		if m == "project" {
			return "replace_holder"
		}
		if d != nil {
			if m == "id" {
				return d.Id()
			}
			v, ok := d.GetOk(m)
			if ok {
				v1, _ := convertToStr(v)
				return v1
			}
		}
		return ""
	}

	s := re.ReplaceAllStringFunc(linkTmpl, replaceFunc)
	return strings.Replace(s, "replace_holder/", "", 1), nil
}

func replaceVarsForTest(rs *terraform.ResourceState, linkTmpl string) (string, error) {
	re := regexp.MustCompile("{([[:word:]]+)}")

	replaceFunc := func(s string) string {
		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return "replace_holder"
		}
		if rs != nil {
			if m == "id" {
				return rs.Primary.ID
			}
			v, ok := rs.Primary.Attributes[m]
			if ok {
				return v
			}
		}
		return ""
	}

	s := re.ReplaceAllStringFunc(linkTmpl, replaceFunc)
	return strings.Replace(s, "replace_holder/", "", 1), nil
}

func navigateMap(d interface{}, index []string) (interface{}, error) {
	for _, i := range index {
		d1, ok := d.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("navigateMap:: Can not convert to map")
		}
		d, ok = d1[i]
		if !ok {
			return nil, fmt.Errorf("navigateMap:: '%s' may not exist", i)
		}
	}
	return d, nil
}

func navigateValue(d interface{}, index []string, arrayIndex map[string]int) (interface{}, error) {
	for n, i := range index {
		if d == nil {
			return nil, nil
		}
		if d1, ok := d.(map[string]interface{}); ok {
			d, ok = d1[i]
			if !ok {
				return nil, fmt.Errorf("navigate:: '%s' may not exist", i)
			}
		} else {
			return nil, fmt.Errorf("navigateValue:: Can not convert (%s) to map, index=%s", reflect.TypeOf(d), strings.Join(index, "."))
		}

		if arrayIndex != nil {
			if j, ok := arrayIndex[strings.Join(index[:n+1], ".")]; ok {
				if d == nil {
					return nil, nil
				}
				if d2, ok := d.([]interface{}); ok {
					if j >= len(d2) {
						return nil, fmt.Errorf("navigate:: The index is out of array")
					}

					d = d2[j]
				} else {
					return nil, fmt.Errorf("navigateValue:: Can not convert (%s) to array, index=%s.%v", reflect.TypeOf(d), i, j)
				}
			}
		}
	}

	return d, nil
}

func isUserInput(d *schema.ResourceData, index []string, arrayIndex map[string]int) bool {
	var r = make([]string, 0, len(index)*2)
	for n, i := range index {
		r = append(r, i)

		if arrayIndex != nil {
			if j, ok := arrayIndex[strings.Join(index[:n+1], ".")]; ok {
				r = append(r, strconv.Itoa(j))
			}
		}
	}
	_, e := d.GetOkExists(strings.Join(r[:len(r)], "."))
	return e
}

func convertToInt(v interface{}) (interface{}, error) {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		return strconv.ParseInt(strVal, 10, 64)
	}
	return nil, fmt.Errorf("can not convert to integer")
}

func convertToStr(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil

	} else if i, ok := v.(int); ok {
		return strconv.Itoa(i), nil

	} else if b, ok := v.(bool); ok {
		return strconv.FormatBool(b), nil
	}

	return "", fmt.Errorf("can't convert to string")
}

func waitToFinish(target, pending []string, timeout, interval time.Duration, f resource.StateRefreshFunc) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Target:     target,
		Pending:    pending,
		Refresh:    f,
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: interval,
	}

	return stateConf.WaitForState()
}

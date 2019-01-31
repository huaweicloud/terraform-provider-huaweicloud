package huaweicloud

import (
	"encoding/json"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacks"
	"gopkg.in/yaml.v2"
)

// Takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}

func pointersMapToStringList(pointers map[string]*string) map[string]interface{} {
	list := make(map[string]interface{}, len(pointers))
	for i, v := range pointers {
		list[i] = *v
	}
	return list
}

// Takes a value containing JSON string and passes it through
// the JSON parser to normalize it, returns either a parsing
// error or normalized JSON string.
func normalizeJsonString(jsonString interface{}) (string, error) {
	var j interface{}

	if jsonString == nil || jsonString.(string) == "" {
		return "", nil
	}

	s := jsonString.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return s, err
	}

	// The error is intentionally ignored here to allow empty policies to passthrough validation.
	// This covers any interpolated values
	bytes, _ := json.Marshal(j)

	return string(bytes[:]), nil
}

// Takes a value containing YAML string and passes it through
// the YAML parser. Returns either a parsing
// error or original YAML string.
func checkYamlString(yamlString interface{}) (string, error) {
	var y interface{}

	if yamlString == nil || yamlString.(string) == "" {
		return "", nil
	}

	s := yamlString.(string)

	err := yaml.Unmarshal([]byte(s), &y)
	if err != nil {
		return s, err
	}

	return s, nil
}

func normalizeStackTemplate(templateString interface{}) (string, error) {
	if looksLikeJsonString(templateString) {
		return normalizeJsonString(templateString.(string))
	}

	return checkYamlString(templateString)
}

func flattenStackOutputs(stackOutputs []*stacks.Output) map[string]string {
	outputs := make(map[string]string, len(stackOutputs))
	for _, o := range stackOutputs {
		outputs[*o.OutputKey] = *o.OutputValue
	}
	return outputs
}

// flattenStackParameters is flattening list of
//  stack Parameters and only returning existing
// parameters to avoid clash with default values
func flattenStackParameters(stackParams map[string]string,
	originalParams map[string]interface{}) map[string]string {
	params := make(map[string]string, len(stackParams))
	for key, value := range stackParams {
		_, isConfigured := originalParams[key]
		if isConfigured {
			params[key] = value
		}
	}
	return params
}

// Takes the result of flatmap.Expand for an array of strings
// and returns a []*string
func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, v.(string))
		}
	}
	return vs
}

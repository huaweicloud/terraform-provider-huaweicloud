package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
)

func main() {
	provider := huaweicloud.Provider()
	CheckResourceMarkdown(provider)
	CheckDataSourcesMarkdown(provider)
}

func CheckResourceMarkdown(provider *schema.Provider) {
	checkMarkdown(provider.ResourcesMap, "../../docs/", "resources")
}

func CheckDataSourcesMarkdown(provider *schema.Provider) {
	checkMarkdown(provider.DataSourcesMap, "../../docs/", "data-sources")
}

func checkMarkdown(resources map[string]*schema.Resource, parentDir, rType string) {
	r := regexp.MustCompile("_v[1-9]$")
	for k, v := range resources {
		if r.MatchString(k) {
			continue
		}

		if v.DeprecationMessage != "" {
			continue
		}

		filePath, err := buildMarkdownFilePath(parentDir, rType, k)
		if err != nil {
			fmt.Printf("\n[WARN] can not generate the markdown file path: %s\n", err)
			continue
		}

		_, err = os.Stat(filePath)
		if isInternalResource(k) {
			if err == nil {
				fmt.Printf("\n[WARN] %s is only used for internal, please delete the file %s!\n", k, filePath)
			}
			continue
		}

		fmt.Printf("\n====== Checking for %s %s =====\n", rType, k)
		if err != nil {
			fmt.Printf("\n[WARN] can not state the markdown file: %s\n", err)
			continue
		}

		mdBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\n[WARN] can not read the markdown file: %s\n", err)
			continue
		}

		mdStr := string(mdBytes)
		checkMarkdownSchemas(v, "", &mdStr)
	}
}

func checkMarkdownSchemas(res *schema.Resource, parent string, mdContent *string) {
	if len(res.Schema) == 0 {
		return
	}

	for name, schemaObject := range res.Schema {
		fieldPath := buildFieldPath(parent, name)

		reg := fmt.Sprintf("[\\*|+] `%s` - ", name)
		isExist := checkArgumentExist(*mdContent, reg)

		extentStr, deprecated := buildAttributeString(schemaObject)
		// check deprecated field
		if deprecated || isDeprecatedField(name) {
			if isExist {
				fmt.Printf("`%s` was deprecated, should be deleted in Markdown\n", fieldPath)
			}
			continue
		}

		if !isExist {
			fmt.Printf("can not find `%s`, please check it\n", fieldPath)
		} else if extentStr != "" {
			// check the format of field
			reg += fmt.Sprintf("\\(%s\\)", extentStr)
			if !checkArgumentExist(*mdContent, reg) {
				fmt.Printf("the format of `%s` is not correct, should be (%s)\n", fieldPath, extentStr)
			}
		}

		// check nested block
		if schemaObject.Elem != nil {
			if nestedRes, ok := schemaObject.Elem.(*schema.Resource); ok {
				checkMarkdownSchemas(nestedRes, fieldPath, mdContent)
			}
		}
	}
}

func checkArgumentExist(mdStr, regStr string) bool {
	r := regexp.MustCompile(regStr)
	return r.MatchString(mdStr)
}

func buildAttributeString(s *schema.Schema) (string, bool) {
	reqd := s.Required
	opt := s.Optional
	computed := s.Computed
	forceNew := s.ForceNew

	// get extent attributes from description
	extent := parseExtentAttribute(s.Description)
	if s.Deprecated != "" || hasExtentAttribute(extent, "Deprecated") || hasExtentAttribute(extent, "Internal") {
		return "", true
	}

	if s.Required || hasExtentAttribute(extent, "Required") {
		reqd = true
		opt = false
		computed = false
	}

	// set the filed as Computed only
	if hasExtentAttribute(extent, "Computed") {
		computed = true
		reqd = false
		opt = false
		forceNew = false
	}

	// computed only
	if computed && !opt && !reqd {
		return "", false
	}

	attrs := make([]string, 0, 3)
	if reqd {
		attrs = append(attrs, "Required")
	} else {
		attrs = append(attrs, "Optional")
	}
	attrs = append(attrs, buildAttrType(s.Type))

	if forceNew {
		attrs = append(attrs, "ForceNew")
	}

	return strings.Join(attrs, ", "), false
}

func buildAttrType(ty schema.ValueType) string {
	switch ty {
	case schema.TypeBool:
		return "Bool"
	case schema.TypeInt:
		return "Int"
	case schema.TypeFloat:
		return "Float"
	case schema.TypeString:
		return "String"
	case schema.TypeList, schema.TypeSet:
		return "List"
	case schema.TypeMap:
		return "Map"
	}
	return ""
}

func buildFieldPath(parent, field string) string {
	if parent == "" {
		return field
	}
	return parent + "." + field
}

func buildMarkdownFilePath(parent, ty, name string) (string, error) {
	subParts := strings.SplitN(name, "_", 2)
	if len(subParts) < 2 {
		return "", fmt.Errorf("the format of %s %s is not correct", ty, name)
	}

	return path.Join(parent, ty, fmt.Sprintf("%s.md", subParts[1])), nil
}

func isInternalResource(key string) bool {
	internalResources := []string{
		"apm_aksk", "aom_alarm_policy", "aom_prometheus_instance",
		"aom_application", "aom_component", "aom_environment", "aom_cmdb_resource_relationships",
		"lts_access_rule", "lts_dashboard", "lts_struct_template", "elb_log",
		// the fellowings are legacy
		"ges_graph", "iam_agency", "networking_eip_associate",
		"vpc_ids", "vpc_route_ids",
	}
	for _, v := range internalResources {
		if strings.HasSuffix(key, v) {
			return true
		}
	}
	return false
}

func isDeprecatedField(field string) bool {
	// deprecatedFields includes the fields that shoud always be ignored.
	var deprecatedFields = []string{"tenant_id", "admin_state_up", "auto_pay"}
	for _, key := range deprecatedFields {
		if field == key {
			return true
		}
	}
	return false
}

func parseExtentAttribute(desc string) map[string]bool {
	extra := make(map[string]bool)
	prefix := "schema:"

	if !strings.HasPrefix(desc, prefix) {
		return extra
	}

	validDesc := strings.SplitN(desc, ";", 2)[0]
	allAttr := strings.Split(validDesc[len(prefix):], ",")
	for _, ext := range allAttr {
		if attr := strings.TrimLeft(ext, " "); attr != "" {
			extra[attr] = true
		}
	}

	return extra
}

func hasExtentAttribute(extra map[string]bool, key string) bool {
	return extra[key]
}

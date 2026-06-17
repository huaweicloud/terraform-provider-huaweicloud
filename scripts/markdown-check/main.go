package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ignoreDeprecatedNames = []string{
	"huaweicloud_identitycenter_bearer_token",
}

func main() {
	provider := huaweicloud.Provider()

	var errCode int
	errCode += CheckResourceMarkdown(provider)
	errCode += CheckDataSourcesMarkdown(provider)
	os.Exit(errCode)
}

func CheckResourceMarkdown(provider *schema.Provider) int {
	errCount := checkMarkdown(provider.ResourcesMap, "../../docs/", "resources")
	fmt.Printf("\n====== Summary: Find <%d> inconsistencies between schemas and docs in resources ======\n", errCount)
	return errCount
}

func CheckDataSourcesMarkdown(provider *schema.Provider) int {
	errCount := checkMarkdown(provider.DataSourcesMap, "../../docs/", "data-sources")
	fmt.Printf("\n====== Summary: Find <%d> inconsistencies between schemas and docs in data sources ======\n", errCount)
	return errCount
}

func checkMarkdown(resources map[string]*schema.Resource, parentDir, rType string) int {
	var totalCount int

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
		if isInternalResource(v, k) {
			if err == nil {
				fmt.Printf("\n[WARN] %s is only used for internal, please check the file %s!\n", k, filePath)
			}
			continue
		}

		fmt.Printf("\n====== Checking for %s %s ======\n", rType, k)
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
		totalCount += checkMarkdownSchemas(v, v, "", &mdStr, k)
	}

	return totalCount
}

func checkMarkdownSchemas(res, rootRes *schema.Resource, parent string, mdContent *string, rName string) int {
	if len(res.Schema) == 0 {
		return 0
	}

	var count int
	for name, schemaObject := range res.Schema {
		fieldPath := buildFieldPath(parent, name)

		reg := fmt.Sprintf("[\\*|+] `%s` - ", name)
		isExist := checkArgumentExist(*mdContent, reg)

		extentStr, deprecated := buildAttributeString(schemaObject)
		deprecated = deprecated || isDeprecatedField(name)
		if deprecated && utils.StrSliceContains(ignoreDeprecatedNames, rName) {
			continue
		}

		// check deprecated field
		if deprecated && !hasIdenticalSchemaField(rootRes, name) {
			if isExist {
				count++
				fmt.Printf("[ERROR] `%s` was deprecated, should be deleted in Markdown\n", fieldPath)
			}
			continue
		}

		// (computed only) nested attributes can have another format:
		// * `<parent>/<attr name>` - The description ...
		if !isExist && parent != "" && extentStr == "" {
			reg1 := fmt.Sprintf("[\\*|+] `%s` - ", strings.ReplaceAll(fieldPath, ".", "/"))
			isExist = checkArgumentExist(*mdContent, reg1)
		}

		if !isExist {
			count++
			fmt.Printf("[ERROR] can not find `%s`, please check it\n", fieldPath)
		} else if extentStr != "" {
			// check the format of field
			reg += fmt.Sprintf("\\(%s", extentStr)
			if !checkArgumentExist(*mdContent, reg) {
				count++
				fmt.Printf("[ERROR] the format of `%s` is not correct, should be (%s)\n", fieldPath, extentStr)
			}
		}

		// check nested block
		if schemaObject.Elem != nil {
			if nestedRes, ok := schemaObject.Elem.(*schema.Resource); ok {
				count += checkMarkdownSchemas(nestedRes, rootRes, fieldPath, mdContent, rName)
			}
		}
	}
	return count
}

// check the identical field which is not deprecated
func hasIdenticalSchemaField(res *schema.Resource, key string) bool {
	if len(res.Schema) == 0 {
		return false
	}

	for name, schemaObject := range res.Schema {
		if _, deprecated := buildAttributeString(schemaObject); deprecated || isDeprecatedField(name) {
			continue
		}

		if name == key {
			return true
		}

		// find key in nested block
		if schemaObject.Elem != nil {
			if nestedRes, ok := schemaObject.Elem.(*schema.Resource); ok {
				if ret := hasIdenticalSchemaField(nestedRes, key); ret {
					return true
				}
			}
		}
	}

	return false
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

func isInternalResource(resource *schema.Resource, key string) bool {
	internalResources := []string{
		"apm_aksk", "aom_alarm_policy", "aom_prometheus_instance",
		"aom_application", "aom_component", "aom_environment", "aom_cmdb_resource_relationships",
		"lts_access_rule", "lts_dashboard", "lts_struct_template", "elb_log",
		// the followings have changed to rds_mysql_xxx after v1.47.0
		"rds_account", "rds_database", "rds_database_privilege",
		"rf_stack", // changed to rfs_stack after v1.47.0
		// the fellowings are legacy
		"iam_agency", "networking_eip_associate", "vpc_ids", "identity_role_assignment",
	}
	for _, v := range internalResources {
		if strings.HasSuffix(key, v) {
			return true
		}
	}

	if resource.Description != "" {
		// get extent attributes from description
		extent := parseExtentAttribute(resource.Description)
		if hasExtentAttribute(extent, "Internal") {
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

func parseExtentAttribute(desc string) map[string]interface{} {
	extra := make(map[string]interface{})
	prefix := "schema:"

	if strings.HasPrefix(desc, prefix) {
		validDesc := strings.SplitN(desc, ";", 2)[0]
		allAttrJson := validDesc[len(prefix):]
		err := json.Unmarshal([]byte(strings.Trim(allAttrJson, " ")), &extra)
		if err == nil {
			return extra
		}

		validDesc = strings.SplitN(desc, ";", 2)[0]
		allAttr := strings.Split(validDesc[len(prefix):], ",")
		for _, ext := range allAttr {
			if attr := strings.TrimLeft(ext, " "); attr != "" {
				extra[attr] = true
			}
		}
		return extra
	}

	return nil
}

func hasExtentAttribute(extra map[string]interface{}, key string) bool {
	if v, ok := extra[key]; ok {
		return v.(bool)
	}

	return false
}

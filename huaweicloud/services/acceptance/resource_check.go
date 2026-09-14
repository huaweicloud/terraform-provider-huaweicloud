package acceptance

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ServiceFunc the resource query functions
type ServiceFunc func(*config.Config, *terraform.ResourceState) (interface{}, error)

// ResourceCheck resource check object
type ResourceCheck struct {
	resourceName    string
	resourceObject  interface{}
	getResourceFunc ServiceFunc
	resourceType    string
}

const (
	resourceTypeCode   = "resource"
	dataSourceTypeCode = "dataSource"

	checkAttrRegexpStr = `^\$\{([^\}]+)\}$`
)

var checkAttrRegexp = regexp.MustCompile(checkAttrRegexpStr)

/*
InitDataSourceCheck build a 'ResourceCheck' object. Only used to check datasource attributes.

	Parameters:
	  dName: The data source name is used to check in the terraform.State. e.g. data.huaweicloud_css_flavors.test
	Return:
	  *ResourceCheck: ResourceCheck object
*/
func InitDataSourceCheck(dName string) *ResourceCheck {
	return &ResourceCheck{
		resourceName: dName,
		resourceType: dataSourceTypeCode,
	}
}

/*
InitResourceCheck build a 'ResourceCheck' object. The common test methods are provided in 'ResourceCheck'.

	Parameters:
	  rName:           The resource name is used to check in the terraform.State. e.g. huaweicloud_waf_domain.domain_1
	  rObject:         Resource object pointer, used to check whether the resource exists
	  getResourceFunc: The function used to get the resource object.
	Return:
	  *ResourceCheck: ResourceCheck object
*/
func InitResourceCheck(rName string, rObject interface{}, getResourceFunc ServiceFunc) *ResourceCheck {
	return &ResourceCheck{
		resourceName:    rName,
		resourceObject:  rObject,
		getResourceFunc: getResourceFunc,
		resourceType:    resourceTypeCode,
	}
}

func parseVariableToName(variable string) (string, string, error) {
	var name, field string

	// Check the format of the variable
	mArr := checkAttrRegexp.FindStringSubmatch(variable)
	if len(mArr) != 2 {
		return name, field, fmt.Errorf("the type of 'variable' is error, "+
			"expected ${resource-type.name.field} but got %s", variable)
	}

	// Get name and field from variable
	strs := strings.Split(mArr[1], ".")
	keyIndex := 2
	if strs[0] == "data" {
		keyIndex = 3
	}

	if len(strs) <= keyIndex {
		return name, field, fmt.Errorf("attribute field is missing: "+
			"expected ${resource-type.name.field} but got %s", variable)
	}

	name = strings.Join(strs[0:keyIndex], ".")
	field = strings.Join(strs[keyIndex:], ".")

	return name, field, nil
}

/*
TestCheckResourceAttrWithVariable validates the pair variable in state for the given name/key combination.

	Parameters:
	  name: The resource or data source name is used to check in the terraform.State.
	  key:  The field name of the resource.
	  pair: The pair name of the value to be checked.

	  pair such like ${huaweicloud_waf_certificate.certificate_1.id}
	  or ${data.huaweicloud_waf_policies.policies_2.policies.0.id}
*/
func TestCheckResourceAttrWithVariable(name, key, pair string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pairName, pairKey, err := parseVariableToName(pair)
		if err != nil {
			return err
		}

		if strings.EqualFold(name, pairName) {
			return fmt.Errorf("meaningless verification: " +
				"The referenced resource cannot be the current resource")
		}

		// Get the value based on pairName and pairKey from the state.
		rs, ok := s.RootModule().Resources[pairName]
		if !ok {
			return fmt.Errorf("can't find %s in state: %v", pairName, ok)
		}
		value := rs.Primary.Attributes[pairKey]

		return resource.TestCheckResourceAttr(name, key, value)(s)
	}
}

// CheckResourceDestroy check whether resources destroyed
func (rc *ResourceCheck) CheckResourceDestroy() resource.TestCheckFunc {
	if strings.Compare(rc.resourceType, dataSourceTypeCode) == 0 {
		return nil
	}

	return func(s *terraform.State) error {
		strs := strings.Split(rc.resourceName, ".")
		resourceType := strs[0]

		if resourceType == "" || resourceType == "data" {
			return fmt.Errorf("the format of the resource name is invalid, please check your configuration")
		}

		if rc.getResourceFunc == nil {
			return fmt.Errorf("the 'getResourceFunc' is nil, please set it during initialization")
		}

		conf := TestAccProvider.Meta().(*config.Config)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			if _, err := rc.getResourceFunc(conf, rs); err == nil {
				return fmt.Errorf("failed to destroy the %s resource: %s still exists",
					resourceType, rs.Primary.ID)
			}
		}
		return nil
	}
}

func (rc *ResourceCheck) checkResourceExists(s *terraform.State) error {
	rs, ok := s.RootModule().Resources[rc.resourceName]
	if !ok {
		return fmt.Errorf("can not found the resource or data source in state: %s", rc.resourceName)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No id set for the resource or data source: %s", rc.resourceName)
	}
	if strings.EqualFold(rc.resourceType, dataSourceTypeCode) {
		return nil
	}

	if rc.getResourceFunc == nil {
		return fmt.Errorf("the 'getResourceFunc' is nil, please set it during initialization")
	}

	conf := TestAccProvider.Meta().(*config.Config)
	r, err := rc.getResourceFunc(conf, rs)
	if err != nil {
		return fmt.Errorf("checking resource %s %s exists error: %s ",
			rc.resourceName, rs.Primary.ID, err)
	}

	b, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshaling resource %s %s error: %s ",
			rc.resourceName, rs.Primary.ID, err)
	}

	// unmarshal the response body into the resourceObject
	if rc.resourceObject != nil {
		return json.Unmarshal(b, rc.resourceObject)
	}

	return nil
}

// CheckResourceExists check whether resources exist
func (rc *ResourceCheck) CheckResourceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return rc.checkResourceExists(s)
	}
}

/*
CheckMultiResourcesExists checks whether multiple resources created by count are both existed.

	Parameters:
	  count: the expected number of resources that will be created.
*/
func (rc *ResourceCheck) CheckMultiResourcesExists(count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var err error
		for i := 0; i < count; i++ {
			rcCopy := *rc
			rcCopy.resourceName = fmt.Sprintf("%s.%d", rcCopy.resourceName, i)
			err = rcCopy.checkResourceExists(s)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

/*
MultiResourcesImportSteps builds ImportState TestSteps for multiple resources created by count.

Unlike CheckMultiResourcesExists, import verification is driven by TestStep rather than TestCheckFunc,
so this helper returns a list of TestSteps that can be appended to the Steps of a TestCase.

	Parameters:
	  count: the expected number of resources that will be created.
	  opts:  optional import behaviors, such as ignored fields or a custom import ID builder.

	Example (default, using resource ID):
	  rc.MultiResourcesImportSteps(2)

	Example (ignore fields):
	  rc.MultiResourcesImportSteps(2, acceptance.WithImportStateVerifyIgnore("password"))

	Example (composite import ID, e.g. <instance_id>/<id>):
	  rc.MultiResourcesImportSteps(2, acceptance.WithImportIDAttributes("instance_id", "id"))

	Example (fully custom composition):
	  rc.MultiResourcesImportSteps(2, acceptance.WithImportStateIDFunc(func(rs *terraform.ResourceState) (string, error) {
	      return fmt.Sprintf("%s/%s", rs.Primary.Attributes["vault_id"], rs.Primary.Attributes["policy_id"]), nil
	  }))
*/
func (rc *ResourceCheck) MultiResourcesImportSteps(count int, opts ...MultiResourcesImportOption) []resource.TestStep {
	steps := make([]resource.TestStep, 0, count)
	for i := 0; i < count; i++ {
		// terraform.State stores counted resources as "type.name.index".
		stateName := fmt.Sprintf("%s.%d", rc.resourceName, i)
		// terraform import requires HCL address syntax: "type.name[index]".
		importAddr := fmt.Sprintf("%s[%d]", rc.resourceName, i)

		step := resource.TestStep{
			ResourceName:      importAddr,
			ImportState:       true,
			ImportStateVerify: true,
			// Always set ImportStateIdFunc so the SDK does not look up ResourceName
			// (HCL address) in terraform.State, which uses the dotted index key.
			ImportStateIdFunc: defaultImportStateIdFunc(stateName),
		}
		for _, opt := range opts {
			opt(stateName, &step)
		}
		steps = append(steps, step)
	}
	return steps
}

func defaultImportStateIdFunc(stateName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[stateName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", stateName)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("no ID is set for the resource: %s", stateName)
		}
		return rs.Primary.ID, nil
	}
}

// MultiResourcesImportOption customizes a generated import TestStep.
// The resourceName argument is the terraform.State key (type.name.index).
type MultiResourcesImportOption func(resourceName string, step *resource.TestStep)

// WithImportStateVerifyIgnore ignores the specified attributes during ImportStateVerify.
func WithImportStateVerifyIgnore(fields ...string) MultiResourcesImportOption {
	return func(_ string, step *resource.TestStep) {
		step.ImportStateVerifyIgnore = fields
	}
}

// WithImportStateIDFunc sets a custom import ID builder based on the current resource state.
// Use this when the import ID is a composite value that cannot be derived from Primary.ID alone.
func WithImportStateIDFunc(build func(rs *terraform.ResourceState) (string, error)) MultiResourcesImportOption {
	return func(resourceName string, step *resource.TestStep) {
		step.ImportStateIdFunc = func(s *terraform.State) (string, error) {
			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return "", fmt.Errorf("resource (%s) not found", resourceName)
			}
			if rs.Primary.ID == "" {
				return "", fmt.Errorf("no ID is set for the resource: %s", resourceName)
			}
			return build(rs)
		}
	}
}

// WithImportIDAttributes builds a composite import ID by joining the specified attribute keys with "/".
// The special key "id" uses rs.Primary.ID. This covers the common "<attr1>/<attr2>" import format.
func WithImportIDAttributes(keys ...string) MultiResourcesImportOption {
	return WithImportStateIDFunc(func(rs *terraform.ResourceState) (string, error) {
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			value := rs.Primary.ID
			if key != "id" {
				value = rs.Primary.Attributes[key]
			}
			if value == "" {
				return "", fmt.Errorf("invalid format specified for import ID, attribute (%s) is empty", key)
			}
			parts = append(parts, value)
		}
		return strings.Join(parts, "/"), nil
	})
}

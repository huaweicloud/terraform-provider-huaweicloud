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
	reg, err := regexp.Compile(checkAttrRegexpStr)
	if err != nil {
		return name, field, err
	}
	mArr := reg.FindStringSubmatch(variable)
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

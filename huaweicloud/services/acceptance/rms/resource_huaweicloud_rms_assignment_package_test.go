package rms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAssignmentPackageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getAssignmentPackage: Query the RMS assignment package
	var (
		getAssignmentPackageHttpUrl = "v1/resource-manager/domains/{domain_id}/conformance-packs/{id}"
		getAssignmentPackageProduct = "rms"
	)
	getAssignmentPackageClient, err := cfg.NewServiceClient(getAssignmentPackageProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RMS client: %s", err)
	}

	getAssignmentPackagePath := getAssignmentPackageClient.Endpoint + getAssignmentPackageHttpUrl
	getAssignmentPackagePath = strings.ReplaceAll(getAssignmentPackagePath, "{domain_id}", cfg.DomainID)
	getAssignmentPackagePath = strings.ReplaceAll(getAssignmentPackagePath, "{id}", state.Primary.ID)

	getAssignmentPackageOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAssignmentPackageResp, err := getAssignmentPackageClient.Request("GET", getAssignmentPackagePath,
		&getAssignmentPackageOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS assignment package: %s", err)
	}

	getAssignmentPackageRespBody, err := utils.FlattenResponse(getAssignmentPackageResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RMS assignment package: %s", err)
	}

	return getAssignmentPackageRespBody, nil
}

func TestAccAssignmentPackage_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rms_assignment_package.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAssignmentPackageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAssignmentPackage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "vars_structure",
						"data.huaweicloud_rms_assignment_package_templates.test", "templates.0.parameters"),
					resource.TestCheckResourceAttrSet(rName, "stack_id"),
					resource.TestCheckResourceAttrSet(rName, "stack_name"),
					resource.TestCheckResourceAttrSet(rName, "deployment_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testAssignmentPackage_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "vars_structure.*", map[string]string{
						"var_key":   "lastBackupAgeValue",
						"var_value": "25",
					}),
					resource.TestCheckResourceAttrSet(rName, "stack_id"),
					resource.TestCheckResourceAttrSet(rName, "stack_name"),
					resource.TestCheckResourceAttrSet(rName, "deployment_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"agency_name", "template_key", "template_body", "template_uri"},
			},
		},
	})
}

func testAssignmentPackage_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_rms_assignment_package_templates" "test" {
  template_key = "Operational-Best-Practices-for-ECS.tf.json"
}

resource "huaweicloud_rms_assignment_package" "test" {
  name         = "%s"
  template_key = data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key

  dynamic "vars_structure" {
    for_each = data.huaweicloud_rms_assignment_package_templates.test.templates.0.parameters
    content {
      var_key   = vars_structure.value["name"]
      var_value = vars_structure.value["default_value"]
    }
  }
}
`, name)
}

func testAssignmentPackage_update(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_rms_assignment_package_templates" "test" {
  template_key = "Operational-Best-Practices-for-ECS.tf.json"
}

resource "huaweicloud_rms_assignment_package" "test" {
  name         = "%s-update"
  template_key = data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key

  dynamic "vars_structure" {
    for_each = data.huaweicloud_rms_assignment_package_templates.test.templates.0.parameters
    content {
      var_key   = vars_structure.value["name"]
      var_value = vars_structure.value["name"] == "lastBackupAgeValue" ? 25 : vars_structure.value["default_value"]
    }
  }
}
`, name)
}

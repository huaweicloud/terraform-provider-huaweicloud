package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/valuelists"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getReferenceTableResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	wafClient, err := cfg.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}
	return valuelists.GetWithEpsID(wafClient, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccReferenceTable_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_waf_reference_table.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = name + "_update"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getReferenceTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafReferenceTable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "tf acc"),
					resource.TestCheckResourceAttr(resourceName, "type", "url"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
				),
			},
			{
				Config: testAccWafReferenceTable_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "url"),
					resource.TestCheckResourceAttr(resourceName, "conditions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccWafReferenceTable_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_reference_table" "test" {
  name                  = "%[1]s"
  type                  = "url"
  description           = "tf acc"
  enterprise_project_id = "%[2]s"

  conditions = [
    "/admin",
    "/manage"
  ]
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafReferenceTable_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_reference_table" "test" {
  name                  = "%[1]s"
  type                  = "url"
  description           = ""
  enterprise_project_id = "%[2]s"

  conditions = [
    "/bill",
    "/sql"
  ]
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

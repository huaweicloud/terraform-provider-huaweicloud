package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getResourceScanTemplateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetScanTemplateInfo(client, "id", state.Primary.ID)
}

func TestAccResourceScanTemplate_basic(t *testing.T) {
	var (
		rName    = "huaweicloud_dsc_scan_template.test"
		scanName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceScanTemplateFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScanTemplate_basic(scanName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", scanName),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
				),
			},
			{
				Config: testAccResourceScanTemplate_update(scanName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", scanName)),
					resource.TestCheckResourceAttr(rName, "description", "test_description_update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"action", "add_built_in_rules"},
			},
		},
	})
}

func TestAccResourceScanTemplate_copy(t *testing.T) {
	var (
		rName    = "huaweicloud_dsc_scan_template.copy"
		scanName = acceptance.RandomAccResourceName()

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceScanTemplateFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScanTemplate_copy(scanName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_copy", scanName)),
					resource.TestCheckResourceAttr(rName, "description", "copy_template"),
				),
			},
			{
				Config: testAccResourceScanTemplate_copyUpdate(scanName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_copy_update", scanName)),
					resource.TestCheckResourceAttr(rName, "description", "copy_template_updated"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action", "add_built_in_rules", "origin_template_id"},
			},
		},
	})
}

func testAccResourceScanTemplate_basic(scanName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_template" "test" {
  action             = "ADD"
  name               = "%s"
  description        = "test_description"
  add_built_in_rules = true
}
`, scanName)
}

func testAccResourceScanTemplate_update(scanName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_template" "test" {
  action             = "ADD"
  name               = "%s_update"
  description        = "test_description_update"
  add_built_in_rules = true
}
`, scanName)
}

func testAccResourceScanTemplate_copy(scanName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_template" "origin" {
  action             = "ADD"
  name               = "%s"
  description        = "origin_template"
  add_built_in_rules = true
}

resource "huaweicloud_dsc_scan_template" "copy" {
  action             = "COPY"
  name               = "%s_copy"
  description        = "copy_template"
  origin_template_id = huaweicloud_dsc_scan_template.origin.id
}
`, scanName, scanName)
}

func testAccResourceScanTemplate_copyUpdate(scanName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_template" "origin" {
  action             = "ADD"
  name               = "%s"
  description        = "origin_template"
  add_built_in_rules = true
}

resource "huaweicloud_dsc_scan_template" "copy" {
  action             = "COPY"
  name               = "%s_copy_update"
  description        = "copy_template_updated"
  origin_template_id = huaweicloud_dsc_scan_template.origin.id
}
`, scanName, scanName)
}

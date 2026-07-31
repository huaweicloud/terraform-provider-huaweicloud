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

func getScanSecurityLevelResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetScanSecurityLevelById(client, state.Primary.ID)
}

// Before this test, please ensure that the DSC instance has been created.
func TestAccResourceScanSecurityLevel_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dsc_scan_security_level.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getScanSecurityLevelResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceScanSecurityLevel_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "security_level_name", name),
					resource.TestCheckResourceAttr(rName, "color_number", "6"),
					resource.TestCheckResourceAttr(rName, "security_level_desc", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(rName, "category"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				Config: testAccResourceScanSecurityLevel_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "security_level_name", updateName),
					resource.TestCheckResourceAttr(rName, "color_number", "7"),
					resource.TestCheckResourceAttr(rName, "security_level_desc", "Updated by terraform script"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceScanSecurityLevel_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_security_level" "test" {
  security_level_name = "%[1]s"
  color_number        = 6
  security_level_desc = "Created by terraform script"
}
`, name)
}

func testAccResourceScanSecurityLevel_basic_step2(updateName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_scan_security_level" "test" {
  security_level_name = "%[1]s"
  color_number        = 7
  security_level_desc = "Updated by terraform script"
}
`, updateName)
}

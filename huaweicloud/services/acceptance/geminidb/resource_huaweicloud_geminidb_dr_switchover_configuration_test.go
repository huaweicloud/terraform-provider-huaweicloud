package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBDRSwitchoverConfiguration_basic(t *testing.T) {
	rName := "huaweicloud_geminidb_dr_switchover_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBDRSwitchoverConfiguration_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "switchover_ratio", "50"),
					resource.TestCheckResourceAttr(rName, "sync_delay", "10"),
				),
			},
			{
				Config: testAccGeminiDBDRSwitchoverConfiguration_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "switchover_ratio", "60"),
					resource.TestCheckResourceAttr(rName, "sync_delay", "300"),
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

func testAccGeminiDBDRSwitchoverConfiguration_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_dr_switchover_configuration" "test" {
  instance_id      = "%s"
  switchover_ratio = 50
  sync_delay       = 10
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}

func testAccGeminiDBDRSwitchoverConfiguration_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_dr_switchover_configuration" "test" {
  instance_id      = "%s"
  switchover_ratio = 60
  sync_delay       = 300
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}

package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccountEnabledControl_basic(t *testing.T) {
	rName := "data.huaweicloud_rgc_account_enabled_control.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAccountEnabledControl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccountEnabledControl_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "control_detail.#"),
					resource.TestCheckResourceAttrSet(rName, "control_detail.0.control_identifier"),
					resource.TestCheckResourceAttrSet(rName, "regions.#"),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestCheckResourceAttrSet(rName, "message"),
					resource.TestCheckResourceAttrSet(rName, "version"),
				),
			},
		},
	})
}

func testAccDataSourceAccountEnabledControl_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_account_enabled_control" "test" {
  control_id         = "%[1]s"
  managed_account_id = "%[2]s"
}
`, acceptance.HW_RGC_CONTROL_ID, acceptance.HW_RGC_ACCOUNT_ID)
}

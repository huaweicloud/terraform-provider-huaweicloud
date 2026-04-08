package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAttackLogTotal_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_attack_log_total.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a firewall instance ID.
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAttackLogTotal_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceAttackLogTotal_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_attack_log_total" "test" {
  fw_instance_id = "%s"
  log_type       = "internet"
  range          = "0"
  start_time     = 1774730600349
  end_time       = 1775076220000
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

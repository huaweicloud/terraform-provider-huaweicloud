package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAttackLogStatistic_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_attack_log_statistic.test"
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
				Config: testDataSourceAttackLogStatistic_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceAttackLogStatistic_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_attack_log_statistic" "test" {
  fw_instance_id = "%s"
  log_type       = "internet"
  item           = "src"
  range          = "1"
  direction      = "in2out"
  start_time     = 1774730600349
  end_time       = 1775076220000
  size           = "50"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessLogStatistic_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_access_log_statistic.test"
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
				Config: testDataSourceAccessLogStatistic_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceAccessLogStatistic_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_access_log_statistic" "test" {
  fw_instance_id = "%s"
  item           = "strategy_dashboard"
  range          = "0"
  direction      = "in2out"
  log_type       = "internet"
  start_time     = 1774730600349
  end_time       = 1775076220000
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

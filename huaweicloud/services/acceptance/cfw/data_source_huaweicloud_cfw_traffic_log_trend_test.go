package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTrafficLogTrend_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_traffic_log_trend.test"
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
				Config: testDataSourceTrafficLogTrend_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.agg_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.in_bps"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.out_bps"),
				),
			},
		},
	})
}

func testDataSourceTrafficLogTrend_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_traffic_log_trend" "test" {
  fw_instance_id = "%s"
  log_type       = "vpc"
  agg_type       = "max"
  range          = "2"
  start_time     = 1774730600349
  end_time       = 1775076220000
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

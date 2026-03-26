package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFlowLogStatisticDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_flow_log_statistic_detail.test"
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
				Config: testDataSourceFlowLogStatisticDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.app_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.bytes"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.dst_ip_count"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.records.#"),
				),
			},
		},
	})
}

func testDataSourceFlowLogStatisticDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_flow_log_statistic_detail" "test" {
  fw_instance_id = "%s"
  log_type       = "internet"
  item           = "dst_ip"
  value          = "100.93.4.158"
  range          = "2"
  direction      = "in2out"
  start_time     = 1774730600349
  end_time       = 1775076220000
  asset_type     = "private"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

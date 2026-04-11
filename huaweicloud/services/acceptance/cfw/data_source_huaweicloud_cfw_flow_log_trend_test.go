package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFlowLogTrend_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_flow_log_trend.test"
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
				Config: testDataSourceFlowLogTrend_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.agg_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.in_bps"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.out_bps"),
				),
			},
		},
	})
}

func testDataSourceFlowLogTrend_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_flow_log_trend" "test" {
  fw_instance_id = "%s"
  log_type       = "vpc"
  range          = "0"
  direction      = "in2out"
  start_time     = 1774730600349
  end_time       = 1775076220000
  vpc            = ["test1", "test2"]
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

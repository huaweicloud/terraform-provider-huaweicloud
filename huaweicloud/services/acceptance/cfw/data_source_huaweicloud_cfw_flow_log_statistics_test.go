package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFlowLogStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_flow_log_statistics.test"
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
				Config: testDataSourceFlowLogStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
				),
			},
		},
	})
}

func testDataSourceFlowLogStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_flow_log_statistics" "test" {
  fw_instance_id = "%s"
  log_type       = "internet"
  item           = "dst_host"
  range          = "2"
  direction      = "in2out"
  start_time     = 1774730600349
  end_time       = 1775076220000
  asset_type     = "private"
  size           = "5"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

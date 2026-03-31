package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessLogStatisticDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_access_log_statistic_detail.test"
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
				Config: testDataSourceAccessLogStatisticDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceAccessLogStatisticDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_access_log_statistic_detail" "test" {
  fw_instance_id = "%s"
  item           = "dst_ip"
  item_id        = "100.93.4.158"
  range          = "1"
  direction      = "in2out"
  log_type       = "internet"
  start_time     = 1774730600349
  end_time       = 1775076220000
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_logs.test"
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
				Config: testDataSourceLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
				),
			},
		},
	})
}

func testDataSourceLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_logs" "test" {
  fw_instance_id = "%s"
  start_time     = 1751952647737
  end_time       = 1751963447737
  log_type       = "internet"
  type           = "flow"

  filters {
    field    = "src_ip"
    operator = "equal"
    values   = ["1", "2"]
  }
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

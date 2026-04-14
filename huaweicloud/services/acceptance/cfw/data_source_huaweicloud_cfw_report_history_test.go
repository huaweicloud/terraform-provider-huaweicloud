package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReportHistory_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cfw_report_history.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a firewall instance ID and a report profile ID.
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwReportProfile(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceReportHistory_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.#"),
				),
			},
		},
	})
}

func testDataSourceReportHistory_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_report_history" "test" {
  fw_instance_id    = "%s"
  report_profile_id = "%s"
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_REPORT_PROFILE_ID)
}

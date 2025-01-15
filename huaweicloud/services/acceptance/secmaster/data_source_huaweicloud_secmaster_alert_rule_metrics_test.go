package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterAlertRuleMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_alert_rule_metrics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterAlertRuleMetrics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "metrics"),
				),
			},
		},
	})
}

func testDataSourceSecmasterAlertRuleMetrics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_alert_rule_metrics" "test" {
  workspace_id = "%s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}

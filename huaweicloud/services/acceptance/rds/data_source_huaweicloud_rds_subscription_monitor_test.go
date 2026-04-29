package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSubscriptionMonitor_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_subscription_monitor.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRdsSubscriptionMonitor_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "latency"),
					resource.TestCheckResourceAttrSet(dataSource, "agent_not_running"),
					resource.TestCheckResourceAttrSet(dataSource, "pending_cmd_count"),
					resource.TestCheckResourceAttrSet(dataSource, "last_dist_sync"),
					resource.TestCheckResourceAttrSet(dataSource, "estimated_process_time"),
				),
			},
		},
	})
}

func testAccDataSourceRdsSubscriptionMonitor_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_subscriptions" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_subscription_monitor" "test" {
  instance_id     = "%[1]s"
  subscription_id = data.huaweicloud_rds_subscriptions.test.subscriptions[0].id
}
`, acceptance.HW_RDS_INSTANCE_ID)
}

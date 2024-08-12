package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmSubscriptions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_alarm_subscriptions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmSubscriptions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "subscriptions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.notification_target_type"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.notification_target"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.notification_target_name"),
					resource.TestCheckResourceAttrSet(dataSource, "subscriptions.0.time_zone"),
					resource.TestCheckOutput("is_exist_current_subscription", "true"),
				),
			},
		},
	})
}

func testDataSourceAlarmSubscriptions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dws_alarm_subscriptions" "test" {
  depends_on = [huaweicloud_dws_alarm_subscription.test]
}

output "is_exist_current_subscription"{
  # The contains method is precision matching.
  value = contains(data.huaweicloud_dws_alarm_subscriptions.test.subscriptions[*].name, huaweicloud_dws_alarm_subscription.test.name)
}
`, testDwsAlarmSubs_basic(name))
}

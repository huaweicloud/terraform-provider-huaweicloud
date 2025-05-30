package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAlarmNotifications_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_alarm_notifications.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, please prepare a WAF alarm notification.
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAlarmNotifications_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.notice_class"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.times"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.sendfreq"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.threat.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.enterprise_project_id"),
				),
			},
		},
	})
}

const testDataSourceAlarmNotifications_basic = `data "huaweicloud_waf_alarm_notifications" "test" {}`

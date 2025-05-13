package cnad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Field `topic_urn` has a value only after configuration. This test case does not test this field yet.
func TestAccDatasourceAlarmNotifications_basic(t *testing.T) {
	rName := "data.huaweicloud_cnad_advanced_alarm_notifications.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAlarmNotifications_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "is_close_attack_source_flag"),
				),
			},
		},
	})
}

const testAccDatasourceAlarmNotifications_basic = `
data "huaweicloud_cnad_advanced_alarm_notifications" "test" {}
`

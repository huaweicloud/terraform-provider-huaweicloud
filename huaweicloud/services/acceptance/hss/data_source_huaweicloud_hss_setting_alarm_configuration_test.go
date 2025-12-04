package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingAlarmConfiguration_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_alarm_configuration.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSettingAlarmConfiguration_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_type"),
					resource.TestCheckResourceAttrSet(dataSource, "display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "daily_alarm"),
					resource.TestCheckResourceAttrSet(dataSource, "realtime_alarm"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_level.#"),
				),
			},
		},
	})
}

const testAccDataSourceSettingAlarmConfiguration_basic = `
data "huaweicloud_hss_setting_alarm_configuration" "test" {
  enterprise_project_id = "0"
}
`

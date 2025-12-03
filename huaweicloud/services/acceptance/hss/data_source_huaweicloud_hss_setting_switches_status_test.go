package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingSwitchesStatus_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_switches_status.test"
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
				Config: testAccDataSourceSettingSwitchesStatus_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enabled"),
				),
			},
		},
	})
}

const testAccDataSourceSettingSwitchesStatus_basic = `
data "huaweicloud_hss_setting_switches_status" "test" {
  code                  = "image_pay_per_scan"
  enterprise_project_id = "0"
}
`

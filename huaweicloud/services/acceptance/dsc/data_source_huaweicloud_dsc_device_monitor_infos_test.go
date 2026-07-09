package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDeviceMonitorInfos_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_device_monitor_infos.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDeviceMonitorInfos_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "monitor_infos.#"),
				),
			},
		},
	})
}

const testAccDataSourceDscDeviceMonitorInfos_basic = `
data "huaweicloud_dsc_device_monitor_infos" "test" {}
`

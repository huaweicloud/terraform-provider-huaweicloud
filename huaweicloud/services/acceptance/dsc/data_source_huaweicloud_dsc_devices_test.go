package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDevices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_devices.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDevices_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "devices.#"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.manage_url"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "devices.0.related_datasource_policy_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceDscDevices_basic = `
data "huaweicloud_dsc_devices" "test" {}
`

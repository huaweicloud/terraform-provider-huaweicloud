package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDeviceProxies_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_device_proxies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Only standard and enterprise IoTDA instances support this resource.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceProxies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxies.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxies.0.space_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxies.0.effective_time_range.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxies.0.effective_time_range.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "proxies.0.effective_time_range.0.end_time"),

					resource.TestCheckOutput("is_space_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceProxies_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_device_proxies" "test" {
  depends_on = [
    huaweicloud_iotda_device_proxy.test,
  ]
}

# Filter using space_id.
locals {
  space_id = data.huaweicloud_iotda_device_proxies.test.proxies[0].space_id
}

data "huaweicloud_iotda_device_proxies" "space_id_filter" {
  space_id = local.space_id
}

output "is_space_id_filter_useful" {
  value = length(data.huaweicloud_iotda_device_proxies.space_id_filter.proxies) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_proxies.space_id_filter.proxies[*].space_id : v == local.space_id]
  )
}

# Filter using name.
locals {
  name = data.huaweicloud_iotda_device_proxies.test.proxies[0].name
}

data "huaweicloud_iotda_device_proxies" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_iotda_device_proxies.name_filter.proxies) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_proxies.name_filter.proxies[*].name : v == local.name]
  )
}

# Filter using non existent name.
data "huaweicloud_iotda_device_proxies" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_device_proxies.not_found.proxies) == 0
}
`, testDeviceProxy_basic(name))
}

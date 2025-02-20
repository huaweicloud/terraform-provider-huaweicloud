package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDevices_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_devices.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDevices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.space_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.space_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.product_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.product_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "devices.0.status"),

					resource.TestCheckOutput("is_product_id_filter_useful", "true"),
					resource.TestCheckOutput("is_node_id_filter_useful", "true"),
					resource.TestCheckOutput("is_device_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceDevices_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_spaces" "test" {
  is_default = true
}

resource "huaweicloud_iotda_product" "test" {
  name        = "%[2]s"
  device_type = "test"
  protocol    = "MQTT"
  space_id    = data.huaweicloud_iotda_spaces.test.spaces[0].id
  data_type   = "json"

  services {
    id   = "service_1"
    type = "serv_type"
  }
}

resource "huaweicloud_iotda_device" "test" {
  node_id    = "%[2]s"
  name       = "%[2]s"
  space_id   = data.huaweicloud_iotda_spaces.test.spaces[0].id
  product_id = huaweicloud_iotda_product.test.id
}
`, buildIoTDAEndpoint(), name)
}

func testAccDataSourceDevices_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_devices" "test" {
  depends_on = [huaweicloud_iotda_device.test]
}

# Filter using product ID.
locals {
  product_id = data.huaweicloud_iotda_devices.test.devices[0].product_id
}

data "huaweicloud_iotda_devices" "product_id_filter" {
  product_id = local.product_id
}

output "is_product_id_filter_useful" {
  value = length(data.huaweicloud_iotda_devices.product_id_filter.devices) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_devices.product_id_filter.devices[*].product_id : v == local.product_id]
  )
}

# Filter using node ID.
locals {
  node_id = data.huaweicloud_iotda_devices.test.devices[0].node_id
}

data "huaweicloud_iotda_devices" "node_id_filter" {
  node_id = local.node_id
}

output "is_node_id_filter_useful" {
  value = length(data.huaweicloud_iotda_devices.node_id_filter.devices) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_devices.node_id_filter.devices[*].node_id : v == local.node_id]
  )
}

# Filter using device ID.
locals {
  device_id = data.huaweicloud_iotda_devices.test.devices[0].id
}

data "huaweicloud_iotda_devices" "device_id_filter" {
  device_id = local.device_id
}

output "is_device_id_filter_useful" {
  value = length(data.huaweicloud_iotda_devices.device_id_filter.devices) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_devices.device_id_filter.devices[*].id : v == local.device_id]
  )
}

# Filter using device name.
locals {
  name = data.huaweicloud_iotda_devices.test.devices[0].name
}

data "huaweicloud_iotda_devices" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_iotda_devices.name_filter.devices) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_devices.name_filter.devices[*].name : v == local.name]
  )
}

# Filter using non existent device name.
data "huaweicloud_iotda_devices" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_devices.not_found.devices) == 0
}
`, testDataSourceDevices_base())
}

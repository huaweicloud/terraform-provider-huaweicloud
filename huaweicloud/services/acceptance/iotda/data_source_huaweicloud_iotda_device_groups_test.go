package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDeviceGroups_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_device_groups.test"
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
				Config: testAccDataSourceDeviceGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.description"),

					resource.TestCheckOutput("group_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceGroups_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_spaces" "test" {
  is_default = "true"
}

resource "huaweicloud_iotda_product" "test" {
  name        = "%[2]s"
  device_type = "test"
  protocol    = "MQTT"
  space_id    = data.huaweicloud_iotda_spaces.test.spaces[0].id
  data_type   = "json"

  services {
    id     = "service_1"
    type   = "serv_type"
    option = "Master"
  }
}

resource "huaweicloud_iotda_device_group" "test" {
  name        = "%[2]s"
  space_id    = data.huaweicloud_iotda_spaces.test.spaces[0].id
  description = "description test"
}
`, buildIoTDAEndpoint(), name)
}

func testAccDataSourceDeviceGroups_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_device_groups" "test" {
  depends_on = [
    huaweicloud_iotda_device_group.test
  ]
}

locals {
  group_id = data.huaweicloud_iotda_device_groups.test.groups[0].id
}

data "huaweicloud_iotda_device_groups" "group_id_filter" {
  group_id = local.group_id
}

output "group_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_groups.group_id_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_groups.group_id_filter.groups[*].id : v == local.group_id]
  )
}

locals {
  name = data.huaweicloud_iotda_device_groups.test.groups[0].name
}

data "huaweicloud_iotda_device_groups" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_groups.name_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_groups.name_filter.groups[*].name : v == local.name]
  )
}

locals {
  type = data.huaweicloud_iotda_device_groups.test.groups[0].type
}

data "huaweicloud_iotda_device_groups" "type_filter" {
  type = local.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_groups.type_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_groups.type_filter.groups[*].type : v == local.type]
  )
}

data "huaweicloud_iotda_device_groups" "not_found" {
  name = "not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_device_groups.not_found.groups) == 0
}
`, testAccDataSourceDeviceGroups_base())
}

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
		deviceName     = acceptance.RandomAccResourceName()
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceGroups_basic(name, deviceName),
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

func TestAccDataSourceDeviceGroups_derived(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_device_groups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		deviceName     = acceptance.RandomAccResourceName()
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDeviceGroups_basic(name, deviceName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "groups.0.description"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDeviceGroups_basic(name, deviceName string) string {
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
`, testDeviceGroup_basic(name, deviceName))
}

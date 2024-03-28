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
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.description"),

					resource.TestCheckOutput("device_group_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("group_type_filter_is_useful", "true"),
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
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "device_groups.0.description"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("group_type_filter_is_useful", "true"),
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
  device_group_id = data.huaweicloud_iotda_device_groups.test.device_groups[0].id
}

data "huaweicloud_iotda_device_groups" "device_group_id_filter" {
  device_group_id = local.device_group_id
}

output "device_group_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_groups.device_group_id_filter.device_groups) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_groups.device_group_id_filter.device_groups[*].id : v == local.device_group_id]
  )
}

locals {
  name = data.huaweicloud_iotda_device_groups.test.device_groups[0].name
}

data "huaweicloud_iotda_device_groups" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_groups.name_filter.device_groups) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_groups.name_filter.device_groups[*].name : v == local.name]
  )
}

locals {
  group_type = data.huaweicloud_iotda_device_groups.test.device_groups[0].group_type
}

data "huaweicloud_iotda_device_groups" "group_type_filter" {
  group_type = local.group_type
}

output "group_type_filter_is_useful" {
  value = length(data.huaweicloud_iotda_device_groups.group_type_filter.device_groups) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_device_groups.group_type_filter.device_groups[*].group_type : v == local.group_type]
  )
}

data "huaweicloud_iotda_device_groups" "not_found" {
  name = "not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_iotda_device_groups.not_found.device_groups) == 0
}
`, testDeviceGroup_basic(name, deviceName))
}

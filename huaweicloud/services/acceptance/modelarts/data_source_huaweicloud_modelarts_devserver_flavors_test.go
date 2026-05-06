package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDevServerFlavors_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_devserver_flavors.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byServerType   = "data.huaweicloud_modelarts_devserver_flavors.filter_by_server_type"
		dcByServerType = acceptance.InitDataSourceCheck(byServerType)

		byArch   = "data.huaweicloud_modelarts_devserver_flavors.filter_by_arch"
		dcByArch = acceptance.InitDataSourceCheck(byArch)

		byChargingMode   = "data.huaweicloud_modelarts_devserver_flavors.filter_by_charging_mode"
		dcByChargingMode = acceptance.InitDataSourceCheck(byChargingMode)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDevServerFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "flavors.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.specification"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.arch"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.server_type"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.sku_code"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.charging_mode"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.roce_num"),
					resource.TestCheckResourceAttrSet(all, "flavors.0.flavor_type"),

					dcByServerType.CheckResourceExists(),
					resource.TestCheckOutput("is_server_type_filter_useful", "true"),

					dcByArch.CheckResourceExists(),
					resource.TestCheckOutput("is_arch_filter_useful", "true"),

					dcByChargingMode.CheckResourceExists(),
					resource.TestCheckOutput("is_charging_mode_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDevServerFlavors_basic() string {
	return `
data "huaweicloud_modelarts_devserver_flavors" "test" {}

# Filter by server_type
locals {
  server_type = data.huaweicloud_modelarts_devserver_flavors.test.flavors[*].server_type
}

data "huaweicloud_modelarts_devserver_flavors" "filter_by_server_type" {
  server_type = local.server_type[0]
}

locals {
  server_type_filter_result = [
    for v in data.huaweicloud_modelarts_devserver_flavors.filter_by_server_type.flavors[*].server_type :
    v == local.server_type[0]
  ]
}

output "is_server_type_filter_useful" {
  value = length(local.server_type_filter_result) > 0 && alltrue(local.server_type_filter_result)
}

# Filter by arch
locals {
  arch = data.huaweicloud_modelarts_devserver_flavors.test.flavors[*].arch
}

data "huaweicloud_modelarts_devserver_flavors" "filter_by_arch" {
  arch = local.arch[0]
}

locals {
  arch_filter_result = [
    for v in data.huaweicloud_modelarts_devserver_flavors.filter_by_arch.flavors[*].arch :
    v == local.arch[0]
  ]
}

output "is_arch_filter_useful" {
  value = length(local.arch_filter_result) > 0 && alltrue(local.arch_filter_result)
}

# Filter by charging_mode
locals {
  charging_mode = data.huaweicloud_modelarts_devserver_flavors.test.flavors[*].charging_mode
}

data "huaweicloud_modelarts_devserver_flavors" "filter_by_charging_mode" {
  charging_mode = local.charging_mode[0]
}

locals {
  charging_mode_filter_result = [
    for v in data.huaweicloud_modelarts_devserver_flavors.filter_by_charging_mode.flavors[*].charging_mode :
    v == local.charging_mode[0]
  ]
}

output "is_charging_mode_filter_useful" {
  value = length(local.charging_mode_filter_result) > 0 && alltrue(local.charging_mode_filter_result)
}
`
}

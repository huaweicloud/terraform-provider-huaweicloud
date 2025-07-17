package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppFlavors_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_workspace_app_flavors.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byProductId   = "data.huaweicloud_workspace_app_flavors.filter_by_product_id"
		dcByProductId = acceptance.InitDataSourceCheck(byProductId)

		byFlavorId   = "data.huaweicloud_workspace_app_flavors.filter_by_flavor_id"
		dcByFlavorId = acceptance.InitDataSourceCheck(byFlavorId)

		byAvailabilityZone   = "data.huaweicloud_workspace_app_flavors.filter_by_az"
		dcByAvailabilityZone = acceptance.InitDataSourceCheck(byAvailabilityZone)

		byChargeMode   = "data.huaweicloud_workspace_app_flavors.filter_by_charge_mode"
		dcByChargeMode = acceptance.InitDataSourceCheck(byChargeMode)

		byOsType   = "data.huaweicloud_workspace_app_flavors.filter_by_os_type"
		dcByOsType = acceptance.InitDataSourceCheck(byOsType)

		byArchitectureX86   = "data.huaweicloud_workspace_app_flavors.filter_by_architecture_x86"
		dcByArchitectureX86 = acceptance.InitDataSourceCheck(byArchitectureX86)

		byArchitectureARM   = "data.huaweicloud_workspace_app_flavors.filter_by_architecture_arm"
		dcByArchitectureARM = acceptance.InitDataSourceCheck(byArchitectureARM)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr("data.huaweicloud_workspace_app_flavors.test", "flavors.#",
						regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByProductId.CheckResourceExists(),
					resource.TestCheckOutput("is_product_id_filter_useful", "true"),
					dcByFlavorId.CheckResourceExists(),
					resource.TestCheckOutput("is_flavor_id_filter_useful", "true"),
					dcByAvailabilityZone.CheckResourceExists(),
					resource.TestCheckOutput("is_az_filter_useful", "true"),
					dcByChargeMode.CheckResourceExists(),
					resource.TestCheckOutput("is_charge_mode_filter_useful", "true"),
					dcByOsType.CheckResourceExists(),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
					dcByArchitectureX86.CheckResourceExists(),
					resource.TestCheckOutput("is_architecture_x86_filter_useful", "true"),
					dcByArchitectureARM.CheckResourceExists(),
					resource.TestCheckOutput("is_architecture_arm_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAppFlavors_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_flavors" "test" {}

# By product ID filter
locals {
  product_id = data.huaweicloud_workspace_app_flavors.test.flavors[0].product_id
}

data "huaweicloud_workspace_app_flavors" "filter_by_product_id" {
  product_id = local.product_id
}

output "is_product_id_filter_useful" {
  value = length(data.huaweicloud_workspace_app_flavors.filter_by_product_id.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_app_flavors.filter_by_product_id.flavors[*].product_id : v == local.product_id]
  )
}

# By flavor ID filter
locals {
  flavor_id = data.huaweicloud_workspace_app_flavors.test.flavors[0].id
}

data "huaweicloud_workspace_app_flavors" "filter_by_flavor_id" {
  flavor_id = local.flavor_id
}

output "is_flavor_id_filter_useful" {
  value = length(data.huaweicloud_workspace_app_flavors.filter_by_flavor_id.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_app_flavors.filter_by_flavor_id.flavors[*].id : v == local.flavor_id]
  )
}

# By charge mode filter
data "huaweicloud_workspace_app_flavors" "filter_by_charge_mode" {
  charge_mode = "1"
}

output "is_charge_mode_filter_useful" {
  value = length(data.huaweicloud_workspace_app_flavors.filter_by_charge_mode.flavors) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_app_flavors.filter_by_charge_mode.flavors[*].charge_mode : v == "1"]
  )
}

# By availability zone filter
data "huaweicloud_workspace_app_flavors" "filter_by_az" {
  availability_zone = "%[1]s"
}

output "is_az_filter_useful" {
  value = data.huaweicloud_workspace_app_flavors.filter_by_az.availability_zone == "%[1]s"
}

# By OS type filter
data "huaweicloud_workspace_app_flavors" "filter_by_os_type" {
  os_type = "Windows"
}

output "is_os_type_filter_useful" {
  value = data.huaweicloud_workspace_app_flavors.filter_by_os_type.os_type == "Windows"
}

# By architecture x86 filter
data "huaweicloud_workspace_app_flavors" "filter_by_architecture_x86" {
  architecture = "x86"
}

output "is_architecture_x86_filter_useful" {
  value = data.huaweicloud_workspace_app_flavors.filter_by_architecture_x86.architecture == "x86"
}

# By architecture arm filter
data "huaweicloud_workspace_app_flavors" "filter_by_architecture_arm" {
  architecture = "arm"
}

output "is_architecture_arm_filter_useful" {
  value = data.huaweicloud_workspace_app_flavors.filter_by_architecture_arm.architecture == "arm"
}
`, acceptance.HW_AVAILABILITY_ZONE)
}

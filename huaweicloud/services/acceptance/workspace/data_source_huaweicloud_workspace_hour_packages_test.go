package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccHourPackagesDataSource_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_hour_packages.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterByDesktopResourceSpecCodeName   = "data.huaweicloud_workspace_hour_packages.filter_by_desktop_resource_spec_code"
		dcFilterByDesktopResourceSpecCodeName = acceptance.InitDataSourceCheck(filterByDesktopResourceSpecCodeName)

		filterByResourceSpecCodeName   = "data.huaweicloud_workspace_hour_packages.filter_by_resource_spec_code"
		dcFilterByResourceSpecCodeName = acceptance.InitDataSourceCheck(filterByResourceSpecCodeName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHourPackagesDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "hour_packages.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.resource_type"),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.desktop_resource_spec_code"),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.descriptions.0.zh_cn"),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.descriptions.0.en_us"),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.package_duration"),
					resource.TestCheckResourceAttrSet(dcName, "hour_packages.0.status"),
					// filter by desktop resource spec code
					dcFilterByDesktopResourceSpecCodeName.CheckResourceExists(),
					resource.TestCheckOutput("is_desktop_resource_spec_code_filter_useful", "true"),
					// filter by resource spec code
					dcFilterByResourceSpecCodeName.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_spec_code_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccHourPackagesDataSource_basic = `
data "huaweicloud_workspace_hour_packages" "all" {}

locals {
  desktop_resource_spec_code = try(data.huaweicloud_workspace_hour_packages.all.hour_packages[0].desktop_resource_spec_code, "NOT_FOUND")
  resource_spec_code         = try(data.huaweicloud_workspace_hour_packages.all.hour_packages[0].resource_spec_code, "NOT_FOUND")
}

# Filter by desktop resource spec code
data "huaweicloud_workspace_hour_packages" "filter_by_desktop_resource_spec_code" {
  desktop_resource_spec_code = local.desktop_resource_spec_code
}

locals {
  desktop_resource_spec_code_filter_result = [
    for v in data.huaweicloud_workspace_hour_packages.filter_by_desktop_resource_spec_code.hour_packages[*].desktop_resource_spec_code
    : v == local.desktop_resource_spec_code
  ]
}

output "is_desktop_resource_spec_code_filter_useful" {
  value = length(local.desktop_resource_spec_code_filter_result) > 0 && alltrue(local.desktop_resource_spec_code_filter_result)
}

# Filter by resource spec code
data "huaweicloud_workspace_hour_packages" "filter_by_resource_spec_code" {
  resource_spec_code = local.resource_spec_code
}

locals {
  resource_spec_code_filter_result = [
    for v in data.huaweicloud_workspace_hour_packages.filter_by_resource_spec_code.hour_packages[*].resource_spec_code
    : v == local.resource_spec_code
  ]
}

output "is_resource_spec_code_filter_useful" {
  value = length(local.resource_spec_code_filter_result) > 0 && alltrue(local.resource_spec_code_filter_result)
}
`

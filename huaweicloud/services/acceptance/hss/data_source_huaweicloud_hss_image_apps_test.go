package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageApps_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_apps.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageApps_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.app_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.app_version"),

					resource.TestCheckOutput("is_image_type_filter_useful", "true"),
					resource.TestCheckOutput("is_app_name_filter_useful", "true"),
					resource.TestCheckOutput("is_is_compliant_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceImageApps_basic() string {
	return `
data "huaweicloud_hss_image_apps" "test" {}

# Filter using image_type.
data "huaweicloud_hss_image_apps" "image_type_filter" {
  image_type = "registry"
}

output "is_image_type_filter_useful" {
  value = length(data.huaweicloud_hss_image_apps.image_type_filter.data_list) > 0
}

# Filter using app_name.
locals {
  app_name = data.huaweicloud_hss_image_apps.test.data_list[0].app_name
}

data "huaweicloud_hss_image_apps" "app_name_filter" {
  app_name = local.app_name
}

output "is_app_name_filter_useful" {
  value = length(data.huaweicloud_hss_image_apps.app_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_apps.app_name_filter.data_list[*].app_name : v == local.app_name]
  )
}

# Filter using is_compliant.
data "huaweicloud_hss_image_apps" "is_compliant_filter" {
  is_compliant = true
}

output "is_is_compliant_filter_useful" {
  value = length(data.huaweicloud_hss_image_apps.is_compliant_filter.data_list) > 0
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_image_apps" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_image_apps.enterprise_project_id_filter.data_list) > 0
}
`
}

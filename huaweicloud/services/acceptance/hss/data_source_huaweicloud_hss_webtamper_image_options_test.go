package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWebTamperImageOptions_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_webtamper_image_options.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWebTamperImageOptions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_full_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_version_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.image_namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.registry_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.registry_type"),

					resource.TestCheckOutput("is_image_namespace_filter_useful", "true"),
					resource.TestCheckOutput("is_registry_name_filter_useful", "true"),
					resource.TestCheckOutput("is_image_name_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceWebTamperImageOptions_basic() string {
	return `
data "huaweicloud_hss_webtamper_image_options" "test" {
  image_type = "registry"
}

# Filter using image_namespace.
locals {
  image_namespace = data.huaweicloud_hss_webtamper_image_options.test.data_list[0].image_namespace
}

data "huaweicloud_hss_webtamper_image_options" "image_namespace_filter" {
  image_type      = "registry"
  image_namespace = local.image_namespace
}

output "is_image_namespace_filter_useful" {
  value = length(data.huaweicloud_hss_webtamper_image_options.image_namespace_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_webtamper_image_options.image_namespace_filter.data_list[*].image_namespace : v == local.image_namespace]
  )
}

# Filter using registry_name.
locals {
  registry_name = data.huaweicloud_hss_webtamper_image_options.test.data_list[0].registry_name
}

data "huaweicloud_hss_webtamper_image_options" "registry_name_filter" {
  image_type    = "registry"
  registry_name = local.registry_name
}

output "is_registry_name_filter_useful" {
  value = length(data.huaweicloud_hss_webtamper_image_options.registry_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_webtamper_image_options.registry_name_filter.data_list[*].registry_name : v == local.registry_name]
  )
}

# Filter using image_name.
locals {
  image_name = data.huaweicloud_hss_webtamper_image_options.test.data_list[0].image_name
}

data "huaweicloud_hss_webtamper_image_options" "image_name_filter" {
  image_type = "registry"
  image_name = local.image_name
}

output "is_image_name_filter_useful" {
  value = length(data.huaweicloud_hss_webtamper_image_options.image_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_webtamper_image_options.image_name_filter.data_list[*].image_name : v == local.image_name]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_webtamper_image_options" "enterprise_project_id_filter" {
  image_type            = "registry"
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_webtamper_image_options.test.data_list) > 0
}
`
}

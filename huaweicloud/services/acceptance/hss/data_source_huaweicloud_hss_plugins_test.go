package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePlugins_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_plugins.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePlugins_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.installed_attachment_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.uninstall_attachment_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.max_cpu_limit"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.max_memory_limit"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.max_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.display_mode"),

					resource.TestCheckOutput("is_code_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePlugins_basic() string {
	return `
data "huaweicloud_hss_plugins" "test" {}

# Filter using code.
locals {
  code = data.huaweicloud_hss_plugins.test.data_list[0].code
}

data "huaweicloud_hss_plugins" "code_filter" {
  code = local.code
}

output "is_code_filter_useful" {
  value = length(data.huaweicloud_hss_plugins.code_filter.data_list) > 0 && alltrue(
  [for v in data.huaweicloud_hss_plugins.code_filter.data_list[*].code : v == local.code]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_plugins" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_plugins.test.data_list) > 0
}
`
}

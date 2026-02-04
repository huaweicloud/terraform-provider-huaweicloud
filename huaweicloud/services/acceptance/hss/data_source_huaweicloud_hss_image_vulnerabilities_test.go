package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageVulnerabilities_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_image_vulnerabilities.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageVulnerabilities_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.vul_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.vul_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.repair_necessity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.solution"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.history_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.undeal_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.data_list.#"),

					resource.TestCheckOutput("repair_necessity_filter_is_useful", "true"),
					resource.TestCheckOutput("vul_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceImageVulnerabilities_basic() string {
	return `
data "huaweicloud_hss_image_vulnerabilities" "test" {
  enterprise_project_id = "all_granted_eps"
}

# Filter using repair_necessity.
locals {
  repair_necessity = data.huaweicloud_hss_image_vulnerabilities.test.data_list[0].repair_necessity
}

data "huaweicloud_hss_image_vulnerabilities" "repair_necessity_filter" {
  repair_necessity = local.repair_necessity
}

output "repair_necessity_filter_is_useful" {
  value = length(data.huaweicloud_hss_image_vulnerabilities.repair_necessity_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_vulnerabilities.repair_necessity_filter.data_list[*].repair_necessity : v == local.repair_necessity]
  )
}

# Filter using vul_id.
locals {
  vul_id = data.huaweicloud_hss_image_vulnerabilities.test.data_list[0].vul_id
}

data "huaweicloud_hss_image_vulnerabilities" "vul_id_filter" {
  vul_id = local.vul_id
}

output "vul_id_filter_is_useful" {
  value = length(data.huaweicloud_hss_image_vulnerabilities.vul_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_vulnerabilities.vul_id_filter.data_list[*].vul_id : v == local.vul_id]
  )
}

# Filter using type.
data "huaweicloud_hss_image_vulnerabilities" "type_filter" {
  type = "linux_vul"
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_hss_image_vulnerabilities.type_filter.data_list) > 0
}
`
}

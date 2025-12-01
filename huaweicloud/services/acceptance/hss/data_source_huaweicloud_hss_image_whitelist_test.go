package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageWhitelists_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_image_whitelists.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need an image whitelist.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceImageWhitelists_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rule_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cves.#"),

					resource.TestCheckOutput("vul_name_filter_useful", "true"),
					resource.TestCheckOutput("vul_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceImageWhitelists_basic = `
data "huaweicloud_hss_image_whitelists" "test" {
  global_image_type = "registry"
  type              = "vulnerability"
	
}

locals {
  vul_name = data.huaweicloud_hss_image_whitelists.test.data_list[0].vul_name
  vul_type = data.huaweicloud_hss_image_whitelists.test.data_list[0].vul_type
}

data "huaweicloud_hss_image_whitelists" "vul_name_filter" {
  global_image_type = "registry"
  type              = "vulnerability"
  vul_name          = local.vul_name
}

output "vul_name_filter_useful" {
  value = length(data.huaweicloud_hss_image_whitelists.vul_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_whitelists.vul_name_filter.data_list[*].vul_name : v == local.vul_name]
  )
}

data "huaweicloud_hss_image_whitelists" "vul_type_filter" {
  global_image_type = "registry"
  type              = "vulnerability"
  vul_type          = local.vul_type
}

output "vul_type_filter_useful" {
  value = length(data.huaweicloud_hss_image_whitelists.vul_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_whitelists.vul_type_filter.data_list[*].vul_type : v == local.vul_type]
  )
}
`

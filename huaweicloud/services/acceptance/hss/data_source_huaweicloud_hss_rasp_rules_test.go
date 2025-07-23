package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRaspRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_rasp_rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceRaspRules_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.chk_feature_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.chk_feature_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.chk_feature_desc"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.feature_configure"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.optional_protective_action"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protective_action"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.editable"),

					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceRaspRules_basic = `
data "huaweicloud_hss_rasp_rules" "test" {}

locals {
  os_type = data.huaweicloud_hss_rasp_rules.test.data_list[0].os_type
}

data "huaweicloud_hss_rasp_rules" "os_type_filter" {
  os_type  = local.os_type
}

output "is_os_type_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_rules.os_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_rasp_rules.os_type_filter.data_list[*].os_type : v == local.os_type]
  )
}

data "huaweicloud_hss_rasp_rules" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_rules.eps_filter.data_list) > 0
}
`

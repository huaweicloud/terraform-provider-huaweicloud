package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAllDataMaskingRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_all_data_masking_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Prepare a WAF data masking rule before test
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAllDataMaskingRules_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.policyid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.timestamp"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.index"),

					resource.TestCheckOutput("policyid_filter_is_useful", "true"),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAllDataMaskingRules_basic = `
data "huaweicloud_waf_all_data_masking_rules" "test" {}

locals {
  policyid = data.huaweicloud_waf_all_data_masking_rules.test.items[0].policyid
}

data "huaweicloud_waf_all_data_masking_rules" "filter_by_policyid" {
  policyids = local.policyid
}

output "policyid_filter_is_useful" {
  value = length(data.huaweicloud_waf_all_data_masking_rules.filter_by_policyid.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_all_data_masking_rules.filter_by_policyid.items[*].policyid : v == local.policyid]
  )
}

data "huaweicloud_waf_all_data_masking_rules" "filter_by_eps" {
  enterprise_project_id = "0"
}

output "eps_filter_is_useful" {
  value = length(data.huaweicloud_waf_all_data_masking_rules.filter_by_eps.items) > 0
}
`

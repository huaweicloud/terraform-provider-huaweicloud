package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcSiteNetworkQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_site_network_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcSiteNetworkQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_key"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.unit"),
					resource.TestCheckOutput("quota_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcSiteNetworkQuotas_basic() string {
	return `
data "huaweicloud_cc_site_network_quotas" "test" {}

data "huaweicloud_cc_site_network_quotas" "quota_type_filter" {
  quota_type = [data.huaweicloud_cc_site_network_quotas.test.quotas[0].quota_key]
}
locals {
  quota_type = data.huaweicloud_cc_site_network_quotas.test.quotas[0].quota_key
}
output "quota_type_filter_is_useful" {
  value = length(data.huaweicloud_cc_site_network_quotas.quota_type_filter.quotas) > 0 && alltrue(
  [for v in data.huaweicloud_cc_site_network_quotas.quota_type_filter.quotas[*].quota_key : v == local.quota_type]
  )
}
`
}

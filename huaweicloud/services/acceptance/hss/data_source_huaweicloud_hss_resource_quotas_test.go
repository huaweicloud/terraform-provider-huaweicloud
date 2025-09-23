package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_resource_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.used_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.available_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.available_resources_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.available_resources_list.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.available_resources_list.0.current_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.available_resources_list.0.shared_quota"),

					resource.TestCheckOutput("is_version_filter_useful", "true"),
					resource.TestCheckOutput("is_charging_mode_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceResourceQuotas_basic string = `
data "huaweicloud_hss_resource_quotas" "test" {}

# Filter using version.
locals {
  version = data.huaweicloud_hss_resource_quotas.test.data_list[0].version
}

data "huaweicloud_hss_resource_quotas" "version_filter" {
  version = local.version
}

output "is_version_filter_useful" {
  value = length(data.huaweicloud_hss_resource_quotas.version_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_resource_quotas.version_filter.data_list[*].version : v == local.version]
  )
}

# Filter using charging mode.
data "huaweicloud_hss_resource_quotas" "charging_mode_filter" {
  charging_mode = "packet_cycle"
}

output "is_charging_mode_filter_useful" {
  value = length(data.huaweicloud_hss_resource_quotas.charging_mode_filter.data_list) > 0
}
`

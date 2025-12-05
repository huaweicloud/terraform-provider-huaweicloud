package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppServerQuotas_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_app_server_quotas.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byFlavorId   = "data.huaweicloud_workspace_app_server_quotas.filter_by_flavor_id"
		dcByFlavorId = acceptance.InitDataSourceCheck(byFlavorId)

		byIsPeriod   = "data.huaweicloud_workspace_app_server_quotas.filter_by_is_period"
		dcByIsPeriod = acceptance.InitDataSourceCheck(byIsPeriod)
		// The deh_id and cluster_id parameters have no test scenarios.
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataAppServerQuotas_notFound,
				ExpectError: regexp.MustCompile("The product 'not_exist' not find"),
			},
			{
				Config: testAccDataAppServerQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(all, "is_enough", "true"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.type"),
					// Filter by 'flavor_id' parameter.
					dcByFlavorId.CheckResourceExists(),
					resource.TestCheckOutput("is_flavor_id_filter_useful", "true"),
					// Filter by 'is_period' parameter.
					dcByIsPeriod.CheckResourceExists(),
					resource.TestMatchResourceAttr(byIsPeriod, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

const testAccDataAppServerQuotas_notFound = `
data "huaweicloud_workspace_app_server_quotas" "test" {
  product_id       = "not_exist"
  subscription_num = 1
  disk_size        = 80
  disk_num         = 1
}
`

const testAccDataAppServerQuotas_basic = `
data "huaweicloud_workspace_app_flavors" "test" {}

locals {
  product_id = try(data.huaweicloud_workspace_app_flavors.test.flavors[0].product_id, null)
  disk_size  = try(data.huaweicloud_workspace_app_flavors.test.flavors[0].system_disk_size, null)
}

data "huaweicloud_workspace_app_server_quotas" "all" {
  product_id       = local.product_id
  subscription_num = 1
  disk_size        = local.disk_size
  disk_num         = 1
}

data "huaweicloud_workspace_app_server_quotas" "filter_by_flavor_id" {
  product_id       = local.product_id
  subscription_num = 1
  disk_size        = try(data.huaweicloud_workspace_app_flavors.test.flavors[0].system_disk_size, null)
  disk_num         = 1
  flavor_id        = try(data.huaweicloud_workspace_app_flavors.test.flavors[0].id, null)
}

locals {
  is_flavor_id_filter_result = [for v in data.huaweicloud_workspace_app_server_quotas.filter_by_flavor_id.quotas :
  v if v.type == "VOLUMES"]
}

output "is_flavor_id_filter_useful" {
  value = alltrue(
    [
      length(local.is_flavor_id_filter_result) > 0,
      local.is_flavor_id_filter_result[0].remainder > 0,
      local.is_flavor_id_filter_result[0].need > 0
    ]
  )
}

data "huaweicloud_workspace_app_server_quotas" "filter_by_is_period" {
  product_id       = local.product_id
  subscription_num = 1
  disk_size        = local.disk_size
  disk_num         = 1
  is_period        = true
}
`

package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcEipPools_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_eip_pools.all"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// there is no openAPI supporting to create EIP pool
			acceptance.TestAccPreCheckVpcEipPoolEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcEipPools_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "pools.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.allow_share_bandwidth_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.shared"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.billing_info.0.product_id"),
					resource.TestCheckResourceAttrSet(dataSource, "pools.0.billing_info.0.order_id"),

					resource.TestCheckOutput("name_filter_validation", "true"),
					resource.TestCheckOutput("size_filter_validation", "true"),
					resource.TestCheckOutput("type_filter_validation", "true"),
					resource.TestCheckOutput("status_filter_validation", "true"),
					resource.TestCheckOutput("public_border_group_filter_validation", "true"),
				),
			},
		},
	})
}

const testDataSourceVpcEipPools_basic = `
data "huaweicloud_vpc_eip_pools" "all" {}

data "huaweicloud_vpc_eip_pools" "test" {
  name                = local.test_refer.name
  size                = local.test_refer.size
  status              = local.test_refer.status
  type                = local.test_refer.type
  public_border_group = local.test_refer.public_border_group
}

locals {
  test_refer   = data.huaweicloud_vpc_eip_pools.all.pools[0]
  test_results = data.huaweicloud_vpc_eip_pools.test
}

output "name_filter_validation" {
  value = alltrue([for v in local.test_results.pools[*].name : v == local.test_refer.name])
}

output "size_filter_validation" {
  value = alltrue([for v in local.test_results.pools[*].size : v == local.test_refer.size])
}

output "type_filter_validation" {
  value = alltrue([for v in local.test_results.pools[*].type : v == local.test_refer.type])
}

output "status_filter_validation" {
  value = alltrue([for v in local.test_results.pools[*].status : v == local.test_refer.status])
}

output "public_border_group_filter_validation" {
  value = alltrue([for v in local.test_results.pools[*].public_border_group : v == local.test_refer.public_border_group])
}
`

package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcEipPublicIpPoolTypes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_vpc_eip_public_ip_pool_types.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcEipPublicIpPoolTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "public_ip_pool_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "public_ip_pool_types.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "public_ip_pool_types.0.type"),

					resource.TestCheckOutput("is_type_id_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceVpcEipPublicIpPoolTypes_basic() string {
	return `
data "huaweicloud_vpc_eip_public_ip_pool_types" "test" {
  fields   = ["id", "type"]
  sort_key = "id"
  sort_dir = "asc"
}

# Filter by type_id.
locals {
  type_id = data.huaweicloud_vpc_eip_public_ip_pool_types.test.public_ip_pool_types[0].id
}

data "huaweicloud_vpc_eip_public_ip_pool_types" "type_id_filter" {
  type_id = local.type_id
}

locals {
  type_id_filter_result = [
    for v in data.huaweicloud_vpc_eip_public_ip_pool_types.type_id_filter.public_ip_pool_types[*].id : v == local.type_id
  ]
}

output "is_type_id_filter_useful" {
  value = alltrue(local.type_id_filter_result) && length(local.type_id_filter_result) > 0
}

# Filter by type.
locals {
  type = data.huaweicloud_vpc_eip_public_ip_pool_types.test.public_ip_pool_types[0].type
}

data "huaweicloud_vpc_eip_public_ip_pool_types" "type_filter" {
  type = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_vpc_eip_public_ip_pool_types.type_id_filter.public_ip_pool_types[*].type : v == local.type
  ]
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`
}

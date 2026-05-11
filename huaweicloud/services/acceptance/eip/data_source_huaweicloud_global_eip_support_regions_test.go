package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipSupportRegions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eip_support_regions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalEipSupportRegions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.instance_type"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.access_site"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.remote_endpoint"),
					resource.TestCheckResourceAttrSet(dataSource, "support_regions.0.status"),

					resource.TestCheckOutput("is_support_region_ids_filter_useful", "true"),
					resource.TestCheckOutput("is_instance_type_filter_useful", "true"),
					resource.TestCheckOutput("is_public_border_group_filter_useful", "true"),
					resource.TestCheckOutput("is_region_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGlobalEipSupportRegions_basic() string {
	return `
data "huaweicloud_global_eip_support_regions" "test" {
  page_reverse = "true"
  fields       = ["id", "instance_type", "access_site", "region_id", "public_border_group", "remote_endpoint", "status"]
  sort_key     = "region_id"
  sort_dir     = "asc"
}

# Filter by support_region_ids.
locals {
  support_region_id = data.huaweicloud_global_eip_support_regions.test.support_regions[0].id
}

data "huaweicloud_global_eip_support_regions" "support_region_ids_filter" {
  support_region_ids = [local.support_region_id]
}

locals {
  support_region_ids_filter_result = [
    for v in data.huaweicloud_global_eip_support_regions.support_region_ids_filter.support_regions[*].id : v
    == local.support_region_id
  ]
}

output "is_support_region_ids_filter_useful" {
  value = alltrue(local.support_region_ids_filter_result) && length(local.support_region_ids_filter_result) > 0
}

# Filter by instance_type.
locals {
  instance_type = data.huaweicloud_global_eip_support_regions.test.support_regions[0].instance_type
}

data "huaweicloud_global_eip_support_regions" "instance_type_filter" {
  instance_type = [local.instance_type]
}

locals {
  instance_type_filter_result = [
    for v in data.huaweicloud_global_eip_support_regions.instance_type_filter.support_regions[*].instance_type : v
    == local.instance_type
  ]
}

output "is_instance_type_filter_useful" {
  value = alltrue(local.instance_type_filter_result) && length(local.instance_type_filter_result) > 0
}

# Filter by public_border_group.
locals {
  public_border_group = data.huaweicloud_global_eip_support_regions.test.support_regions[0].public_border_group
}

data "huaweicloud_global_eip_support_regions" "public_border_group_filter" {
  public_border_group = [local.public_border_group]
}

locals {
  public_border_group_filter_result = [
    for v in data.huaweicloud_global_eip_support_regions.public_border_group_filter.support_regions[*].public_border_group
    : v == local.public_border_group
  ]
}

output "is_public_border_group_filter_useful" {
  value = alltrue(local.public_border_group_filter_result) && length(local.public_border_group_filter_result) > 0
}

# Filter by region_id.
locals {
  region_id = data.huaweicloud_global_eip_support_regions.test.support_regions[0].region_id
}

data "huaweicloud_global_eip_support_regions" "region_id_filter" {
  region_id = [local.region_id]
}

locals {
  region_id_filter_result = [
    for v in data.huaweicloud_global_eip_support_regions.region_id_filter.support_regions[*].region_id : v
    == local.region_id
  ]
}

output "is_region_id_filter_useful" {
  value = alltrue(local.region_id_filter_result) && length(local.region_id_filter_result) > 0
}
`
}

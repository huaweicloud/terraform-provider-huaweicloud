package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipSegments_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eip_segments.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalEipSegments_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "global_eip_segments.#"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eip_segments.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eip_segments.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eip_segments.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eip_segments.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eip_segments.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eip_segments.0.updated_at"),

					resource.TestCheckOutput("is_segment_ids_filter_useful", "true"),
					resource.TestCheckOutput("is_ip_version_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGlobalEipSegments_basic() string {
	return `
data "huaweicloud_global_eip_segments" "test" {
  page_reverse = "true"
  sort_key     = "id"
  sort_dir     = "asc"
}

# Filter by segment_ids.
locals {
  segment_id = data.huaweicloud_global_eip_segments.test.global_eip_segments[0].id
}

data "huaweicloud_global_eip_segments" "segment_ids_filter" {
  fields      = ["id"]
  segment_ids = [local.segment_id]
}

locals {
  segment_ids_filter_result = [
    for v in data.huaweicloud_global_eip_segments.segment_ids_filter.global_eip_segments[*].id : v == local.segment_id
  ]
}

output "is_segment_ids_filter_useful" {
  value = alltrue(local.segment_ids_filter_result) && length(local.segment_ids_filter_result) > 0
}

# Filter by ip_version.
locals {
  ip_version = data.huaweicloud_global_eip_segments.test.global_eip_segments[0].ip_version
}

data "huaweicloud_global_eip_segments" "ip_version_filter" {
  ip_version = [local.ip_version]
}

locals {
  ip_version_filter_result = [
    for v in data.huaweicloud_global_eip_segments.ip_version_filter.global_eip_segments[*].ip_version : v == local.ip_version
  ]
}

output "is_ip_version_filter_useful" {
  value = alltrue(local.ip_version_filter_result) && length(local.ip_version_filter_result) > 0
}

# Filter by name.
locals {
  segment_name = data.huaweicloud_global_eip_segments.test.global_eip_segments[0].name
}

data "huaweicloud_global_eip_segments" "name_filter" {
  name = [local.segment_name]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_global_eip_segments.name_filter.global_eip_segments[*].name : v == local.segment_name
  ]
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`
}

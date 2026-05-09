package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipSegmentSupportMasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eip_segment_support_masks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGlobalEipSegmentSupportMasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "support_masks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "support_masks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "support_masks.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "support_masks.0.mask"),
					resource.TestCheckResourceAttrSet(dataSource, "support_masks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "support_masks.0.updated_at"),

					resource.TestCheckOutput("is_mask_ids_filter_useful", "true"),
					resource.TestCheckOutput("is_ip_version_filter_useful", "true"),
					resource.TestCheckOutput("is_mask_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGlobalEipSegmentSupportMasks_basic() string {
	return `
data "huaweicloud_global_eip_segment_support_masks" "test" {
  page_reverse = "true"
  fields       = ["id", "mask", "ip_version", "created_at", "updated_at"]
  sort_key     = "mask"
  sort_dir     = "asc"
}

# Filter by mask_ids.
locals {
  mask_id = data.huaweicloud_global_eip_segment_support_masks.test.support_masks[0].id
}

data "huaweicloud_global_eip_segment_support_masks" "mask_ids_filter" {
  mask_ids = [local.mask_id]
}

locals {
  mask_ids_filter_result = [
    for v in data.huaweicloud_global_eip_segment_support_masks.mask_ids_filter.support_masks[*].id : v == local.mask_id
  ]
}

output "is_mask_ids_filter_useful" {
  value = alltrue(local.mask_ids_filter_result) && length(local.mask_ids_filter_result) > 0
}

# Filter by ip_version.
locals {
  ip_version = data.huaweicloud_global_eip_segment_support_masks.test.support_masks[0].ip_version
}

data "huaweicloud_global_eip_segment_support_masks" "ip_version_filter" {
  ip_version = [local.ip_version]
}

locals {
  ip_version_filter_result = [
    for v in data.huaweicloud_global_eip_segment_support_masks.ip_version_filter.support_masks[*].ip_version : v == local.ip_version
  ]
}

output "is_ip_version_filter_useful" {
  value = alltrue(local.ip_version_filter_result) && length(local.ip_version_filter_result) > 0
}

# Filter by mask.
locals {
  mask = data.huaweicloud_global_eip_segment_support_masks.test.support_masks[0].mask
}

data "huaweicloud_global_eip_segment_support_masks" "mask_filter" {
  mask = local.mask
}

locals {
  mask_filter_result = [
    for v in data.huaweicloud_global_eip_segment_support_masks.mask_filter.support_masks[*].mask : v == local.mask
  ]
}

output "is_mask_filter_useful" {
  value = alltrue(local.mask_filter_result) && length(local.mask_filter_result) > 0
}
`
}

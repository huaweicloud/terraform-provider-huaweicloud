package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTrackedResourceTags_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_rms_tracked_resource_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTrackedResourceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.value.#"),
					resource.TestCheckOutput("key_filter_is_useful", "true"),
					resource.TestCheckOutput("resource_deleted_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTrackedResourceTags_basic() string {
	return `
data "huaweicloud_rms_tracked_resource_tags" "test" {}

data "huaweicloud_rms_tracked_resource_tags" "key_filter" {
  key = data.huaweicloud_rms_tracked_resource_tags.test.tags[0].key
}
locals {
  key = data.huaweicloud_rms_tracked_resource_tags.test.tags[0].key
}
output "key_filter_is_useful" {
  value = length(data.huaweicloud_rms_tracked_resource_tags.key_filter.tags) > 0 && alltrue(
  [for v in data.huaweicloud_rms_tracked_resource_tags.key_filter.tags[*].key : strcontains(lower(v), lower(local.key))]
  )
}

data "huaweicloud_rms_tracked_resource_tags" "resource_deleted_filter" {
  resource_deleted = true
}
output "resource_deleted_filter_is_useful" {
  value = length(data.huaweicloud_rms_tracked_resource_tags.resource_deleted_filter.tags) > 0
}
`
}

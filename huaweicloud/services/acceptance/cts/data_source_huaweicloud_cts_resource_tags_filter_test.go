package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceTagsFilter_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cts_resource_tags_filter.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceTagsFilter_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.name"),
					resource.TestCheckOutput("is_resources_exist", "true"),
					resource.TestCheckOutput("is_resource_id_filter_useful", "false"),
					resource.TestCheckOutput("is_resource_name_filter_useful", "true"),
					resource.TestCheckOutput("is_first_resource_id_set", "true"),
					resource.TestCheckOutput("is_first_resource_name_set", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceResourceTagsFilter_withTags(t *testing.T) {
	dataSourceName := "data.huaweicloud_cts_resource_tags_filter.with_tags"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceTagsFilter_withTags(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.#"),
					resource.TestCheckOutput("is_resources_exist", "true"),
				),
			},
		},
	})
}

func testAccDataSourceResourceTagsFilter_base() string {
	return `
data "huaweicloud_cts_trackers" "test" {}

locals {
  resource_id   = try(data.huaweicloud_cts_trackers.test.trackers[0].id, "")
  resource_name = try(data.huaweicloud_cts_trackers.test.trackers[0].name, "")
}
`
}

func testAccDataSourceResourceTagsFilter_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cts_resource_tags_filter" "test" {
  resource_type = "cts-tracker"

  depends_on = [data.huaweicloud_cts_trackers.test]
}

locals {
  resources = data.huaweicloud_cts_resource_tags_filter.test.resources
}

output "is_resources_exist" {
  value = length(local.resources) >= 0
}

output "is_first_resource_id_set" {
  value = length(local.resources) > 0 ? try(local.resources[0].id != "", false) : true
}

output "is_first_resource_name_set" {
  value = length(local.resources) > 0 ? try(local.resources[0].name != "", false) : true
}

# Filter by resource ID
data "huaweicloud_cts_resource_tags_filter" "filter_by_resource_id" {
  resource_type = "cts-tracker"

  matches {
    key   = "resource_name"
    value = local.resource_id
  }

  depends_on = [data.huaweicloud_cts_trackers.test]
}

locals {
  resource_id_filter_result = [
    for v in data.huaweicloud_cts_resource_tags_filter.filter_by_resource_id.resources[*].id : v == local.resource_id
  ]
}

output "is_resource_id_filter_useful" {
  value = length(local.resource_id_filter_result) > 0 && alltrue(local.resource_id_filter_result)
}

# Filter by resource name
data "huaweicloud_cts_resource_tags_filter" "filter_by_resource_name" {
  resource_type = "cts-tracker"

  matches {
    key   = "resource_name"
    value = local.resource_name
  }

  depends_on = [data.huaweicloud_cts_trackers.test]
}

locals {
  resource_name_filter_result = [
    for v in data.huaweicloud_cts_resource_tags_filter.filter_by_resource_name.resources[*].name : v == local.resource_name
  ]
}

output "is_resource_name_filter_useful" {
  value = length(local.resource_name_filter_result) > 0 && alltrue(local.resource_name_filter_result)
}
`, testAccDataSourceResourceTagsFilter_base())
}

func testAccDataSourceResourceTagsFilter_withTags() string {
	return `
data "huaweicloud_cts_resource_tags_filter" "with_tags" {
  resource_type = "cts-tracker"

  tags {
    key    = "test_key"
    values = ["test_value"]
  }
}

locals {
  resources = data.huaweicloud_cts_resource_tags_filter.with_tags.resources
}

output "is_resources_exist" {
  value = length(local.resources) >= 0
}
`
}

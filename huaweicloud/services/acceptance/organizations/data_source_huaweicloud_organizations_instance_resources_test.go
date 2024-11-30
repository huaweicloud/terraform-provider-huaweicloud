package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsResourceInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_resource_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationsResourceInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.value"),
					resource.TestCheckOutput("without_any_tag_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsResourceInstances_basic() string {
	return `
data "huaweicloud_organizations_resource_instances" "test" {
  resource_type = "organizations:accounts"
}

data "huaweicloud_organizations_resource_instances" "without_any_tag_filter" {
  resource_type   = "organizations:accounts"
  without_any_tag = true
}
	
output "without_any_tag_filter_is_useful" {
  value = length(data.huaweicloud_organizations_resource_instances.without_any_tag_filter.resources) > 0 && alltrue(
    [for v in data.huaweicloud_organizations_resource_instances.without_any_tag_filter.resources[*].tags : length(v) == 0]
  )
}

data "huaweicloud_organizations_resource_instances" "tags_filter" {
  resource_type = "organizations:accounts"

  tags {
    key    = "key3"
    values = ["value3"]
  }
}

locals {
  tag_key   = "key3"
  tag_value = "value3"
}
	
output "tags_filter_is_useful" {
  value = length(data.huaweicloud_organizations_resource_instances.tags_filter.resources) > 0 && alltrue(
    [for v in data.huaweicloud_organizations_resource_instances.tags_filter.resources[*].tags : anytrue(
    [for vv in v[*].key : vv == local.tag_key]) && anytrue([for vv in v[*].value : vv == local.tag_value])]
  )
}

data "huaweicloud_organizations_resource_instances" "matches_filter" {
  resource_type = "organizations:accounts"

  matches {
    key   = "key3"
    value = "value3"
  }
}

locals {
  match_key   = "resource_name"
  match_value = "tf_test"
}
	
output "matches_filter_is_useful" {
  value = length(data.huaweicloud_organizations_resource_instances.matches_filter.resources) > 0
}
`
}

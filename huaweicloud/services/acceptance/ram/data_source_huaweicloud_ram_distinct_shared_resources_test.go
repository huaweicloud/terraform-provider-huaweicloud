package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDistinctSharedResources_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_distinct_shared_resources.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDistinctSharedReources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "distinct_shared_resources.0.resource_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "distinct_shared_resources.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "distinct_shared_resources.0.updated_at"),

					resource.TestCheckOutput("is_resource_urns_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDistinctSharedReources_basic() string {
	return `
data "huaweicloud_ram_distinct_shared_resources" "test" {
  resource_owner = "self"
}

# Filter using resource_urns.
locals {
  resource_urn = data.huaweicloud_ram_distinct_shared_resources.test.distinct_shared_resources[0].resource_urn
}

data "huaweicloud_ram_distinct_shared_resources" "resource_urns_filter" {
  resource_owner = "self"
  resource_urns  = [local.resource_urn]
}

output "is_resource_urns_filter_useful" {
  value = length(data.huaweicloud_ram_distinct_shared_resources.resource_urns_filter.distinct_shared_resources) > 0 && alltrue(
    [for v in data.huaweicloud_ram_distinct_shared_resources.resource_urns_filter.distinct_shared_resources[*].resource_urn : v
    == local.resource_urn]
  )
}
`
}

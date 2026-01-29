package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceInstancesFilter_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_resource_instances_filter.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceInstancesFilter_without_any_tag_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
			{
				Config: testDataSourceResourceInstancesFilter_with_tags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
			{
				Config: testDataSourceResourceInstancesFilter_with_matches_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
		},
	})
}

func testDataSourceResourceInstancesFilter_without_any_tag_basic() string {
	return `
data "huaweicloud_ram_resource_instances_filter" "test" {
  without_any_tag = true
}
`
}

func testDataSourceResourceInstancesFilter_with_tags_basic() string {
	return `
data "huaweicloud_ram_resource_instances_filter" "test" {
  tags {
    key    = "test"
    values = ["value"]
  }
}
`
}

func testDataSourceResourceInstancesFilter_with_matches_basic() string {
	return `
data "huaweicloud_ram_resource_instances_filter" "test" {
  matches {
    key   = "test"
    value = "value"
  }
}
`
}

package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceInstancesCount_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_resource_instances_count.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceInstancesCount_without_any_tag_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
			{
				Config: testDataSourceResourceInstancesCount_with_tags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
			{
				Config: testDataSourceResourceInstancesCount_with_matches_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
		},
	})
}

func testDataSourceResourceInstancesCount_without_any_tag_basic() string {
	return `
data "huaweicloud_ram_resource_instances_count" "test" {
  without_any_tag = true
}
`
}

func testDataSourceResourceInstancesCount_with_tags_basic() string {
	return `
data "huaweicloud_ram_resource_instances_count" "test" {
  tags {
    key    = "abc"
    values = ["abc"]
  }
}
`
}

func testDataSourceResourceInstancesCount_with_matches_basic() string {
	return `
data "huaweicloud_ram_resource_instances_count" "test" {
  matches {
    key   = "test_key"
    value = "test_value"
  }
}
`
}

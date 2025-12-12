package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRAMResourceInstanceFilter_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ram_resource_instances_filter.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRAMResourceInstanceFilter_without_any_tag_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
			{
				Config: testDataSourceRAMResourceInstanceFilter_with_tag_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
		},
	})
}

func testDataSourceRAMResourceInstanceFilter_without_any_tag_basic() string {
	return `
data "huaweicloud_ram_resource_instances_filter" "test" {
	without_any_tag = true
}
`
}

func testDataSourceRAMResourceInstanceFilter_with_tag_basic() string {
	return `
data "huaweicloud_ram_resource_instances_filter" "test" {
	tags {
		key   = "abc"
		values = ["abc"]
	}
}
`
}

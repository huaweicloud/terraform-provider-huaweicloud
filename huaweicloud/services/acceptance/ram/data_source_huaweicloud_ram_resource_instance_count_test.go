package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRAMResourceInstanceCount_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ram_resource_instances_count.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRAMResourceInstanceCount_without_any_tag_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
			{
				Config: testDataSourceRAMResourceInstanceCount_with_tag_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
				),
			},
		},
	})
}

func testDataSourceRAMResourceInstanceCount_without_any_tag_basic() string {
	return `
data "huaweicloud_ram_resource_instances_count" "test" {
	without_any_tag = true
}
`
}

func testDataSourceRAMResourceInstanceCount_with_tag_basic() string {
	return `
data "huaweicloud_ram_resource_instances_count" "test" {
	tags {
		key   = "abc"
		values = ["abc"]
	}
}
`
}

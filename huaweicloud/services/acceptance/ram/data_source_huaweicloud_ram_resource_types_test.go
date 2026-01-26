package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceTypes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_resource_types.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resource_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_types.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_types.0.resource_type"),
				),
			},
		},
	})
}

func testDataSourceResourceTypes_basic() string {
	return `
  data "huaweicloud_ram_resource_types" "test" {}
`
}

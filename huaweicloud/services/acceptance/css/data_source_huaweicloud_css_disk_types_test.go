package css

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDiskTypes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_disk_types.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDiskTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "disk_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_types.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "disk_types.0.volume_names.#"),
				),
			},
		},
	})
}

func testDataSourceDiskTypes_basic() string {
	return `
data "huaweicloud_css_disk_types" "test" {}
`
}

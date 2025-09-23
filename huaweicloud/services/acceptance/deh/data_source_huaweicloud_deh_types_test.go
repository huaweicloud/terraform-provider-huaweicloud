package deh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDehTypes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_deh_types.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDehTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_host_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_host_types.0.host_type"),
					resource.TestCheckResourceAttrSet(dataSource, "dedicated_host_types.0.host_type_name"),
				),
			},
		},
	})
}

func testDataSourceDehTypes_basic() string {
	return `
data "huaweicloud_deh_types" "test" {
  availability_zone = "cn-north-4a"
}
`
}

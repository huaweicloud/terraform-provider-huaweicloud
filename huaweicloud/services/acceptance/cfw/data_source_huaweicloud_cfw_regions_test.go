package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwRegions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_regions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwRegions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data"),
				),
			},
		},
	})
}

func testDataSourceCfwRegions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_regions" "test" {
  fw_instance_id = "%[1]s"
}
`, acceptance.HW_CFW_INSTANCE_ID)
}

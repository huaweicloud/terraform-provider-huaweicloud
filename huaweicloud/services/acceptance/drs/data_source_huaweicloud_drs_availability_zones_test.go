package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAZs_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAZs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "names.#"),
				),
			},
		},
	})
}

const testAccDataSourceAZs_basic string = `data "huaweicloud_drs_availability_zones" "test" {
  engine_type = "mysql"
  type        = "migration"
  direction   = "up"
  node_type   = "high"
}`

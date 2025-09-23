package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNodeTypes_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_node_types.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNodeTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "node_types.#"),
				),
			},
		},
	})
}

const testAccDataSourceNodeTypes_basic string = `data "huaweicloud_drs_node_types" "test" {
  engine_type = "mysql"
  type        = "migration"
  direction   = "up"
  multi_write = false
}`

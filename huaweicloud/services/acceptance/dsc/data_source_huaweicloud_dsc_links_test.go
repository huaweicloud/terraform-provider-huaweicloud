package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLinks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_links.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLinks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "links.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.#"),
				),
			},
		},
	})
}

const testDataSourceLinks_basic = `
data "huaweicloud_dsc_links" "test" {
  db_name         = "db_name_test"
  ecs_name        = "ecs_name_test"
  labels          = ["key_test"]
  sensitive_level = 0
}
`

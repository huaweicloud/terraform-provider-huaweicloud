package cnad

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceAdvancedBlockList_basic(t *testing.T) {
	rName := "data.huaweicloud_cnad_advanced_block_list.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAdvancedBlockList_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "blocking_list.#"),
				),
			},
		},
	})
}

const testAccDatasourceAdvancedBlockList_basic = `
data "huaweicloud_cnad_advanced_block_list" "test" {
}
`

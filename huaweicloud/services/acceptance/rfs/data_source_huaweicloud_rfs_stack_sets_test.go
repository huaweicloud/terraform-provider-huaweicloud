package rfs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsStackSets_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_stack_sets.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRfsEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRfsStackSets_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "stack_sets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_sets.0.stack_set_id"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_sets.0.stack_set_name"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_sets.0.permission_model"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_sets.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_sets.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_sets.0.update_time"),
				),
			},
		},
	})
}

const testAccDataSourceRfsStackSets_basic = `
data "huaweicloud_rfs_stack_sets" "test" {}
`

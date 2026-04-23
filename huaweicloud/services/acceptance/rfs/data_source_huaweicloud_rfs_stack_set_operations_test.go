package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceStackSetOperations_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_stack_set_operations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the stack set in the environment before running this test case.
			acceptance.TestAccPreCheckRfsStackSetName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceStackSetOperations_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "stack_set_operations.#"),
				),
			},
		},
	})
}

func testDataSourceStackSetOperations_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rfs_stack_set_operations" "test" {
  stack_set_name = "%[1]s"
  sort_key       = "create_time"
  sort_dir       = "desc"
}
`, acceptance.HW_RFS_STACK_SET_NAME)
}

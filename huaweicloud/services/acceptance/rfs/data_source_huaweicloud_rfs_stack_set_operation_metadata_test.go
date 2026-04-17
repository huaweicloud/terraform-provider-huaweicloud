package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsStackSetOperationMetadata_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rfs_stack_set_operation_metadata.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the stack set operation data in the environment before running this test case.
			acceptance.TestAccPreCheckRfsStackSetName(t)
			acceptance.TestAccPreCheckRfsStackSetOperationId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRfsStackSetOperationMetadata_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "action"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "deployment_targets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "operation_preferences.#"),
				),
			},
		},
	})
}

func testDataSourceRfsStackSetOperationMetadata_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rfs_stack_set_operation_metadata" "test" {
  stack_set_name         = "%[1]s"
  stack_set_operation_id = "%[2]s"
}
`, acceptance.HW_RFS_STACK_SET_NAME, acceptance.HW_RFS_STACK_SET_OPERATION_ID)
}

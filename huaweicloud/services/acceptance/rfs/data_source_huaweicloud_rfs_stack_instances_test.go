package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsStackInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rfs_stack_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the stack instance data in the environment before running this test case.
			acceptance.TestAccPreCheckRfsStackSetName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRfsStackInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.stack_set_id"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.stack_set_name"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.latest_stack_set_operation_id"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.stack_domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "stack_instances.0.update_time"),

					resource.TestCheckOutput("is_stack_set_id_filter_result", "true"),
					resource.TestCheckOutput("is_stack_domain_id_filter_result", "true"),
				),
			},
		},
	})
}

func testDataSourceRfsStackInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rfs_stack_instances" "test" {
  stack_set_name = "%[1]s"
  sort_key       = "create_time"
  sort_dir       = "asc"
}

locals {
  stack_set_id    = data.huaweicloud_rfs_stack_instances.test.stack_instances.0.stack_set_id
  stack_domain_id = data.huaweicloud_rfs_stack_instances.test.stack_instances.0.stack_domain_id
}

# Filter by stack_set_id
data "huaweicloud_rfs_stack_instances" "stack_set_id_filter" {
  stack_set_name = "%[1]s"
  stack_set_id   = local.stack_set_id
}

locals {
  stack_set_id_filter_result = [
    for v in data.huaweicloud_rfs_stack_instances.stack_set_id_filter.stack_instances[*].stack_set_id : v == local.stack_set_id
  ]
}

output "is_stack_set_id_filter_result" {
  value = length(local.stack_set_id_filter_result) > 0 && alltrue(local.stack_set_id_filter_result)
}

# Filter by stack_domain_id
data "huaweicloud_rfs_stack_instances" "stack_domain_id_filter" {
  stack_set_name = "%[1]s"
  filter         = "stack_domain_id==${local.stack_domain_id}"
}

locals {
  stack_domain_id_filter_result = [
    for v in data.huaweicloud_rfs_stack_instances.stack_domain_id_filter.stack_instances[*].stack_domain_id : v == local.stack_domain_id
  ]
}

output "is_stack_domain_id_filter_result" {
  value = length(local.stack_domain_id_filter_result) > 0 && alltrue(local.stack_domain_id_filter_result)
}
`, acceptance.HW_RFS_STACK_SET_NAME)
}

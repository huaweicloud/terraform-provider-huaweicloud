package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseReplicationPolicyExecutionRecordTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_replication_policy_execution_record_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseReplicationPolicyExecutionRecordTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.execution_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.operation"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.src_resource"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.dst_resource"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status_revision"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.end_time"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseReplicationPolicyExecutionRecordTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_replication_policy_execute" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  policy_id   = huaweicloud_swr_enterprise_replication_policy.test.policy_id

  provisioner "local-exec" {
    command = "sleep 60"
  }
}

data "huaweicloud_swr_enterprise_replication_policy_execution_record_tasks" "test" {
  depends_on = [huaweicloud_swr_enterprise_replication_policy_execute.test]

  instance_id  = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  execution_id = huaweicloud_swr_enterprise_replication_policy_execute.test.execution_id
}
`, testAccSwrEnterpriseReplicationPolicy_basic(name))
}

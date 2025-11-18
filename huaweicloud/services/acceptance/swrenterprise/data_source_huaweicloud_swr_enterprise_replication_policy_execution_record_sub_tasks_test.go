package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseReplicationPolicyExecutionRecordSubTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_replication_policy_execution_record_sub_tasks.test"
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
				Config: testDataSourceSwrEnterpriseReplicationPolicyExecutionRecordSubTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.tag"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.status"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseReplicationPolicyExecutionRecordSubTasks_basic(name string) string {
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

data "huaweicloud_swr_enterprise_replication_policy_execution_record_sub_tasks" "test" {
  instance_id  = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  execution_id = huaweicloud_swr_enterprise_replication_policy_execute.test.execution_id
  task_id      = data.huaweicloud_swr_enterprise_replication_policy_execution_record_tasks.test.tasks[0].id
}
`, testAccSwrEnterpriseReplicationPolicy_basic(name))
}

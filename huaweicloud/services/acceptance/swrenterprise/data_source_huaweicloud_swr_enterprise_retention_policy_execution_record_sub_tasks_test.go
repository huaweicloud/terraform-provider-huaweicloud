package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseRetentionPolicyExecutionRecordSubTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseRetentionPolicyExecutionRecordSubTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.tag"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.op_time"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.status"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseRetentionPolicyExecutionRecordSubTasks_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_retention_policies" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_retention_policy_execution_records" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = data.huaweicloud_swr_enterprise_retention_policies.test.policies[0].id
}

data "huaweicloud_swr_enterprise_retention_policy_execution_record_tasks" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = data.huaweicloud_swr_enterprise_retention_policies.test.policies[0].id
  execution_id   = data.huaweicloud_swr_enterprise_retention_policy_execution_records.test.executions[0].id
}

data "huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = data.huaweicloud_swr_enterprise_retention_policies.test.policies[0].id
  execution_id   = data.huaweicloud_swr_enterprise_retention_policy_execution_records.test.executions[0].id
  task_id        = data.huaweicloud_swr_enterprise_retention_policy_execution_record_tasks.test.tasks[0].id
  status         = ""
}

data "huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks" "filter_by_status" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = data.huaweicloud_swr_enterprise_retention_policies.test.policies[0].id
  execution_id   = data.huaweicloud_swr_enterprise_retention_policy_execution_records.test.executions[0].id
  task_id        = data.huaweicloud_swr_enterprise_retention_policy_execution_record_tasks.test.tasks[0].id
  status         = data.huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks.test.sub_tasks[0].status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks.filter_by_status.sub_tasks) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks.filter_by_status.sub_tasks[*].status : 
	  v == data.huaweicloud_swr_enterprise_retention_policy_execution_record_sub_tasks.test.sub_tasks[0].status]
  )
}
`
}

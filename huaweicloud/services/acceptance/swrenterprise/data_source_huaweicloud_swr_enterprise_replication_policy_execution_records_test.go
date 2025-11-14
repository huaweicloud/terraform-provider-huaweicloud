package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseReplicationPolicyExecutionRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_replication_policy_execution_records.test"
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
				Config: testDataSourceSwrEnterpriseReplicationPolicyExecutionRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "executions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.trigger"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.total"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.stopped"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.succeed"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.failed"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.in_progress"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.status"),

					resource.TestCheckOutput("policy_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseReplicationPolicyExecutionRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_replication_policy_execution_records" "test" {
  depends_on = [huaweicloud_swr_enterprise_replication_policy_execute.test]

  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_replication_policy_execution_records" "filter_by_policy_id" {
  depends_on = [huaweicloud_swr_enterprise_replication_policy_execute.test]

  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  policy_id   = huaweicloud_swr_enterprise_replication_policy.test.policy_id
}

output "policy_id_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_replication_policy_execution_records.filter_by_policy_id.executions) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_replication_policy_execution_records.filter_by_policy_id.executions[*].policy_id :
	  v == data.huaweicloud_swr_enterprise_replication_policy_execution_records.test.executions[0].policy_id]
  )
}
`, testAccSwrEnterpriseReplicationPolicyExecute_basic(name))
}

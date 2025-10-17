package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseRetentionPolicyExecutionRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_retention_policy_execution_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseRetentionPolicyExecutionRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "executions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.trigger"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.dry_run"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.end_time"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseRetentionPolicyExecutionRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_retention_policy_execution_records" "test" {
  depends_on = [huaweicloud_swr_enterprise_retention_policy_execute.test]

  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_retention_policy.test.policy_id
}
`, testAccSwrEnterpriseRetentionPolicyExecute_basic(name))
}

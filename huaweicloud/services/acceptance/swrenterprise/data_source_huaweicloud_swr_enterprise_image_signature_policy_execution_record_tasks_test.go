package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecordTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks.test"
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
				Config: testDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecordTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.execution_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecordTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_image_signature_policy.test.policy_id
  execution_id   = data.huaweicloud_swr_enterprise_image_signature_policy_execution_records.test.executions[0].id
}
`, testDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecords_basic(name))
}

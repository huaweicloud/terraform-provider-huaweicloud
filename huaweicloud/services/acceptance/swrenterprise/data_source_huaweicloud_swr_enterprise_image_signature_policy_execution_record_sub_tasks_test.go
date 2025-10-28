package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecordSubTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_image_signature_policy_execution_record_sub_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecordSubTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.repository"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.digest"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.tags"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "sub_tasks.0.status"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecordSubTasks_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_image_signature_policies" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
}

data "huaweicloud_swr_enterprise_image_signature_policy_execution_records" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = data.huaweicloud_swr_enterprise_image_signature_policies.test.policies[0].id
}

data "huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = data.huaweicloud_swr_enterprise_image_signature_policies.test.policies[0].id
  execution_id   = data.huaweicloud_swr_enterprise_image_signature_policy_execution_records.test.executions[0].id
}

data "huaweicloud_swr_enterprise_image_signature_policy_execution_record_sub_tasks" "test" {
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  namespace_name = "library"
  policy_id      = data.huaweicloud_swr_enterprise_image_signature_policies.test.policies[0].id
  execution_id   = data.huaweicloud_swr_enterprise_image_signature_policy_execution_records.test.executions[0].id
  task_id        = data.huaweicloud_swr_enterprise_image_signature_policy_execution_record_tasks.test.tasks[0].id
}
`
}

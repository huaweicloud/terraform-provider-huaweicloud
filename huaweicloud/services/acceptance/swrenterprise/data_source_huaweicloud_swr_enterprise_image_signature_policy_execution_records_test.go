package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_image_signature_policy_execution_records.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecords_basic(rName),
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
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "executions.0.status_text"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseImageSignaturePolicyExecutionRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_image_signature_policy_execution_records" "test" {
  depends_on = [huaweicloud_swr_enterprise_image_signature_policy_execute.test]

  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"
  policy_id      = huaweicloud_swr_enterprise_image_signature_policy.test.id
}
`, testAccSwrEnterpriseImageSignaturePolicyExecute_basic(name))
}

package rms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRemediationExecutionStatuses_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_remediation_execution_statuses.test"
	basicConfig := testResourceRmsRemediationExecution_base()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRMSTargetIDForFGS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRemediationExecutionStatuses_forAllStatuses(basicConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policy_assignment_id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.resource_key.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.resource_key.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.resource_key.0.resource_provider"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.invocation_time"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.message"),
				),
			},
			{
				Config: testDataSourceRemediationExecutionStatuses_forSpecificStatuses(basicConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policy_assignment_id"),
					resource.TestCheckResourceAttrPair(dataSource, "value.0.resource_key.0.resource_id", "huaweicloud_vpc.test.0", "id"),
					resource.TestCheckResourceAttr(dataSource, "value.0.resource_key.0.resource_type", "vpcs"),
					resource.TestCheckResourceAttr(dataSource, "value.0.resource_key.0.resource_provider", "vpc"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.invocation_time"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "value.0.message"),
				),
			},
		},
	})
}

func testDataSourceRemediationExecutionStatuses_forAllStatuses(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_remediation_execution_statuses" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id

  depends_on = [huaweicloud_rms_remediation_execution.test]
}
`, testResourceRemediationExecution_basic(baseConfig))
}

func testDataSourceRemediationExecutionStatuses_forSpecificStatuses(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rms_remediation_execution_statuses" "test" {
  policy_assignment_id = huaweicloud_rms_policy_assignment.test.id
  
  resource_keys {
    resource_type     = "vpcs"
    resource_id       = huaweicloud_vpc.test[0].id
    resource_provider = "vpc"
  }

  depends_on = [huaweicloud_rms_remediation_execution.test]
}  
`, testResourceRemediationExecution_basic(baseConfig))
}

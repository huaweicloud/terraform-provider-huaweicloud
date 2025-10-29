package swrenterprise

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseInstanceAuditLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_instance_audit_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSwrEnterpriseInstanceAuditLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.operation"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.resource"),
					resource.TestCheckResourceAttrSet(dataSource, "audit_logs.0.op_time"),
					resource.TestCheckResourceAttrSet(dataSource, "total"),
				),
			},
		},
	})
}

func testDataSourceSwrEnterpriseInstanceAuditLogs_basic() string {
	return `
data "huaweicloud_swr_enterprise_instances" "test" {}

data "huaweicloud_swr_enterprise_instance_audit_logs" "test" {
  instance_id = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  operation   = "create"
}
`
}

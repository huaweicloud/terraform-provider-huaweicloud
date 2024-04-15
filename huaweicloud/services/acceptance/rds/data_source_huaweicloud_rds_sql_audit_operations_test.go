package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSqlAuditTypes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_sql_audit_operations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsSqlAuditOperations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "operations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "operations.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "operations.0.actions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "operations.0.actions.0"),

					resource.TestCheckOutput("operation_types_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsSqlAuditOperations_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_sql_audit_operations" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}

data "huaweicloud_rds_sql_audit_operations" "operation_types_filter" {
  instance_id     = huaweicloud_rds_instance.test.id
  operation_types = ["DDL", "DML"]
}
output "operation_types_filter_is_useful" {
  value = length(data.huaweicloud_rds_sql_audit_operations.operation_types_filter.operations) > 0 && alltrue(
  [for v in data.huaweicloud_rds_sql_audit_operations.operation_types_filter.operations[*].type : contains(["DDL", "DML"], v)]
  )
}
`, testAccRdsInstance_mysql_step1(name))
}

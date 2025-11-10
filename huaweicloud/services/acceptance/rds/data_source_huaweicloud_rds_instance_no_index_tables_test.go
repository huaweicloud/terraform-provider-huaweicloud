package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsInstanceNoIndexTables_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_instance_no_index_tables.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsInstanceNoIndexTables_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tables.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tables.0.table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "last_diagnose_timestamp"),
				),
			},
		},
	})
}

func testDataSourceRdsInstanceNoIndexTables_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_instance_no_index_tables" "test" {
  instance_id = "%s"
  table_type  = "no_primary_key"
  newest      = true
}
`, acceptance.HW_RDS_INSTANCE_ID)
}

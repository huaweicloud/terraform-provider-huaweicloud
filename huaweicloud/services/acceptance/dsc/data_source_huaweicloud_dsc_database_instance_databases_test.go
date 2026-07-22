package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDatabaseInstanceDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_database_instance_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscDbInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDatabaseInstanceDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.db_name"),
				),
			},
		},
	})
}

func testAccDataSourceDscDatabaseInstanceDatabases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_database_instance_databases" "test" {
  instance_id = "%s"
  type        = "MySQL"
}
`, acceptance.HW_DSC_DB_INSTANCE_ID)
}

package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBTableRestoredTables_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_table_restored_tables.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBTableRestoredTables_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_names.#"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBTableRestoredTables_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_table_restored_tables" "test" {
  instance_id   = "%s"
  database_name = "test_db"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}

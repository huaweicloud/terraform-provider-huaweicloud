package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBTableRestoredDatabases_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_table_restored_databases.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBTableRestoredDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "database_names.#"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBTableRestoredDatabases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_table_restored_databases" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}

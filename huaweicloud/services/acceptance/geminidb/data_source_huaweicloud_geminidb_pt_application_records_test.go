package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBPtApplicationRecords_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_pt_appliction_records.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbConfigID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBPtApplicationRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.applied_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "histories.0.apply_result"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBPtApplicationRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_pt_appliction_records" "test" {
  config_id = "%[1]s"
}
`, acceptance.HW_GEMINIDB_CONFIG_ID)
}

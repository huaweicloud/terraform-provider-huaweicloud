package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBigdataInstanceDatabases_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_bigdata_instance_databases.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDscInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBigdataInstanceDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "datasources.#"),
				),
			},
		},
	})
}

func testDataSourceBigdataInstanceDatabases_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_bigdata_instance_databases" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_DSC_INSTANCE_ID)
}

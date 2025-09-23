package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussErrorLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_error_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussInstanceId(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussErrorLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "log_files.#"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussErrorLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_error_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, acceptance.HW_GAUSSDB_OPENGAUSS_START_TIME, acceptance.HW_GAUSSDB_OPENGAUSS_END_TIME)
}

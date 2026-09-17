package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsSqlExportStatus_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_sql_export_status.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsSqlExportStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "export_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "progress_percentage"),
				),
			},
		},
	})
}

func testAccDataSourceDrsSqlExportStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_sql_export_status" "test" {
  job_id    = "%[1]s"
  file_type = "abnormal_sql"
}
`, acceptance.HW_DRS_JOB_ID)
}

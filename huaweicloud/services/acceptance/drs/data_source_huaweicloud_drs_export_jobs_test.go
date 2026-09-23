package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceExportJobs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_export_jobs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceExportJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "job_type", "backupMigration"),
					resource.TestCheckResourceAttrSet(dataSourceName, "async_job_id"),
				),
			},
		},
	})
}

// The export API is an asynchronous operation that returns an async job ID.
func testDataSourceExportJobs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_export_jobs" "test" {
  job_type = "backupMigration"
}
`)
}

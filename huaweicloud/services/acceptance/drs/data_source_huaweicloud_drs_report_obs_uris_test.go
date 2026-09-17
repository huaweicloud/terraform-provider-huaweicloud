package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReportObsUris_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_report_obs_uris.test"
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
				Config: testAccDataSourceReportObsUris_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "report_type", "abnormal_sql"),
				),
			},
		},
	})
}

func testAccDataSourceReportObsUris_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_report_obs_uris" "test" {
  job_id      = "%[1]s"
  report_type = "abnormal_sql"
}
`, acceptance.HW_DRS_JOB_ID)
}

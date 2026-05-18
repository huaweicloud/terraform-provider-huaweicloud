package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsBatchProgresses_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_batch_progresses.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsBatchProgresses_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.progress"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.task_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.transfer_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.remaining_time"),
				),
			},
		},
	})
}

func testAccDataSourceDrsBatchProgresses_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_batch_progresses" "test" {
  job_ids = split(",", "%s")
}
`, acceptance.HW_DRS_JOB_IDS)
}

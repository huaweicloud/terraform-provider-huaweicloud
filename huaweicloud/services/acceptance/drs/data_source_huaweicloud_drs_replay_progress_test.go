package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsReplayProgress_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_replay_progress.test"
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
				Config: testAccDataSourceDrsReplayProgress_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "progress"),
					resource.TestCheckResourceAttrSet(dataSourceName, "task_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "transfer_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "process_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replay_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "min_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "max_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "now_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "min_export_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "max_export_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "parse_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replay_sql_now_list.#"),
				),
			},
		},
	})
}

func testAccDataSourceDrsReplayProgress_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_replay_progress" "test" {
  job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID)
}

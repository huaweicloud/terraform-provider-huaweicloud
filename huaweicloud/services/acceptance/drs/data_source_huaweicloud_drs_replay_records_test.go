package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceReplayRecords_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_replay_records.test"
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
				Config: testAccDataSourceReplayRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sql_records.#"),
				),
			},
		},
	})
}

func testAccDataSourceReplayRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_replay_records" "test" {
  job_id     = "%[1]s"
  start_time = 1700421943
  end_time   = 1700439961
}
`, acceptance.HW_DRS_JOB_ID)
}

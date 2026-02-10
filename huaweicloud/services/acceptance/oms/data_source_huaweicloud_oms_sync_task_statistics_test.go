package oms

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSyncTaskStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_oms_sync_task_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		startTime  = time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime    = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOmsSyncTaskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSyncTaskStatistics_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "statistic_datas.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestCheckResourceAttrSet(dataSource, "statistic_datas.0.data_type"),
					resource.TestCheckResourceAttrSet(dataSource, "task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "statistic_time_type"),
				),
			},
		},
	})
}

func testAccDataSourceSyncTaskStatistics_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_oms_sync_task_statistics" "test" {
  sync_task_id = "%[1]s"
  data_type    = "SUCCESS"
  start_time   = %[2]d
  end_time     = %[3]d
}
`, acceptance.HW_OMS_SYNC_TASK_ID, startTime, endTime)
}

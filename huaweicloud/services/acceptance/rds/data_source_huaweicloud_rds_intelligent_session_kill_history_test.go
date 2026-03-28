package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIntelligentSessionKillHistory_basic(t *testing.T) {
	rName := "data.huaweicloud_rds_intelligent_session_kill_history.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIntelligentSessionKillHistory_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "history.#"),
					resource.TestCheckResourceAttrSet(rName, "history.0.task_id"),
					resource.TestCheckResourceAttrSet(rName, "history.0.start_time"),
					resource.TestCheckResourceAttrSet(rName, "history.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "history.0.download_link"),
				),
			},
		},
	})
}

func testAccDatasourceIntelligentSessionKillHistory_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_intelligent_session_kill_history" "test" {
  instance_id = "%s"
  start_time  = "2026-03-25 15:04:05"
  end_time    = "2026-03-28 15:04:05"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}

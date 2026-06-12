package das

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSlowLogExportTask_basic(t *testing.T) {
	var (
		rName = "huaweicloud_das_slow_log_export_task.test"

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSlowLogExportTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "bucket_name"),
				),
			},
		},
	})
}

func testAccSlowLogExportTask_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

locals {
  instance_ids = split(",", "%[2]s")
}
`, name, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccSlowLogExportTask_basic(name string) string {
	// 1.The earliest `start_time` is at most `2` days earlier than the current time.
	// 2.The latest `end_time` is at most `1` day later than the current time.
	// 3.The `end_time` must be greater than the `start_time`.
	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, -1).Format(time.RFC3339)
	endTime := currentTime.AddDate(0, 0, 1).Format(time.RFC3339)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_das_slow_log_export_task" "test" {
  instance_id = local.instance_ids[0]
  bucket_name = huaweicloud_obs_bucket.test.bucket
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  time_zone   = "Asia/Shanghai"
  export_type = "slowsqldetails"
}
`, testAccSlowLogExportTask_base(name), startTime, endTime)
}

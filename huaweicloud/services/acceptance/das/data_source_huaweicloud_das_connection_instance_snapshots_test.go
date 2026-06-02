package das

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccConnectionInstanceSnapshots_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_connection_instance_snapshots.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionInstanceSnapshots_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "snapshots.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.id"),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.status"),
					resource.TestCheckResourceAttrSet(all, "snapshots.0.find_lock"),
					resource.TestMatchResourceAttr(all, "snapshots.0.created_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccConnectionInstanceSnapshots_base() string {
	// 1.The earliest `start_time` and `end_time` is `7` days earlier than the current time.
	// 2.The `end_time` must be greater than the `start_time`.
	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, -3).Format(time.RFC3339)
	endTime := currentTime.AddDate(0, 0, 3).Format(time.RFC3339)

	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}

data "huaweicloud_das_database_users" "test" {
  instance_id = local.instance_ids[0]
}

locals {
  user_id    = try(data.huaweicloud_das_database_users.test.users.0.id, "")
  start_time = "%[2]s"
  end_time   = "%[3]s"
}
`, acceptance.HW_DAS_INSTANCE_IDS, startTime, endTime)
}

func testAccConnectionInstanceSnapshots_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_connection_instance_snapshots" "all" {
  user_id    = local.user_id
  module     = 0
  start_time = local.start_time
  end_time   = local.end_time
}
`, testAccConnectionInstanceSnapshots_base())
}

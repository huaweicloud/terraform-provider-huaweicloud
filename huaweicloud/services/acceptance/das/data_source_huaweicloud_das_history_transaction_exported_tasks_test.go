package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataHistoryTransactionExportedTasks_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_history_transaction_exported_tasks.all"
		dc  = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataHistoryTransactionExportedTasks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tasks.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.export_line_num"),
					resource.TestMatchResourceAttr(all, "tasks.0.start_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "tasks.0.end_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "tasks.0.created_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataHistoryTransactionExportedTasks_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

locals {
  instance_ids = split(",", "%[2]s")
}

resource "huaweicloud_das_history_transaction_export_task" "test" {
  instance_id = local.instance_ids[0]
  bucket_name = huaweicloud_obs_bucket.test.bucket
  start_time  = "2000-06-01T00:00:00+08:00"
  end_time    = "2099-06-02T00:00:00+08:00"
  time_zone   = "GMT+8"
}
`, name, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccDataHistoryTransactionExportedTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_history_transaction_exported_tasks" "all" {
  instance_id = local.instance_ids[0]

  depends_on = [
    huaweicloud_das_history_transaction_export_task.test
  ]
}
`, testAccDataHistoryTransactionExportedTasks_base(name))
}

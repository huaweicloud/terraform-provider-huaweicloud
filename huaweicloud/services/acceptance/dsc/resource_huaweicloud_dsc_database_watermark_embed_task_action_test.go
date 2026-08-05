package dsc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatabaseWatermarkEmbedTaskAction_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_dsc_database_watermark_embed_task_action.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscDbInstanceId(t)
			acceptance.TestAccPreCheckDscDbTableName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseWatermarkEmbedTaskAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "action", "DISABLE"),
				),
			},
			{
				Config: testAccDatabaseWatermarkEmbedTaskAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "action", "ENABLE"),
				),
			},
			{
				Config: testAccDatabaseWatermarkEmbedTaskAction_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "action", "START"),
				),
			},
		},
	})
}

func testAccDatabaseWatermarkEmbedTaskAction_base(name string) string {
	startTime := time.Now().Add(4 * time.Hour).Format("2006-01-02T15:04:05+08:00")
	return fmt.Sprintf(`
data "huaweicloud_dsc_authorize_databases" "test" {
  instance_type = "RDS"
}

locals {
  instance_id   = "%[1]s"
  instance      = try([for v in data.huaweicloud_dsc_authorize_databases.test.databases : v if v.ins_id == local.instance_id][0], {})
  instance_name = try(local.instance.ins_name, "NOT_FOUND")
  database_id   = try(local.instance.id, "NOT_FOUND")
  database_name = try(local.instance.db_name, "NOT_FOUND")
}

resource "huaweicloud_dsc_database_watermark_embed_task" "test" {
  task_name         = "%[2]s"
  water_mark        = "TF"
  watermark_version = "V2"

  db_water_param {
    embed_mode = "EMBED_FAKE_COLUMN"

    params {
      new_column_name = "test"
      new_column_type = "varchar"
      fake_strategy   = "name"
    }
  }

  source_db_info {
    db_id      = local.database_id
    db_name    = local.database_name
    db_type    = "MySQL"
    ins_id     = local.instance_id
    ins_name   = local.instance_name
    table_name = "%[4]s"
  }

  target_db_info {
    db_id      = local.database_id
    db_name    = local.database_name
    db_type    = "MySQL"
    ins_id     = local.instance_id
    ins_name   = local.instance_name
    table_name = "%[2]s"
  }

  error_code    = 1
  schedule_type = "MONTH"
  start_time    = "%[3]s"

  # Each time the task is restarted, 'start_time' will be changed to the time when the current task started execution.
  lifecycle {
    ignore_changes = [start_time]
  }
}
`, acceptance.HW_DSC_DB_INSTANCE_ID, name, startTime, acceptance.HW_DSC_DB_TABLE_NAME)
}

func testAccDatabaseWatermarkEmbedTaskAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_database_watermark_embed_task_action" "test" {
  task_id = huaweicloud_dsc_database_watermark_embed_task.test.id
  action  = "DISABLE"
}
`, testAccDatabaseWatermarkEmbedTaskAction_base(name))
}

func testAccDatabaseWatermarkEmbedTaskAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_database_watermark_embed_task_action" "test" {
  task_id = huaweicloud_dsc_database_watermark_embed_task.test.id
  action  = "ENABLE"

  enable_force_new = "true"
}
`, testAccDatabaseWatermarkEmbedTaskAction_base(name))
}

func testAccDatabaseWatermarkEmbedTaskAction_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_database_watermark_embed_task_action" "test" {
  task_id = huaweicloud_dsc_database_watermark_embed_task.test.id
  action  = "START"

  enable_force_new = "true"
}
`, testAccDatabaseWatermarkEmbedTaskAction_base(name))
}

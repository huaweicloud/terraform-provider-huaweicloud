package dsc

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, please ensure that the DSC instance has been created
// and at least one database watermark embed task exists.
func TestAccDataSourceDatabaseWatermarkEmbedTasks_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dsc_database_watermark_embed_tasks.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byTaskId   = "data.huaweicloud_dsc_database_watermark_embed_tasks.filter_by_task_id"
		dcByTaskId = acceptance.InitDataSourceCheck(byTaskId)

		byStatus   = "data.huaweicloud_dsc_database_watermark_embed_tasks.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byStartAndEnd   = "data.huaweicloud_dsc_database_watermark_embed_tasks.filter_by_start_and_end"
		dcByStartAndEnd = acceptance.InitDataSourceCheck(byStartAndEnd)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscDbInstanceId(t)
			acceptance.TestAccPreCheckDscDbTableName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseWatermarkEmbedTasks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tasks.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'task_id' parameter.
					dcByTaskId.CheckResourceExists(),
					resource.TestCheckOutput("is_task_id_filter_useful", "true"),
					resource.TestCheckResourceAttr(byTaskId, "tasks.#", "1"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.id", "huaweicloud_dsc_database_watermark_embed_task.test", "id"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.task_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "task_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.water_mark",
						"huaweicloud_dsc_database_watermark_embed_task.test", "water_mark"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.watermark_version",
						"huaweicloud_dsc_database_watermark_embed_task.test", "watermark_version"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.watermark_describe",
						"huaweicloud_dsc_database_watermark_embed_task.test", "watermark_describe"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.start_now",
						"huaweicloud_dsc_database_watermark_embed_task.test", "start_now"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.schedule_type",
						"huaweicloud_dsc_database_watermark_embed_task.test", "schedule_type"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.schedule_switch",
						"huaweicloud_dsc_database_watermark_embed_task.test", "schedule_switch"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.task_state",
						"huaweicloud_dsc_database_watermark_embed_task.test", "task_state"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.start_time",
						"huaweicloud_dsc_database_watermark_embed_task.test", "start_time"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.embed_mode",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.embed_mode"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.#",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.#"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.0.new_column_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.0.new_column_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.0.new_column_type",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.0.new_column_type"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.0.fake_strategy",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.0.fake_strategy"),
					resource.TestCheckResourceAttrSet(byTaskId, "tasks.0.db_water_param.0.params.0.fake_param.0.date_begin"),
					resource.TestCheckResourceAttrSet(byTaskId, "tasks.0.db_water_param.0.params.0.fake_param.0.date_end"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.1.new_column_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.1.new_column_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.1.new_column_type",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.1.new_column_type"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.1.fake_strategy",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.1.fake_strategy"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.1.fake_param.0.random_begin",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.1.fake_param.0.random_begin"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.1.fake_param.0.random_end",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.1.fake_param.0.random_end"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.db_water_param.0.params.1.fake_param.0.random_accuracy",
						"huaweicloud_dsc_database_watermark_embed_task.test", "db_water_param.0.params.1.fake_param.0.random_accuracy"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.source_db_info.0.db_id",
						"huaweicloud_dsc_database_watermark_embed_task.test", "source_db_info.0.db_id"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.source_db_info.0.db_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "source_db_info.0.db_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.source_db_info.0.db_type",
						"huaweicloud_dsc_database_watermark_embed_task.test", "source_db_info.0.db_type"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.source_db_info.0.ins_id",
						"huaweicloud_dsc_database_watermark_embed_task.test", "source_db_info.0.ins_id"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.source_db_info.0.ins_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "source_db_info.0.ins_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.source_db_info.0.table_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "source_db_info.0.table_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.target_db_info.0.db_id",
						"huaweicloud_dsc_database_watermark_embed_task.test", "target_db_info.0.db_id"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.target_db_info.0.db_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "target_db_info.0.db_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.target_db_info.0.db_type",
						"huaweicloud_dsc_database_watermark_embed_task.test", "target_db_info.0.db_type"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.target_db_info.0.ins_id",
						"huaweicloud_dsc_database_watermark_embed_task.test", "target_db_info.0.ins_id"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.target_db_info.0.ins_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "target_db_info.0.ins_name"),
					resource.TestCheckResourceAttrPair(byTaskId, "tasks.0.target_db_info.0.table_name",
						"huaweicloud_dsc_database_watermark_embed_task.test", "target_db_info.0.table_name"),
					resource.TestMatchResourceAttr(byTaskId, "tasks.0.task_create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byTaskId, "tasks.0.task_end_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Filter by 'status' parameter.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Filter by 'start' and 'end' parameters.
					dcByStartAndEnd.CheckResourceExists(),
					resource.TestCheckOutput("is_start_and_end_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDatabaseWatermarkEmbedTasks_base(name string) string {
	currentTime := time.Now()
	beginTime := currentTime.Format("2006-01-02T15:04:05+08:00")
	endTime := currentTime.Add(1 * time.Hour).Format("2006-01-02T15:04:05+08:00")
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
  task_name          = "%[2]s"
  water_mark         = "TF"
  watermark_version  = "V2"
  watermark_describe = "Created by Terraform"

  db_water_param {
    embed_mode = "EMBED_FAKE_COLUMN"

    params {
      new_column_name = "field1"
      new_column_type = "date"
      fake_strategy   = "date"

      fake_param {
        date_begin = "%[3]s"
        date_end   = "%[4]s"
      }
    }
    params {
      new_column_name = "field2"
      new_column_type = "number"
      fake_strategy   = "number_random"

      fake_param {
        random_begin    = "2"
        random_end      = "3"
        random_accuracy = 1
      }
    }
  }

  source_db_info {
    db_id      = local.database_id
    db_name    = local.database_name
    db_type    = "MySQL"
    ins_id     = local.instance_id
    ins_name   = local.instance_name
    table_name = "%[5]s"
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
  start_now     = true
  schedule_type = "ONCE"
}
`, acceptance.HW_DSC_DB_INSTANCE_ID, name, beginTime, endTime, acceptance.HW_DSC_DB_TABLE_NAME)
}

func testAccDataSourceDatabaseWatermarkEmbedTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_dsc_database_watermark_embed_tasks" "test" {
  depends_on = [huaweicloud_dsc_database_watermark_embed_task.test]
}

# Filter by 'task_id' parameter.
locals {
  task_id = huaweicloud_dsc_database_watermark_embed_task.test.id
}

data "huaweicloud_dsc_database_watermark_embed_tasks" "filter_by_task_id" {
  task_id = local.task_id
}

locals {
  task_id_filter_result = [
    for v in data.huaweicloud_dsc_database_watermark_embed_tasks.filter_by_task_id.tasks[*].id : v == local.task_id
  ]
}

output "is_task_id_filter_useful" {
  value = length(local.task_id_filter_result) == 1 && alltrue(local.task_id_filter_result)
}

# Filter by 'status' parameter.
locals {
  task_state = huaweicloud_dsc_database_watermark_embed_task.test.task_state
}

data "huaweicloud_dsc_database_watermark_embed_tasks" "filter_by_status" {
  status = local.task_state
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dsc_database_watermark_embed_tasks.filter_by_status.tasks[*].task_state : v == local.task_state
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by 'start' and 'end' parameters.
locals {
  start_time = huaweicloud_dsc_database_watermark_embed_task.test.start_time
  end_time   = timeadd(local.start_time, "1h")
}

data "huaweicloud_dsc_database_watermark_embed_tasks" "filter_by_start_and_end" {
  start = local.start_time
  end   = local.end_time
}

locals {
  start_and_end_filter_result = [
    for v in data.huaweicloud_dsc_database_watermark_embed_tasks.filter_by_start_and_end.tasks[*].start_time :
    timecmp(v, local.start_time) >= 0 && timecmp(v, local.end_time) <= 0
  ]
}

output "is_start_and_end_filter_useful" {
  value = length(local.start_and_end_filter_result) > 0 && alltrue(local.start_and_end_filter_result)
}
`, testAccDataSourceDatabaseWatermarkEmbedTasks_base(name))
}

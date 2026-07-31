package dsc

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dsc"
)

func getDatabaseWatermarkEmbedTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dsc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DSC client: %s", err)
	}

	return dsc.GetDatabaseWatermarkEmbedTaskById(client, state.Primary.ID)
}

// Before this test, please ensure that the DSC instance has been created and authorized RDS MySQL database assets.
// HW_DSC_DB_TABLE_NAME must belong to the currently provided authorized database.
func TestAccDatabaseWatermarkEmbedTask_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dsc_database_watermark_embed_task.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getDatabaseWatermarkEmbedTaskResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscDbInstanceId(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPreCheckDscDbTableName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseWatermarkEmbedTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", name),
					resource.TestCheckResourceAttr(rName, "water_mark", "TF"),
					resource.TestCheckResourceAttr(rName, "watermark_version", "V1"),
					resource.TestCheckResourceAttr(rName, "watermark_describe", "Created by Terraform"),
					resource.TestCheckResourceAttr(rName, "db_water_param.#", "1"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.embed_mode", "EMBED_FAKE_COLUMN"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.watermark_key", acceptance.HW_PROJECT_ID),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.#", "2"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.0.new_column_name", "field1"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.0.new_column_type", "date"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.0.fake_strategy", "date"),
					resource.TestCheckResourceAttrSet(rName, "db_water_param.0.params.0.fake_param.0.date_begin"),
					resource.TestCheckResourceAttrSet(rName, "db_water_param.0.params.0.fake_param.0.date_end"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.1.new_column_name", "field2"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.1.new_column_type", "number"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.1.fake_strategy", "number_random"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.1.fake_param.0.random_begin", "2"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.1.fake_param.0.random_end", "3"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.1.fake_param.0.random_accuracy", "1"),
					resource.TestCheckResourceAttr(rName, "source_db_info.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "source_db_info.0.db_id"),
					resource.TestCheckResourceAttrSet(rName, "source_db_info.0.db_name"),
					resource.TestCheckResourceAttr(rName, "source_db_info.0.db_type", "MySQL"),
					resource.TestCheckResourceAttr(rName, "source_db_info.0.ins_id", acceptance.HW_DSC_DB_INSTANCE_ID),
					resource.TestCheckResourceAttrSet(rName, "source_db_info.0.ins_name"),
					resource.TestCheckResourceAttr(rName, "source_db_info.0.table_name", acceptance.HW_DSC_DB_TABLE_NAME),
					resource.TestCheckResourceAttr(rName, "target_db_info.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "target_db_info.0.db_id"),
					resource.TestCheckResourceAttrSet(rName, "target_db_info.0.db_name"),
					resource.TestCheckResourceAttr(rName, "target_db_info.0.db_type", "MySQL"),
					resource.TestCheckResourceAttr(rName, "target_db_info.0.ins_id", acceptance.HW_DSC_DB_INSTANCE_ID),
					resource.TestCheckResourceAttrSet(rName, "target_db_info.0.ins_name"),
					resource.TestCheckResourceAttr(rName, "target_db_info.0.table_name", name),
					resource.TestCheckResourceAttr(rName, "error_code", "1"),
					resource.TestCheckResourceAttr(rName, "start_now", "true"),
					resource.TestCheckResourceAttr(rName, "schedule_type", "ONCE"),
					resource.TestCheckResourceAttr(rName, "schedule_switch", "true"),
					resource.TestMatchResourceAttr(rName, "start_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "task_state"),
					resource.TestCheckResourceAttrSet(rName, "task_create_time"),
					resource.TestMatchResourceAttr(rName, "task_create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "task_end_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccDatabaseWatermarkEmbedTask_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "task_name", updateName),
					resource.TestCheckResourceAttr(rName, "water_mark", "TFUpdate"),
					resource.TestCheckResourceAttr(rName, "watermark_version", "V2"),
					resource.TestCheckResourceAttr(rName, "watermark_describe", "update description"),
					resource.TestCheckResourceAttr(rName, "db_water_param.#", "1"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.embed_mode", "EMBED_FAKE_COLUMN"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.#", "1"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.0.new_column_name", "test_update"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.0.new_column_type", "varchar"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.0.fake_strategy", "name"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.0.fake_param.#", "0"),
					resource.TestCheckResourceAttr(rName, "target_db_info.0.table_name", updateName),
					resource.TestCheckResourceAttr(rName, "error_code", "2"),
					resource.TestCheckResourceAttr(rName, "schedule_switch", "true"),
					resource.TestCheckResourceAttr(rName, "schedule_type", "MONTH"),
				),
			},
			{
				Config: testAccDatabaseWatermarkEmbedTask_basic_step3(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "watermark_describe", ""),
					resource.TestCheckResourceAttr(rName, "db_water_param.#", "1"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.embed_mode", "EMBED_FAKE_ROW"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.row_spacing", "2"),
					resource.TestCheckResourceAttr(rName, "db_water_param.0.params.#", "0"),
					resource.TestCheckResourceAttr(rName, "schedule_switch", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDatabaseWatermarkEmbedTask_base() string {
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
`, acceptance.HW_DSC_DB_INSTANCE_ID)
}

func testAccDatabaseWatermarkEmbedTask_basic(name string) string {
	currentTime := time.Now()
	beginTime := currentTime.Format("2006-01-02T15:04:05+08:00")
	endTime := currentTime.Add(1 * time.Hour).Format("2006-01-02T15:04:05+08:00")

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_database_watermark_embed_task" "test" {
  task_name          = "%[2]s"
  water_mark         = "TF"
  watermark_version  = "V1"
  watermark_describe = "Created by Terraform"

  db_water_param {
    embed_mode    = "EMBED_FAKE_COLUMN"
    watermark_key = "%[3]s"

    params {
      new_column_name = "field1"
      new_column_type = "date"
      fake_strategy   = "date"

      fake_param {
        date_begin = "%[4]s"
        date_end   = "%[5]s"
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
    table_name = "%[6]s"
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
`, testAccDatabaseWatermarkEmbedTask_base(), name, acceptance.HW_PROJECT_ID, beginTime, endTime, acceptance.HW_DSC_DB_TABLE_NAME)
}

func testAccDatabaseWatermarkEmbedTask_basic_step2(updateName string) string {
	startTime := time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05+08:00")
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_database_watermark_embed_task" "test" {
  task_name         = "%[2]s"
  water_mark        = "TFUpdate"
  watermark_version = "V2"

  db_water_param {
    embed_mode = "EMBED_FAKE_COLUMN"

    params {
      new_column_name = "test_update"
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

  error_code         = 2
  watermark_describe = "update description"
  schedule_switch    = true
  schedule_type      = "MONTH"
  start_time         = "%[3]s"
}
`, testAccDatabaseWatermarkEmbedTask_base(), updateName, startTime, acceptance.HW_DSC_DB_TABLE_NAME)
}

func testAccDatabaseWatermarkEmbedTask_basic_step3(updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dsc_database_watermark_embed_task" "test" {
  task_name         = "%[2]s"
  water_mark        = "TFUpdate"
  watermark_version = "V2"

  db_water_param {
    embed_mode  = "EMBED_FAKE_ROW"
    row_spacing = "2"
  }

  source_db_info {
    db_id      = local.database_id
    db_name    = local.database_name
    db_type    = "MySQL"
    ins_id     = local.instance_id
    ins_name   = local.instance_name
    table_name = "%[3]s"
  }

  target_db_info {
    db_id      = local.database_id
    db_name    = local.database_name
    db_type    = "MySQL"
    ins_id     = local.instance_id
    ins_name   = local.instance_name
    table_name = "%[2]s"
  }

  error_code      = 1
  schedule_switch = false
}
`, testAccDatabaseWatermarkEmbedTask_base(), updateName, acceptance.HW_DSC_DB_TABLE_NAME)
}

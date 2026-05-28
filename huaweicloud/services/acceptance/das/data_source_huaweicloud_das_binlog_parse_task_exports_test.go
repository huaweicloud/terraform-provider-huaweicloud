package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataBinlogParseTaskExports_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_binlog_parse_task_exports.all"
		dc  = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataBinlogParseTaskExports_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tasks.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "tasks.0.exported_task_id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.status"),
					resource.TestMatchResourceAttr(all, "tasks.0.start_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "tasks.0.end_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "tasks.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "tasks.0.export_line_num"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.source_file_name"),
					resource.TestCheckResourceAttrSet(all, "tasks.0.parsed_task_id"),
				),
			},
		},
	})
}

func testAccDataBinlogParseTaskExports_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

locals {
  instance_ids = split(",", "%[2]s")
}

data "huaweicloud_das_database_users" "test" {
  instance_id = local.instance_ids[0]
}

locals {
  user_id = try(data.huaweicloud_das_database_users.test.users.0.id, "")
}

data "huaweicloud_das_binlogs" "test" {
  user_id     = local.user_id
  binlog_type = "latest"
}

resource "huaweicloud_das_binlog_parse_task" "test" {
  user_id     = local.user_id
  binlog_type = "latest"
  file_name   = try(data.huaweicloud_das_binlogs.test.binlogs.0.file_name, "")
}

resource "huaweicloud_das_binlog_parse_task_export" "test" {
  user_id     = local.user_id
  task_id     = huaweicloud_das_binlog_parse_task.test.id
  bucket_name = huaweicloud_obs_bucket.test.bucket

  filter_condition {
    types               = ["insert", "update", "delete", "ddl"]
    start_time          = "2000-06-01T00:00:00+08:00"
    end_time            = "2099-06-02T00:00:00+08:00"
    parse_double_insert = true
  }
}
`, name, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccDataBinlogParseTaskExports_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# After exporting the parsing task, it takes some time before its value can be obtained.
# Although the POST interface for exporting the parsing task is a synchronous interface.
resource "time_sleep" "wait_10_second" {
  create_duration = "10s"
  depends_on      = [huaweicloud_das_binlog_parse_task_export.test]
}

data "huaweicloud_das_binlog_parse_task_exports" "all" {
  user_id = local.user_id

  depends_on = [
    huaweicloud_das_binlog_parse_task_export.test,
    time_sleep.wait_10_second
  ]
}
`, testAccDataBinlogParseTaskExports_base(name))
}

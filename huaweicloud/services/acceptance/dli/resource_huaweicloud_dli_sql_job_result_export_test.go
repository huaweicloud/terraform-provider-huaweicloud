package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSqlJobResultExport_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliSQLQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlJobResultExport_basic_step1(name),
			},
			{
				Config: testAccSqlJobResultExport_basic_step2(name),
			},
		},
	})
}

func testAccSqlJobResultExport_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = replace("%[1]s", "_", "-")
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_dli_database" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dli_table" "test" {
  database_name = huaweicloud_dli_database.test.name
  name          = "%[1]s"
  data_location = "DLI"

  columns {
    name = "field1"
    type = "string"
  }

  columns {
    name = "field2"
    type = "string"
  }
}

resource "huaweicloud_dli_sql_job" "test" {
  sql           = "SELECT * FROM ${huaweicloud_dli_table.test.name}"
  database_name = huaweicloud_dli_database.test.name
  queue_name    = "%[2]s"
}

# Wait for the SQL query job to complete.
resource "time_sleep" "wait_20_seconds" {
  depends_on = [huaweicloud_dli_sql_job.test]

  create_duration = "20s"
}
`, name, acceptance.HW_DLI_SQL_QUEUE_NAME)
}

func testAccSqlJobResultExport_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_sql_job_result_export" "test" {
  job_id    = huaweicloud_dli_sql_job.test.id
  data_path = "obs://${huaweicloud_obs_bucket.test.bucket}/%[2]s"
  data_type = "json"

  depends_on = [time_sleep.wait_20_seconds]
}
`, testAccSqlJobResultExport_base(name), name)
}

func testAccSqlJobResultExport_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dli_sql_job_result_export" "test" {
  job_id             = huaweicloud_dli_sql_job.test.id
  data_path          = "obs://${huaweicloud_obs_bucket.test.bucket}/%[2]s"
  data_type          = "csv"
  queue_name         = "%[3]s"
  export_mode        = "Overwrite"
  limit_num          = 10
  with_column_header = true
  quote_char         = "\""
  escape_char        = "\\"

  enable_force_new = "true"
}
`, testAccSqlJobResultExport_base(name), name, acceptance.HW_DLI_SQL_QUEUE_NAME)
}

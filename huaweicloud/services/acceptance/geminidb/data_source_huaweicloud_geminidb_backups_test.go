package geminidb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGeminiDBBackups_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_geminidb_backups.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGeminiDBBackups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "backups.0.database_tables.#"),

					resource.TestCheckOutput("is_instance_filter_useful", "true"),
					resource.TestCheckOutput("is_datastore_type_filter_useful", "true"),
					resource.TestCheckOutput("is_backup_id_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_backup_type_filter_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGeminiDBBackups_basic(rName string) string {
	beginTime := time.Now().UTC().Add(-5 * time.Minute)
	beginTimeString := beginTime.Format("2006-01-02T15:04:05+0000")
	endTime := time.Now().UTC().Add(3 * time.Minute)
	endTimeString := endTime.Format("2006-01-02T15:04:05+0000")
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_geminidb_backups" "test" {}

data "huaweicloud_geminidb_backups" "instance_filter" {
 instance_id = huaweicloud_geminidb_backup.test.instance_id
}

output "is_instance_filter_useful" {
  value = length(data.huaweicloud_geminidb_backups.instance_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_backups.instance_filter.backups[*].instance_id :
    v == huaweicloud_geminidb_backup.test.instance_id]
  )
}

data "huaweicloud_geminidb_backups" "datastore_type_filter" {
 datastore_type = "cassandra"

  depends_on = [huaweicloud_geminidb_backup.test]
}

output "is_datastore_type_filter_useful" {
  value = length(data.huaweicloud_geminidb_backups.datastore_type_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_backups.datastore_type_filter.backups[*].datastore[0].type :
    v == "cassandra"]
  )
}

data "huaweicloud_geminidb_backups" "backup_id_filter" {
  backup_id  = huaweicloud_geminidb_backup.test.id
  type       = "DatabaseTable"
  depends_on = [huaweicloud_geminidb_backup.test]
}

output "is_backup_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_backups.backup_id_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_backups.backup_id_filter.backups[*].id :
    v == huaweicloud_geminidb_backup.test.id]
  )
}

data "huaweicloud_geminidb_backups" "type_filter" {
  type       = "DatabaseTable"
  depends_on = [huaweicloud_geminidb_backup.test]
}

output "is_type_filter_useful" {
 value = length(data.huaweicloud_geminidb_backups.type_filter.backups) > 0
}

data "huaweicloud_geminidb_backups" "backup_type_filter" {
  backup_type = "Manual"
  type        = "DatabaseTable"
  depends_on  = [huaweicloud_geminidb_backup.test]
}

output "is_backup_type_filter_useful" {
  value = length(data.huaweicloud_geminidb_backups.backup_type_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_backups.backup_type_filter.backups[*].type :
    v == "Manual"]
  )
}

data "huaweicloud_geminidb_backups" "time_filter" {
  begin_time = "%[2]s"
  end_time   = "%[3]s"
  type       = "DatabaseTable"

  depends_on = [huaweicloud_geminidb_backup.test]
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_geminidb_backups.time_filter.backups) > 0
}
`, testAccGeminiDBBackup_withDatabaseTables(rName), beginTimeString, endTimeString)
}

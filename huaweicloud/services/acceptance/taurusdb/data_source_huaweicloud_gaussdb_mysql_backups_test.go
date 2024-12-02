package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlBackups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_backups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussdbMysqlBackups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.take_up_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),

					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("backup_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("backup_type_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_name_filter_is_useful", "true"),
					resource.TestCheckOutput("begin_time_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussdbMysqlBackups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_mysql_backups" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_backup.test]
}

locals {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}
data "huaweicloud_gaussdb_mysql_backups" "instance_id_filter" {
  depends_on  = [huaweicloud_gaussdb_mysql_backup.test]
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_backups.instance_id_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_backups.instance_id_filter.backups[*].instance_id : v == local.instance_id]
  )
}

locals {
  backup_id = huaweicloud_gaussdb_mysql_backup.test.id
}
data "huaweicloud_gaussdb_mysql_backups" "backup_id_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_backup.test]
  backup_id  = huaweicloud_gaussdb_mysql_backup.test.id
}
output "backup_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_backups.backup_id_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_backups.backup_id_filter.backups[*].id : v == local.backup_id]
  )
}

data "huaweicloud_gaussdb_mysql_backups" "name_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_backup.test]
  name       = huaweicloud_gaussdb_mysql_backup.test.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_backups.name_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_backups.name_filter.backups[*].name : length(regex("%[2]s", v))>0]
  )
}

locals {
  backup_type = "manual"
}
data "huaweicloud_gaussdb_mysql_backups" "backup_type_filter" {
  depends_on  = [huaweicloud_gaussdb_mysql_backup.test]
  backup_type = "manual"
}
output "backup_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_backups.backup_type_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_backups.backup_type_filter.backups[*].type : v == local.backup_type]
  )
}

data "huaweicloud_gaussdb_mysql_backups" "instance_name_filter" {
  depends_on    = [huaweicloud_gaussdb_mysql_backup.test]
  instance_name = huaweicloud_gaussdb_mysql_instance.test.name
}
output "instance_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_backups.instance_name_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_backups.instance_name_filter.backups[*].instance_name : v == "%[2]s"]
  )
}

locals {
  begin_time = huaweicloud_gaussdb_mysql_backup.test.begin_time
}
data "huaweicloud_gaussdb_mysql_backups" "begin_time_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_backup.test]
  begin_time = huaweicloud_gaussdb_mysql_backup.test.begin_time
}
output "begin_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_backups.begin_time_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_backups.begin_time_filter.backups[*].begin_time : v == local.begin_time]
  )
}

locals {
  status = "COMPLETED"
}
data "huaweicloud_gaussdb_mysql_backups" "status_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_backup.test]
  status     = "COMPLETED"
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_backups.status_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_backups.status_filter.backups[*].status : v == local.status]
  )
}
`, testBackup_basic(name), name)
}

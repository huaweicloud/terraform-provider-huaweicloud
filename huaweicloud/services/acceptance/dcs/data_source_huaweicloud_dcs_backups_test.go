package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDcsBackups_basic(t *testing.T) {
	rName := "data.huaweicloud_dcs_backups.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDcsBackups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "backups.0.error_code", ""),
					resource.TestCheckResourceAttrSet(rName, "backups.#"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.id"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.name"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.size"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.type"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.status"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.description"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.is_support_restore"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.backup_format"),
					resource.TestCheckResourceAttrSet(rName, "backups.0.progress"),
					resource.TestCheckOutput("backup_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("is_support_restore_filter_is_useful", "true"),
					resource.TestCheckOutput("backup_format_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDcsBackups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_backups" "test" {
  depends_on  = [huaweicloud_dcs_backup.test]
  instance_id = huaweicloud_dcs_instance.instance_1.id
}

data "huaweicloud_dcs_backups" "backup_id_filter" {
  depends_on  = [huaweicloud_dcs_backup.test]
  instance_id = huaweicloud_dcs_instance.instance_1.id
  backup_id   = local.id
}

locals {
  id = split("/",huaweicloud_dcs_backup.test.id).1
}
	
output "backup_id_filter_is_useful" {
  value = length(data.huaweicloud_dcs_backups.backup_id_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_backups.backup_id_filter.backups[*].id : v == local.id]
  )
}

data "huaweicloud_dcs_backups" "name_filter" {
  depends_on  = [huaweicloud_dcs_backup.test]
  instance_id = huaweicloud_dcs_instance.instance_1.id
  name        = huaweicloud_dcs_backup.test.name
}

locals {
  name = huaweicloud_dcs_backup.test.name
}
	
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dcs_backups.name_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_backups.name_filter.backups[*].name : v == local.name]
  )
}

data "huaweicloud_dcs_backups" "status_filter" {
  depends_on  = [huaweicloud_dcs_backup.test]
  instance_id = huaweicloud_dcs_instance.instance_1.id
  status      = huaweicloud_dcs_backup.test.status
}

locals {
  status = huaweicloud_dcs_backup.test.status
}
	
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dcs_backups.status_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_backups.status_filter.backups[*].status : v == local.status]
  )
}

data "huaweicloud_dcs_backups" "type_filter" {
  depends_on  = [huaweicloud_dcs_backup.test]
  instance_id = huaweicloud_dcs_instance.instance_1.id
  type        = huaweicloud_dcs_backup.test.type
}

locals {
  type = huaweicloud_dcs_backup.test.type
}
	
output "type_filter_is_useful" {
  value = length(data.huaweicloud_dcs_backups.type_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_backups.type_filter.backups[*].type : v == local.type]
  )
}

data "huaweicloud_dcs_backups" "is_support_restore_filter" {
  depends_on         = [huaweicloud_dcs_backup.test]
  instance_id        = huaweicloud_dcs_instance.instance_1.id
  is_support_restore = huaweicloud_dcs_backup.test.is_support_restore
}

locals {
  is_support_restore = huaweicloud_dcs_backup.test.is_support_restore
}
	
output "is_support_restore_filter_is_useful" {
  value = length(data.huaweicloud_dcs_backups.is_support_restore_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_backups.is_support_restore_filter.backups[*].is_support_restore : v == local.is_support_restore]
  )
}

data "huaweicloud_dcs_backups" "backup_format_filter" {
  depends_on    = [huaweicloud_dcs_backup.test]
  instance_id   = huaweicloud_dcs_instance.instance_1.id
  backup_format = huaweicloud_dcs_backup.test.backup_format
}

locals {
  backup_format = huaweicloud_dcs_backup.test.backup_format
}
	
output "backup_format_filter_is_useful" {
  value = length(data.huaweicloud_dcs_backups.backup_format_filter.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dcs_backups.backup_format_filter.backups[*].backup_format : v == local.backup_format]
  )
}
`, testDcsBackup_basic(name))
}

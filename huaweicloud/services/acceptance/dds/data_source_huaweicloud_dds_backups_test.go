package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsBackups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_backups.all"
	dc := acceptance.InitDataSourceCheck(dataSource)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),

					resource.TestCheckOutput("is_backup_id_filter_useful", "true"),
					resource.TestCheckOutput("is_backup_name_filter_useful", "true"),
					resource.TestCheckOutput("is_backup_type_filter_useful", "true"),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("is_instance_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceBackups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dds_backups" "all" {
  depends_on = [huaweicloud_dds_backup.test]
}

// filter by backup_id
data "huaweicloud_dds_backups" "filter_by_backup_id" {
  backup_id = huaweicloud_dds_backup.test.id
}

locals {
  filter_result_by_backup_id = [for v in data.huaweicloud_dds_backups.filter_by_backup_id.backups[*].id : 
    v == huaweicloud_dds_backup.test.id]
}

output "is_backup_id_filter_useful" {
  value = length(local.filter_result_by_backup_id) == 1 && alltrue(local.filter_result_by_backup_id) 
}

// filter by backup_name
data "huaweicloud_dds_backups" "filter_by_backup_name" {
  backup_name = huaweicloud_dds_backup.test.name
}

locals {
  filter_result_by_backup_name = [for v in data.huaweicloud_dds_backups.filter_by_backup_name.backups[*].name : 
    v == huaweicloud_dds_backup.test.name]
}

output "is_backup_name_filter_useful" {
  value = length(local.filter_result_by_backup_name) > 0 && alltrue(local.filter_result_by_backup_name) 
}

// filter by backup_type
data "huaweicloud_dds_backups" "filter_by_backup_type" {
  backup_type = huaweicloud_dds_backup.test.type
}

locals {
  filter_result_by_backup_type = [for v in data.huaweicloud_dds_backups.filter_by_backup_type.backups[*].type : 
    v == huaweicloud_dds_backup.test.type]
}

output "is_backup_type_filter_useful" {
  value = length(local.filter_result_by_backup_type) > 0 && alltrue(local.filter_result_by_backup_type) 
}

// filter by instance_id
data "huaweicloud_dds_backups" "filter_by_instance_id" {
  depends_on = [huaweicloud_dds_backup.test]

  instance_id = huaweicloud_dds_instance.instance.id
}

locals {
  filter_result_by_instance_id = [for v in data.huaweicloud_dds_backups.filter_by_instance_id.backups[*].instance_id : 
    v == huaweicloud_dds_instance.instance.id]
}

output "is_instance_id_filter_useful" {
  value = length(local.filter_result_by_instance_id) > 0 && alltrue(local.filter_result_by_instance_id) 
}

// filter by instance_name
data "huaweicloud_dds_backups" "filter_by_instance_name" {
  depends_on = [huaweicloud_dds_backup.test]

  instance_name = huaweicloud_dds_instance.instance.name
}

locals {
  filter_result_by_instance_name = [for v in data.huaweicloud_dds_backups.filter_by_instance_name.backups[*].instance_name : 
    v == huaweicloud_dds_instance.instance.name]
}

output "is_instance_name_filter_useful" {
  value = length(local.filter_result_by_instance_name) > 0 && alltrue(local.filter_result_by_instance_name) 
}

// filter by status
data "huaweicloud_dds_backups" "filter_by_status" {
  status = huaweicloud_dds_backup.test.status
}

locals {
  filter_result_by_status = [for v in data.huaweicloud_dds_backups.filter_by_status.backups[*].status : 
    v == huaweicloud_dds_backup.test.status]
}

output "is_status_filter_useful" {
  value = length(local.filter_result_by_status) > 0 && alltrue(local.filter_result_by_status) 
}

// filter by description
data "huaweicloud_dds_backups" "filter_by_description" {
  description = huaweicloud_dds_backup.test.description
}

locals {
  filter_result_by_description = [for v in data.huaweicloud_dds_backups.filter_by_description.backups[*].description : 
    v == huaweicloud_dds_backup.test.description]
}

output "is_description_filter_useful" {
  value = length(local.filter_result_by_description) > 0 && alltrue(local.filter_result_by_description) 
}
`, testDdsBackup_basic(name))
}

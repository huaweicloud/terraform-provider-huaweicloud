package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsBackups_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dds_backups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceBackups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.is_instance_restoring"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.backup_method"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.kms_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.deletable"),

					resource.TestCheckOutput("is_backup_id_filter_useful", "true"),
					resource.TestCheckOutput("is_backup_name_filter_useful", "true"),
					resource.TestCheckOutput("is_backup_type_filter_useful", "true"),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					resource.TestCheckOutput("is_instance_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					resource.TestCheckOutput("is_time_filter_useful", "true"),
					resource.TestCheckOutput("is_mode_filter_useful", "true"),
					resource.TestCheckOutput("is_sort_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceBackups_basic() string {
	return `
data "huaweicloud_dds_backups" "test" {}

// filter by backup_id
locals {
  backup_id = data.huaweicloud_dds_backups.test.backups[0].id
}

data "huaweicloud_dds_backups" "filter_by_backup_id" {
  backup_id = local.backup_id
}

output "is_backup_id_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_backup_id.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_backup_id.backups[*].id : v == local.backup_id]
  )
}

// filter by backup_name
locals {
  backup_name = data.huaweicloud_dds_backups.test.backups[0].name
}

data "huaweicloud_dds_backups" "filter_by_backup_name" {
  backup_name = local.backup_name
}

output "is_backup_name_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_backup_name.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_backup_name.backups[*].name : v == local.backup_name]
  )
}

// filter by backup_type
locals {
  backup_type = data.huaweicloud_dds_backups.test.backups[0].type
}

data "huaweicloud_dds_backups" "filter_by_backup_type" {
  backup_type = local.backup_type
}

output "is_backup_type_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_backup_type.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_backup_type.backups[*].type : v == local.backup_type]
  )
}

// filter by instance_id
locals {
  instance_id = data.huaweicloud_dds_backups.test.backups[0].instance_id
}

data "huaweicloud_dds_backups" "filter_by_instance_id" {
  instance_id = local.instance_id
}

output "is_instance_id_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_instance_id.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_instance_id.backups[*].instance_id : v == local.instance_id]
  )
}

// filter by instance_name
locals {
  instance_name = data.huaweicloud_dds_backups.test.backups[0].instance_name
}

data "huaweicloud_dds_backups" "filter_by_instance_name" {
  instance_name = local.instance_name
}

output "is_instance_name_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_instance_name.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_instance_name.backups[*].instance_name : v == local.instance_name]
  )
}

// filter by status
locals {
  status = data.huaweicloud_dds_backups.test.backups[0].status
}

data "huaweicloud_dds_backups" "filter_by_status" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_status.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_status.backups[*].status : v == local.status]
  )
}

// filter by description
locals {
  description = data.huaweicloud_dds_backups.test.backups[0].description
}

data "huaweicloud_dds_backups" "filter_by_description" {
  description = local.description
}

output "is_description_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_description.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_description.backups[*].description : v == local.description]
  )
}

// filter by begin_time and end_time
locals {
  begin_time = data.huaweicloud_dds_backups.test.backups[0].begin_time
  end_time   = data.huaweicloud_dds_backups.test.backups[0].end_time
}

data "huaweicloud_dds_backups" "filter_by_time" {
  begin_time = local.begin_time
  end_time   = local.end_time
}

output "is_time_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_description.backups) > 0
}

// filter by mode
locals {
  mode = data.huaweicloud_dds_backups.test.backups[0].instance_mode
}

data "huaweicloud_dds_backups" "filter_by_mode" {
  mode = local.mode
}

output "is_mode_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_mode.backups) > 0 && alltrue(
    [for v in data.huaweicloud_dds_backups.filter_by_mode.backups[*].instance_mode : v == local.mode]
  )
}

// filter by order_field and order_rule
data "huaweicloud_dds_backups" "filter_by_sort" {
  order_field = "beginTime"
  order_rule  = "asc"
}

locals {
  len     = length(data.huaweicloud_dds_backups.filter_by_mode.backups)
  result1 = data.huaweicloud_dds_backups.test.backups[local.len - 1].id
  result2 = data.huaweicloud_dds_backups.filter_by_sort.backups[0].id
}

output "is_sort_filter_useful" {
  value = length(data.huaweicloud_dds_backups.filter_by_sort.backups) > 0 && (
    local.result1 == local.result2
  )
}
`
}

package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsRecyclingInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_recycling_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsBackupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsRecyclingInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.ha_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.pay_model"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.data_vip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.retained_until"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.recycle_backup_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.recycle_status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.is_serverless"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.deleted_at"),

					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("ha_mode_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_name_filter_is_useful", "true"),
					resource.TestCheckOutput("engine_version_filter_is_useful", "true"),
					resource.TestCheckOutput("pay_model_filter_is_useful", "true"),
					resource.TestCheckOutput("volume_type_filter_is_useful", "true"),
					resource.TestCheckOutput("volume_size_filter_is_useful", "true"),
					resource.TestCheckOutput("data_vip_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("recycle_backup_id_filter_is_useful", "true"),
					resource.TestCheckOutput("recycle_status_filter_is_useful", "true"),
					resource.TestCheckOutput("is_serverless_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsRecyclingInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_recycling_instances" "test" {}

locals {
  instance_id = "%[1]s"
}
data "huaweicloud_rds_recycling_instances" "instance_id_filter" {
  instance_id = "%[1]s"
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.instance_id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.instance_id_filter.instances[*].id : v == local.instance_id]
  )
}

locals {
  name = "test_terraform"
}
data "huaweicloud_rds_recycling_instances" "name_filter" {
  name = "test_terraform"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.name_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.name_filter.instances[*].name : v == local.name]
  )
}

locals {
  ha_mode = "Single"
}
data "huaweicloud_rds_recycling_instances" "ha_mode_filter" {
  ha_mode = "Single"
}
output "ha_mode_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.ha_mode_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.ha_mode_filter.instances[*].ha_mode : v == local.ha_mode]
  )
}

locals {
  engine_name = "postgresql"
}
data "huaweicloud_rds_recycling_instances" "engine_name_filter" {
  engine_name = "postgresql"
}
output "engine_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.engine_name_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.engine_name_filter.instances[*].engine_name : v == local.engine_name]
  )
}

locals {
  engine_version = "12.18"
}
data "huaweicloud_rds_recycling_instances" "engine_version_filter" {
  engine_version = "12.18"
}
output "engine_version_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.engine_version_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.engine_version_filter.instances[*].engine_version : v == local.engine_version]
  )
}

locals {
  pay_model = "0"
}
data "huaweicloud_rds_recycling_instances" "pay_model_filter" {
  pay_model = "0"
}
output "pay_model_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.pay_model_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.pay_model_filter.instances[*].pay_model : v == local.pay_model]
  )
}

locals {
  volume_type = "SSD"
}
data "huaweicloud_rds_recycling_instances" "volume_type_filter" {
  volume_type = "SSD"
}
output "volume_type_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.volume_type_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.volume_type_filter.instances[*].volume_type : v == local.volume_type]
  )
}

locals {
  volume_size = 40
}
data "huaweicloud_rds_recycling_instances" "volume_size_filter" {
  volume_size = "40"
}
output "volume_size_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.volume_size_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.volume_size_filter.instances[*].volume_size : v == local.volume_size]
  )
}

locals {
  data_vip = "192.168.0.235"
}
data "huaweicloud_rds_recycling_instances" "data_vip_filter" {
  data_vip = "192.168.0.235"
}
output "data_vip_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.data_vip_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.data_vip_filter.instances[*].data_vip : v == local.data_vip]
  )
}

locals {
  enterprise_project_id = "0"
}
data "huaweicloud_rds_recycling_instances" "enterprise_project_id_filter" {
  enterprise_project_id = "0"
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.enterprise_project_id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.enterprise_project_id_filter.instances[*].enterprise_project_id
  : v == local.enterprise_project_id]
  )
}

locals {
  recycle_backup_id = "%[2]s"
}
data "huaweicloud_rds_recycling_instances" "recycle_backup_id_filter" {
  recycle_backup_id = "%[2]s"
}
output "recycle_backup_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.recycle_backup_id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.recycle_backup_id_filter.instances[*].recycle_backup_id
  : v == local.recycle_backup_id]
  )
}

locals {
  recycle_status = "COMPLETED"
}
data "huaweicloud_rds_recycling_instances" "recycle_status_filter" {
  recycle_status = "COMPLETED"
}
output "recycle_status_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.recycle_status_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.recycle_status_filter.instances[*].recycle_status : v == local.recycle_status]
  )
}

locals {
  is_serverless = false
}
data "huaweicloud_rds_recycling_instances" "is_serverless_filter" {
  is_serverless = "false"
}
output "is_serverless_filter_is_useful" {
  value = length(data.huaweicloud_rds_recycling_instances.is_serverless_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_recycling_instances.is_serverless_filter.instances[*].is_serverless : v == local.is_serverless]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_BACKUP_ID)
}

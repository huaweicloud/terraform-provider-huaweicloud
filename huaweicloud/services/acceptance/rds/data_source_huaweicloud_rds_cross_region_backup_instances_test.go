package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsCrossRegionBackupInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_cross_region_backup_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckReplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsCrossRegionBackupInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.source_region"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.source_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.destination_region"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.destination_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "backup_instances.0.keep_days"),

					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("source_region_filter_is_useful", "true"),
					resource.TestCheckOutput("source_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("destination_region_filter_is_useful", "true"),
					resource.TestCheckOutput("destination_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("keep_days_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsCrossRegionBackupInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_cross_region_backup_instances" "test" {
  depends_on = [huaweicloud_rds_cross_region_backup_strategy.test] 
}

locals {
  instance_id = huaweicloud_rds_instance.test.id
}
data "huaweicloud_rds_cross_region_backup_instances" "instance_id_filter" {
  depends_on  = [huaweicloud_rds_cross_region_backup_strategy.test]
  instance_id = huaweicloud_rds_instance.test.id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_backup_instances.instance_id_filter.backup_instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_cross_region_backup_instances.instance_id_filter.backup_instances[*].id : v == local.instance_id]
  )
}

locals {
  name = huaweicloud_rds_instance.test.name
}
data "huaweicloud_rds_cross_region_backup_instances" "name_filter" {
  depends_on = [huaweicloud_rds_cross_region_backup_strategy.test]
  name       = huaweicloud_rds_instance.test.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_backup_instances.name_filter.backup_instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_cross_region_backup_instances.name_filter.backup_instances[*].name : v == local.name]
  )
}

locals {
  source_region = "%[2]s"
}
data "huaweicloud_rds_cross_region_backup_instances" "source_region_filter" {
  depends_on    = [huaweicloud_rds_cross_region_backup_strategy.test]
  source_region = "%[2]s"
}
output "source_region_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_backup_instances.source_region_filter.backup_instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_cross_region_backup_instances.source_region_filter.backup_instances[*].source_region :
  v == local.source_region]
  )
}

locals {
  source_project_id = "%[3]s"
}
data "huaweicloud_rds_cross_region_backup_instances" "source_project_id_filter" {
  depends_on        = [huaweicloud_rds_cross_region_backup_strategy.test]
  source_project_id = "%[3]s"
}
output "source_project_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_backup_instances.source_project_id_filter.backup_instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_cross_region_backup_instances.source_project_id_filter.backup_instances[*].source_project_id :
  v == local.source_project_id]
  )
}

locals {
  destination_region = huaweicloud_rds_cross_region_backup_strategy.test.destination_region
}
data "huaweicloud_rds_cross_region_backup_instances" "destination_region_filter" {
  depends_on         = [huaweicloud_rds_cross_region_backup_strategy.test]
  destination_region = huaweicloud_rds_cross_region_backup_strategy.test.destination_region
}
output "destination_region_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_backup_instances.destination_region_filter.backup_instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_cross_region_backup_instances.destination_region_filter.backup_instances[*].destination_region :
  v == local.destination_region]
  )
}

locals {
  destination_project_id = huaweicloud_rds_cross_region_backup_strategy.test.destination_project_id
}
data "huaweicloud_rds_cross_region_backup_instances" "destination_project_id_filter" {
  depends_on             = [huaweicloud_rds_cross_region_backup_strategy.test]
  destination_project_id = huaweicloud_rds_cross_region_backup_strategy.test.destination_project_id
}
output "destination_project_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_backup_instances.destination_project_id_filter.backup_instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_cross_region_backup_instances.destination_project_id_filter.backup_instances[*].destination_project_id :
  v == local.destination_project_id]
  )
}

locals {
  keep_days = huaweicloud_rds_cross_region_backup_strategy.test.keep_days
}
data "huaweicloud_rds_cross_region_backup_instances" "keep_days_filter" {
  depends_on = [huaweicloud_rds_cross_region_backup_strategy.test]
  keep_days  = huaweicloud_rds_cross_region_backup_strategy.test.keep_days
}
output "keep_days_filter_is_useful" {
  value = length(data.huaweicloud_rds_cross_region_backup_instances.keep_days_filter.backup_instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_cross_region_backup_instances.keep_days_filter.backup_instances[*].keep_days : v == local.keep_days]
  )
}
`, testBackupStrategy_basic(name), acceptance.HW_REGION_NAME, acceptance.HW_PROJECT_ID)
}

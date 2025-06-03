package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDrRelationships_basic(t *testing.T) {
	rName := "data.huaweicloud_rds_dr_relationships.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDrRelationships_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.#"),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.0.id"),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.0.status"),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.0.master_instance_id"),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.0.master_region"),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.0.slave_region"),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.0.build_process"),
					resource.TestCheckResourceAttrSet(rName, "instance_dr_infos.0.time"),
					resource.TestCheckOutput("relationship_id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("master_instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("master_region_filter_is_useful", "true"),
					resource.TestCheckOutput("slave_instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("slave_region_filter_is_useful", "true"),
					resource.TestCheckOutput("create_at_start_filter_is_useful", "true"),
					resource.TestCheckOutput("create_at_end_filter_is_useful", "true"),
					resource.TestCheckOutput("order_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_field_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDrRelationships_basic() string {
	return `
data "huaweicloud_rds_dr_relationships" "test" {}

data "huaweicloud_rds_dr_relationships" "relationship_id_filter" {
  relationship_id = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].id
}

locals {
  relationship_id = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].id
}

output "relationship_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.relationship_id_filter.instance_dr_infos) > 0 && alltrue(
    [for v in data.huaweicloud_rds_dr_relationships.relationship_id_filter.instance_dr_infos[*].id : v == local.relationship_id]
  )
}

data "huaweicloud_rds_dr_relationships" "status_filter" {
  status = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].status
}

locals {
  status = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.status_filter.instance_dr_infos) > 0 && alltrue(
    [for v in data.huaweicloud_rds_dr_relationships.status_filter.instance_dr_infos[*].status : v == local.status]
  )
}

data "huaweicloud_rds_dr_relationships" "master_instance_id_filter" {
  master_instance_id = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].master_instance_id
}

locals {
  master_instance_id = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].master_instance_id
}

output "master_instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.master_instance_id_filter.instance_dr_infos) > 0 && alltrue(
    [for v in data.huaweicloud_rds_dr_relationships.master_instance_id_filter.instance_dr_infos[*].master_instance_id :
    v == local.master_instance_id]
  )
}

data "huaweicloud_rds_dr_relationships" "master_region_filter" {
  master_region = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].master_region
}

locals {
  master_region = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].master_region
}

output "master_region_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.master_region_filter.instance_dr_infos) > 0 && alltrue(
    [for v in data.huaweicloud_rds_dr_relationships.master_region_filter.instance_dr_infos[*].master_region :
    v == local.master_region]
  )
}

data "huaweicloud_rds_dr_relationships" "slave_instance_id_filter" {
  slave_instance_id = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].slave_instance_id
}

locals {
  slave_instance_id = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].slave_instance_id
}

output "slave_instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.slave_instance_id_filter.instance_dr_infos) > 0 && alltrue(
    [for v in data.huaweicloud_rds_dr_relationships.slave_instance_id_filter.instance_dr_infos[*].slave_instance_id :
    v == local.slave_instance_id]
  )
}

data "huaweicloud_rds_dr_relationships" "slave_region_filter" {
  slave_region = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].slave_region
}

locals {
  slave_region = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].slave_region
}

output "slave_region_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.slave_region_filter.instance_dr_infos) > 0 && alltrue(
    [for v in data.huaweicloud_rds_dr_relationships.slave_region_filter.instance_dr_infos[*].slave_region : v == local.slave_region]
  )
}

data "huaweicloud_rds_dr_relationships" "create_at_start_filter" {
  create_at_start = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].time
}

output "create_at_start_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.create_at_start_filter.instance_dr_infos) > 0
}

data "huaweicloud_rds_dr_relationships" "create_at_end_filter" {
  create_at_end = data.huaweicloud_rds_dr_relationships.test.instance_dr_infos[0].time
}

output "create_at_end_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.create_at_end_filter.instance_dr_infos) > 0
}

data "huaweicloud_rds_dr_relationships" "order_filter" {
  order = "DESC"
}

output "order_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.order_filter.instance_dr_infos) > 0
}

data "huaweicloud_rds_dr_relationships" "sort_field_filter" {
  sort_field = "time"
}

output "sort_field_filter_is_useful" {
  value = length(data.huaweicloud_rds_dr_relationships.sort_field_filter.instance_dr_infos) > 0
}
`
}

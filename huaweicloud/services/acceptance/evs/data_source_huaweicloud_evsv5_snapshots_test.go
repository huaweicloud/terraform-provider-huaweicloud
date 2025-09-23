package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsv5Snapshots_basic(t *testing.T) {
	var (
		dataSourceName   = "data.huaweicloud_evsv5_snapshots.test"
		dc               = acceptance.InitDataSourceCheck(dataSourceName)
		dataSourceWithId = "data.huaweicloud_evsv5_snapshots.test_with_id"
		dcWithId         = acceptance.InitDataSourceCheck(dataSourceWithId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEvsv5Snapshots_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.encrypted"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.instant_access"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.incremental"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.progress"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.snapshot_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.capacity"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.category"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshots.0.snapshot_chains.0.updated_at"),
					dcWithId.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceWithId, "snapshots.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceWithId, "snapshots.0.id", dataSourceName, "snapshots.0.id"),
					resource.TestCheckOutput("is_snapshot_id_useful", "true"),
					resource.TestCheckOutput("is_volume_id_useful", "true"),
					resource.TestCheckOutput("is_status_useful", "true"),
					resource.TestCheckOutput("is_availability_zone_useful", "true"),
					resource.TestCheckOutput("is_name_useful", "true"),
					resource.TestCheckOutput("is_sort_key_useful", "true"),
					resource.TestCheckOutput("is_sort_dir_useful", "true"),
					resource.TestCheckOutput("is_snapshot_type_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_useful", "true"),
					resource.TestCheckOutput("is_snapshot_chain_id_useful", "true"),
					resource.TestCheckOutput("is_snapshot_group_id_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceEvsv5Snapshots_basic() string {
	return `
data "huaweicloud_evsv5_snapshots" "test" {}

locals {
  snapshot_id           = data.huaweicloud_evsv5_snapshots.test.snapshots.0.id
  volume_id             = data.huaweicloud_evsv5_snapshots.test.snapshots.0.volume_id
  status                = data.huaweicloud_evsv5_snapshots.test.snapshots.0.status
  availability_zone     = data.huaweicloud_evsv5_snapshots.test.snapshots.0.availability_zone
  name                  = data.huaweicloud_evsv5_snapshots.test.snapshots.0.name
  sort_key              = "created_at"
  sort_dir              = "desc"
  snapshot_type         = data.huaweicloud_evsv5_snapshots.test.snapshots.0.snapshot_type
  enterprise_project_id = data.huaweicloud_evsv5_snapshots.test.snapshots.0.enterprise_project_id
  snapshot_chain_id     = data.huaweicloud_evsv5_snapshots.test.snapshots.0.snapshot_chains[0].id
  snapshot_group_id     = data.huaweicloud_evsv5_snapshots.test.snapshots.0.snapshot_group_id
}

data "huaweicloud_evsv5_snapshots" "test_with_id" {
  id = local.snapshot_id
}

output "is_snapshot_id_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_id.snapshots.*.id) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_id.snapshots.*.id : v == local.snapshot_id
  ])
}

data "huaweicloud_evsv5_snapshots" "test_with_volume_id" {
  volume_id = local.volume_id
}

output "is_volume_id_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_volume_id.snapshots.*.volume_id) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_volume_id.snapshots.*.volume_id : v == local.volume_id
  ])
}

data "huaweicloud_evsv5_snapshots" "test_with_status" {
  status = local.status
}

output "is_status_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_status.snapshots.*.status) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_status.snapshots.*.status : v == local.status
  ])
}

data "huaweicloud_evsv5_snapshots" "test_with_availability_zone" {
  availability_zone = local.availability_zone
}

output "is_availability_zone_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_availability_zone.snapshots.*.availability_zone) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_availability_zone.snapshots.*.availability_zone : v == local.availability_zone
  ])
}

data "huaweicloud_evsv5_snapshots" "test_with_name" {
  name = local.name
}

output "is_name_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_name.snapshots.*.name) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_name.snapshots.*.name : v == local.name
  ])
}

data "huaweicloud_evsv5_snapshots" "sort_desc_filter" {
  sort_key = "created_at"
  sort_dir = "desc"
}

data "huaweicloud_evsv5_snapshots" "sort_asc_filter" {
  sort_key = "created_at"
  sort_dir = "asc"
}

locals {
  asc_first_id  = try(data.huaweicloud_evsv5_snapshots.sort_asc_filter.snapshots[0].id, "")
  desc_length   = length(data.huaweicloud_evsv5_snapshots.sort_desc_filter.snapshots)
  desc_last_id  = try(data.huaweicloud_evsv5_snapshots.sort_desc_filter.snapshots[local.desc_length - 1].id, "")
}

output "is_sort_key_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.sort_desc_filter.snapshots) > 0
}

output "is_sort_dir_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.sort_desc_filter.snapshots) > 0 && local.asc_first_id == local.desc_last_id
}

data "huaweicloud_evsv5_snapshots" "test_with_snapshot_type" {
  snapshot_type = local.snapshot_type
}

output "is_snapshot_type_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_snapshot_type.snapshots.*.snapshot_type) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_snapshot_type.snapshots.*.snapshot_type : v == local.snapshot_type
  ])
}

data "huaweicloud_evsv5_snapshots" "test_with_enterprise_project_id" {
  enterprise_project_id = local.enterprise_project_id
}

output "is_enterprise_project_id_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_enterprise_project_id.snapshots.*.enterprise_project_id) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_enterprise_project_id.snapshots.*.enterprise_project_id : v == local.enterprise_project_id
  ])
}

data "huaweicloud_evsv5_snapshots" "test_with_snapshot_chain_id" {
  snapshot_chain_id = local.snapshot_chain_id
}

output "is_snapshot_chain_id_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_snapshot_chain_id.snapshots.*.snapshot_chains) > 0 && alltrue([
    for chains in data.huaweicloud_evsv5_snapshots.test_with_snapshot_chain_id.snapshots.*.snapshot_chains :
      length(chains) == 0 || contains([for c in chains : c.id], local.snapshot_chain_id)
  ])
}

data "huaweicloud_evsv5_snapshots" "test_with_snapshot_group_id" {
  snapshot_group_id = local.snapshot_group_id
}

output "is_snapshot_group_id_useful" {
  value = length(data.huaweicloud_evsv5_snapshots.test_with_snapshot_group_id.snapshots.*.snapshot_group_id) > 0 && alltrue([
    for v in data.huaweicloud_evsv5_snapshots.test_with_snapshot_group_id.snapshots.*.snapshot_group_id : v == local.snapshot_group_id
  ])
}
`
}

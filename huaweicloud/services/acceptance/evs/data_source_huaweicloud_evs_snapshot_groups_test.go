package evs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEvsSnapshotGroups_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_evs_snapshot_groups.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test needs to create an EVS snapshot group before running.
			acceptance.TestAccPreCheckEVSFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEvsSnapshotGroups_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_groups.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_groups.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_groups.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_groups.0.enterprise_project_id"),
					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_key_is_useful", "true"),
					resource.TestCheckOutput("sort_dir_is_useful", "true"),
					resource.TestCheckOutput("server_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceEvsSnapshotGroups_basic = `
data "huaweicloud_evs_snapshot_groups" "test" {}

locals {
  group_id                = data.huaweicloud_evs_snapshot_groups.test.snapshot_groups[0].id
  name                    = data.huaweicloud_evs_snapshot_groups.test.snapshot_groups[0].name
  status                  = data.huaweicloud_evs_snapshot_groups.test.snapshot_groups[0].status
  enterprise_project_id   = data.huaweicloud_evs_snapshot_groups.test.snapshot_groups[0].enterprise_project_id
  server_id               = data.huaweicloud_evs_snapshot_groups.test.snapshot_groups[0].server_id
}

data "huaweicloud_evs_snapshot_groups" "id_filter" {
  id = local.group_id
}

output "id_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshot_groups.id_filter.snapshot_groups.*.id) > 0 && alltrue([
    for v in data.huaweicloud_evs_snapshot_groups.id_filter.snapshot_groups.*.id : v == local.group_id
  ])
}

data "huaweicloud_evs_snapshot_groups" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshot_groups.name_filter.snapshot_groups.*.name) > 0 && alltrue([
    for v in data.huaweicloud_evs_snapshot_groups.name_filter.snapshot_groups.*.name : v == local.name
  ])
}

data "huaweicloud_evs_snapshot_groups" "status_filter" {
  status = local.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshot_groups.status_filter.snapshot_groups.*.status) > 0 && alltrue([
    for v in data.huaweicloud_evs_snapshot_groups.status_filter.snapshot_groups.*.status : v == local.status
  ])
}

data "huaweicloud_evs_snapshot_groups" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(
    data.huaweicloud_evs_snapshot_groups.enterprise_project_id_filter.snapshot_groups.*.enterprise_project_id
  ) > 0 && alltrue([
    for v in data.huaweicloud_evs_snapshot_groups.enterprise_project_id_filter.snapshot_groups.*.enterprise_project_id : 
    v == local.enterprise_project_id
  ])
}

data "huaweicloud_evs_snapshot_groups" "sort_desc_filter" {
  sort_key = "created_at"
  sort_dir = "desc"
}

data "huaweicloud_evs_snapshot_groups" "sort_asc_filter" {
  sort_key = "created_at"
  sort_dir = "asc"
}

locals {
  asc_first_id  = try(data.huaweicloud_evs_snapshot_groups.sort_asc_filter.snapshot_groups[0].id, "")
  desc_length   = length(data.huaweicloud_evs_snapshot_groups.sort_desc_filter.snapshot_groups)
  desc_last_id  = try(data.huaweicloud_evs_snapshot_groups.sort_desc_filter.snapshot_groups[local.desc_length - 1].id, "")
}

output "sort_key_is_useful" {
  value = length(data.huaweicloud_evs_snapshot_groups.sort_desc_filter.snapshot_groups) > 0
}

output "sort_dir_is_useful" {
  value = length(data.huaweicloud_evs_snapshot_groups.sort_desc_filter.snapshot_groups) > 0 && local.asc_first_id == local.desc_last_id
}

data "huaweicloud_evs_snapshot_groups" "server_id_filter" {
  server_id = local.server_id
}

output "server_id_filter_is_useful" {
  value = length(data.huaweicloud_evs_snapshot_groups.server_id_filter.snapshot_groups.*.server_id) > 0 && alltrue([
    for v in data.huaweicloud_evs_snapshot_groups.server_id_filter.snapshot_groups.*.server_id : v == local.server_id
  ])
}
`

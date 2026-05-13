package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataResourcePoolWorkloads_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelartsv2_resource_pool_workloads.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_modelartsv2_resource_pool_workloads.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byStatus   = "data.huaweicloud_modelartsv2_resource_pool_workloads.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		bySort   = "data.huaweicloud_modelartsv2_resource_pool_workloads.filter_by_sort"
		dcBySort = acceptance.InitDataSourceCheck(bySort)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataResourcePoolWorkloads_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "workloads.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "workloads.0.api_version"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.kind"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.type"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.namespace"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.name"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.uid"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.job_uuid"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.status"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.priority"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.running_duration"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.pending_duration"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.pending_position"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.create_time"),
					resource.TestCheckResourceAttrSet(all, "workloads.0.gvk"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					dcBySort.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataResourcePoolWorkloads_base() string {
	return fmt.Sprintf(`
locals {
  resource_pood_ids = split(",", "%[1]s")
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_IDS)
}

func testAccDataResourcePoolWorkloads_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_modelartsv2_resource_pool_workloads" "all" {
  pool_id = local.resource_pood_ids[0]
}

# Filter by type
locals {
  type = data.huaweicloud_modelartsv2_resource_pool_workloads.all.workloads[*].type
}

data "huaweicloud_modelartsv2_resource_pool_workloads" "filter_by_type" {
  pool_id = local.resource_pood_ids[0]
  type    = local.type[0]
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_modelartsv2_resource_pool_workloads.filter_by_type.workloads[*].type :
    v == local.type[0]
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by status
locals {
  status = data.huaweicloud_modelartsv2_resource_pool_workloads.all.workloads[*].status
}

data "huaweicloud_modelartsv2_resource_pool_workloads" "filter_by_status" {
  pool_id = local.resource_pood_ids[0]
  status  = local.status[0]
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_modelartsv2_resource_pool_workloads.filter_by_status.workloads[*].status :
    v == local.status[0]
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by sort
data "huaweicloud_modelartsv2_resource_pool_workloads" "filter_by_sort" {
  pool_id = local.resource_pood_ids[0]
  sort    = "create_time"
  ascend  = true
}

locals {
  sorted_times = data.huaweicloud_modelartsv2_resource_pool_workloads.filter_by_sort.workloads[*].create_time
  is_ascending = length(local.sorted_times) > 1 ? timecmp(local.sorted_times[0], local.sorted_times[1]) <= 0 : true
}

output "is_sort_filter_useful" {
  value = local.is_ascending
}`, testAccDataResourcePoolWorkloads_base())
}

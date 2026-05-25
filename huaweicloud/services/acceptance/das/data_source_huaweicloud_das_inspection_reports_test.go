package das

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccInspectionReports_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_inspection_reports.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByHealthRank   = "data.huaweicloud_das_inspection_reports.filter_by_health_rank"
		dcFilterByHealthRank = acceptance.InitDataSourceCheck(filterByHealthRank)

		filterBySortField   = "data.huaweicloud_das_inspection_reports.filter_by_sort_field"
		dcFilterBySortField = acceptance.InitDataSourceCheck(filterBySortField)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInspectionReports_basic,
				Check: resource.ComposeTestCheckFunc(
					// All reports
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "reports.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "reports.0.task_id"),
					resource.TestCheckResourceAttrSet(all, "reports.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "reports.0.instance_name"),
					resource.TestCheckResourceAttrSet(all, "reports.0.cpu"),
					resource.TestCheckResourceAttrSet(all, "reports.0.mem"),
					resource.TestCheckResourceAttrSet(all, "reports.0.disk_size"),
					resource.TestCheckResourceAttrSet(all, "reports.0.score"),
					resource.TestMatchResourceAttr(all, "reports.0.lost_points_details.#",
						regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "reports.0.lost_points_details.0.deducted_condition"),
					resource.TestCheckResourceAttrSet(all, "reports.0.lost_points_details.0.deducted_formula"),
					resource.TestCheckResourceAttrSet(all, "reports.0.lost_points_details.0.deducted_points"),
					resource.TestCheckResourceAttrSet(all, "reports.0.lost_points_details.0.metric"),
					resource.TestCheckResourceAttrSet(all, "reports.0.lost_points_details.0.metric_value"),
					resource.TestCheckResourceAttrSet(all, "reports.0.lost_points_details.0.risk_level"),
					resource.TestCheckResourceAttrSet(all, "reports.0.lost_points_details.0.suggestions"),

					// Filter by health rank
					dcFilterByHealthRank.CheckResourceExists(),
					resource.TestCheckOutput("is_health_rank_filter_useful", "true"),

					// Filter by sort field
					dcFilterBySortField.CheckResourceExists(),
					resource.TestCheckOutput("is_field_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccInspectionReports_basic = `
locals {
  start_time = "2000-01-01T00:00:00+08:00"
  end_time   = "2099-01-01T00:00:00+08:00"
}

data "huaweicloud_das_inspection_reports" "all" {
  start_time     = local.start_time
  end_time       = local.end_time
  datastore_type = "MySQL"
}

# Filter by health rank
data "huaweicloud_das_inspection_reports" "filter_by_health_rank" {
  start_time     = local.start_time
  end_time       = local.end_time
  datastore_type = "MySQL"
  health_rank    = "healthy"
}

locals {
  health_rank_filter_result = [
    for v in data.huaweicloud_das_inspection_reports.filter_by_health_rank.reports : v.health_rank == "healthy"
  ]
}

output "is_health_rank_filter_useful" {
  value = length(local.health_rank_filter_result) > 0 && alltrue(local.health_rank_filter_result)
}

# Filter by sort field
data "huaweicloud_das_inspection_reports" "filter_by_sort_field" {
  start_time     = local.start_time
  end_time       = local.end_time
  datastore_type = "MySQL"
  sort_field     = "create_at"
  asc            = false
}

locals {
  field_filter_result = [
    for v in data.huaweicloud_das_inspection_reports.filter_by_sort_field.reports : v.score != null
  ]
}

output "is_field_filter_useful" {
  value = length(local.field_filter_result) > 0 && alltrue(local.field_filter_result)
}
`

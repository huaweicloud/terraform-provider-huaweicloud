package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsCompareLineOverview_basic(t *testing.T) {
	dataSource := "data.huaweicloud_drs_compare_line_overview.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
			acceptance.TestAccPreCheckDrsCompareJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsCompareLineOverview_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.0.source_db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.0.target_db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.0.compare_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.0.compare_end_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.0.data_inconsistent_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_compare_overview_infos.0.uncomparable_num"),
					resource.TestCheckOutput("filter_by_status_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDrsCompareLineOverview_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_line_overview" "test" {
  job_id         = "%[1]s"
  compare_job_id = "%[2]s"
}

locals {
  status = data.huaweicloud_drs_compare_line_overview.test.data_compare_overview_infos[0].status
}

# Filter by status
data "huaweicloud_drs_compare_line_overview" "filter_by_status" {
  job_id         = "%[1]s"
  compare_job_id = "%[2]s"
  status         = local.status
}

output "filter_by_status_is_useful" {
  value = length(data.huaweicloud_drs_compare_line_overview.filter_by_status.data_compare_overview_infos) > 0 && alltrue(
    [for info in data.huaweicloud_drs_compare_line_overview.filter_by_status.data_compare_overview_infos : 
      info.status == local.status]
  )
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}

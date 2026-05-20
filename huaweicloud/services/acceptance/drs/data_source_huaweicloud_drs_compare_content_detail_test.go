package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsCompareContentDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_drs_compare_content_detail.test"
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
				Config: testAccDataSourceDrsCompareContentDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.source_db"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.target_db"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.source_table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.target_table_name"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.source_row_num"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.target_row_num"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.difference_row_num"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.line_compare_result"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.content_compare_result"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.complete_shard_count"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.total_shard_count"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.progress"),

					resource.TestCheckOutput("filter_by_target_db_is_useful", "true"),
					resource.TestCheckOutput("filter_by_db_name_is_useful", "true"),
					resource.TestCheckOutput("filter_by_type_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDrsCompareContentDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_content_detail" "test" {
  job_id         = "%s"
  compare_job_id = "%s"
}

locals {
  target_db_name = data.huaweicloud_drs_compare_content_detail.test.content_compare_result_infos[0].target_db
}

# Filter by target_db_name
data "huaweicloud_drs_compare_content_detail" "filter_by_target_db" {
  job_id         = "%[1]s"
  compare_job_id = "%[2]s"
  target_db_name = local.target_db_name
}

output "filter_by_target_db_is_useful" {
  value = length(data.huaweicloud_drs_compare_content_detail.filter_by_target_db.content_compare_result_infos) > 0 && alltrue(
    [for info in data.huaweicloud_drs_compare_content_detail.filter_by_target_db.content_compare_result_infos : 
      info.target_db == local.target_db_name]
  )
}

locals {
  db_name = data.huaweicloud_drs_compare_content_detail.test.content_compare_result_infos[0].source_db
}

# Filter by db_name
data "huaweicloud_drs_compare_content_detail" "filter_by_db_name" {
  job_id         = "%[1]s"
  compare_job_id = "%[2]s"
  db_name        = local.db_name
}

output "filter_by_db_name_is_useful" {
  value = length(data.huaweicloud_drs_compare_content_detail.filter_by_db_name.content_compare_result_infos) > 0 && alltrue(
    [for info in data.huaweicloud_drs_compare_content_detail.filter_by_db_name.content_compare_result_infos : 
      info.source_db == local.db_name]
  )
}

# Filter by type
data "huaweicloud_drs_compare_content_detail" "filter_by_type" {
  job_id         = "%[1]s"
  compare_job_id = "%[2]s"
  type           = "compare"
}

output "filter_by_type_is_useful" {
  value = length(data.huaweicloud_drs_compare_content_detail.filter_by_type.content_compare_result_infos) > 0
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}

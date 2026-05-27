package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsCompareLineDetail_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_drs_compare_line_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobId(t)
			acceptance.TestAccPreCheckDrsCompareJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDrsCompareLineDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_line_compare_result_infos.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_line_compare_result_infos.0.source_table_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_line_compare_result_infos.0.source_row_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_line_compare_result_infos.0.target_table_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_line_compare_result_infos.0.target_row_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_line_compare_result_infos.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "table_line_compare_result_infos.0.difference_row_num"),
				),
			},
		},
	})
}

func testDataSourceDrsCompareLineDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_line_detail" "test" {
  job_id         = "%s"
  compare_job_id = "%s"
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}

package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsCompareContentOverview_basic(t *testing.T) {
	dataSource := "data.huaweicloud_drs_compare_content_overview.test"
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
				Config: testAccDataSourceDrsCompareContentOverview_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.source_db"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.target_db"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.compare_num"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.compare_end_num"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.data_inconsistent_num"),
					resource.TestCheckResourceAttrSet(dataSource, "content_compare_result_infos.0.uncomparable_num"),
				),
			},
		},
	})
}

func testAccDataSourceDrsCompareContentOverview_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_content_overview" "test" {
  job_id         = "%[1]s"
  compare_job_id = "%[2]s"
}
`, acceptance.HW_DRS_JOB_ID, acceptance.HW_DRS_COMPARE_JOB_ID)
}

package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDrsCompareResult_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_compare_result.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDrsJobIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDrsCompareResult_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "object_level_compare_results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "line_compare_results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "content_compare_results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "compare_task_list_results.#"),
				),
			},
		},
	})
}

func testAccDataSourceDrsCompareResult_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_drs_compare_result" "test" {
  job_id       = "%s"
  current_page = 1
  per_page     = 10
}
`, acceptance.HW_DRS_JOB_IDS)
}

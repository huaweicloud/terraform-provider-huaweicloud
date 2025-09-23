package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsTaskConsistencyResults_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_task_consistency_results.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSmsTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsTaskConsistencyResults_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.finished_time"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.check_result"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.consistency_result.#"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.consistency_result.0.dir_check"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.consistency_result.0.num_total_files"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.consistency_result.0.num_different_files"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.consistency_result.0.num_target_miss_files"),
					resource.TestCheckResourceAttrSet(dataSource, "result_list.0.consistency_result.0.num_target_more_files"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsTaskConsistencyResults_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_sms_task_consistency_results" "test" {
  task_id = "%s"
}
`, acceptance.HW_SMS_TASK_ID)
}

package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsBackgroundTaskDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_background_task_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcsBackgroundTaskDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "progress"),
					resource.TestCheckResourceAttrSet(dataSource, "remain_time"),
					resource.TestCheckResourceAttrSet(dataSource, "step_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "step_details.0.step_id"),
					resource.TestCheckResourceAttrSet(dataSource, "step_details.0.step_name"),
					resource.TestCheckResourceAttrSet(dataSource, "step_details.0.step_status"),
					resource.TestCheckResourceAttrSet(dataSource, "step_details.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "step_details.0.sub_step_details.#"),
				),
			},
		},
	})
}

func testDataSourceDcsBackgroundTaskDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dcs_background_task_detail" "test" {
  instance_id = "%s"
  task_id     = "%s"
}
`, acceptance.HW_DCS_INSTANCE_ID, acceptance.HW_DCS_BACKGROUND_TASK_ID)
}

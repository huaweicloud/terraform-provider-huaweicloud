package rocketmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataBackgroundTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dms_rocketmq_background_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byTime   = "data.huaweicloud_dms_rocketmq_background_tasks.filter_by_time"
		dcByTime = acceptance.InitDataSourceCheck(byTime)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataBackgroundTasks_notFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccDataBackgroundTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "tasks.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.params"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.updated_at"),
					dcByTime.CheckResourceExists(),
					resource.TestCheckOutput("is_filter_by_time_useful", "true"),
				),
			},
		},
	})
}

func testAccDataBackgroundTasks_notFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_background_tasks" "test" {
  instance_id = "%s"
}
`, randomId)
}

func testAccDataBackgroundTasks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_background_tasks" "test" {
  instance_id = "%[1]s"
}

locals {
  created_at = data.huaweicloud_dms_rocketmq_background_tasks.test.tasks[0].created_at
  begin_time = formatdate("YYYYMMDDhhmmss", local.created_at)
}

data "huaweicloud_dms_rocketmq_background_tasks" "filter_by_time" {
  instance_id = "%[1]s"
  begin_time  = local.begin_time
  end_time    = formatdate("YYYYMMDDhhmmss", timeadd(local.created_at, "2h"))
}

locals {
  filter_by_time_result = [for v in data.huaweicloud_dms_rocketmq_background_tasks.filter_by_time.tasks :
  formatdate("YYYYMMDDhhmmss", v.created_at) >= local.begin_time]
}

output "is_filter_by_time_useful" {
  value = length(local.filter_by_time_result) > 0 && alltrue(local.filter_by_time_result)
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID)
}

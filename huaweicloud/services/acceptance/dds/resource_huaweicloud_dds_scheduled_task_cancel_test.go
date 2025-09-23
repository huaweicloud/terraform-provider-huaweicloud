package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSScheduledTaskCancel_basic(t *testing.T) {
	resourceName := "huaweicloud_dds_scheduled_task_cancel.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSScheduledTasksEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDDSScheduledTaskCancel_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "job_status", "Canceled"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_status"),
					resource.TestCheckResourceAttrSet(resourceName, "job_name"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "end_time"),
				),
			},
		},
	})
}

const testAccDDSScheduledTaskCancel_basic = `
data "huaweicloud_dds_scheduled_tasks" "test" {
  job_status = "Pending"
}

resource "huaweicloud_dds_scheduled_task_cancel" "test" {
  job_id = try(data.huaweicloud_dds_scheduled_tasks.test.schedules[0].job_id, "")

  lifecycle {
    ignore_changes = [
      job_id,
    ]
  }
}`

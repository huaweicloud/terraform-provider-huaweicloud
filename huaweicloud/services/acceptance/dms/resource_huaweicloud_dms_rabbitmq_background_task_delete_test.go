package dms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsRabbitMQBackgroundTaskDelete_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRabbitMQBackgroundTaskDelete_basic(rName),
			},
		},
	})
}

func testAccDmsRabbitMQBackgroundTaskDelete_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_rabbitmq_background_task_delete" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  task_id     = try(data.huaweicloud_dms_rabbitmq_background_tasks.test.tasks[0].id, "")

  lifecycle {
    ignore_changes = [
      task_id,
    ]
  }
}`, testDataSourceDmsRabbitMQBackgroundTasks_basic(rName))
}

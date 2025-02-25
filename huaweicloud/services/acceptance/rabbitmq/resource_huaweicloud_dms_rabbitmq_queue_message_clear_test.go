package rabbitmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRabbitmqQueueMessageClear_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRabbitmqQueueMessageClear_basic(rName),
			},
		},
	})
}

func testAccRabbitmqQueueMessageClear_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rabbitmq_queue_message_clear" "test" {
  depends_on = [huaweicloud_dms_rabbitmq_queue.test]

  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  vhost       = huaweicloud_dms_rabbitmq_vhost.test.name
  queue       = huaweicloud_dms_rabbitmq_queue.test.name
}`, testRabbitmqQueue_basic(rName))
}

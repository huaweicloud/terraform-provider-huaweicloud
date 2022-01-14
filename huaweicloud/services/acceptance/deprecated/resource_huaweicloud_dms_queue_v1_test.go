package deprecated

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/dms/v1/queues"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDmsQueueFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DMS client(V1): %s", err)
	}

	return queues.Get(client, state.Primary.ID, false).Extract()
}

func TestAccDmsQueuesV1_basic(t *testing.T) {
	var queue queues.Queue
	var queueName = acceptance.RandomAccResourceName()
	var resourceName = "huaweicloud_dms_queue_v1.queue_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&queue,
		getDmsQueueFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckDeprecated(t)
			acceptance.TestAccPreCheckDms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Queue_basic(queueName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", queueName),
					resource.TestCheckResourceAttr(resourceName, "queue_mode", "NORMAL"),
				),
			},
		},
	})
}

func TestAccDmsQueuesV1_FIFOmode(t *testing.T) {
	var queue queues.Queue
	var queueName = acceptance.RandomAccResourceName()
	var resourceName = "huaweicloud_dms_queue_v1.queue_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&queue,
		getDmsQueueFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckDeprecated(t)
			acceptance.TestAccPreCheckDms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Queue_FIFOmode(queueName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", queueName),
					resource.TestCheckResourceAttr(resourceName, "description", "test create dms queue"),
					resource.TestCheckResourceAttr(resourceName, "queue_mode", "FIFO"),
					resource.TestCheckResourceAttr(resourceName, "redrive_policy", "enable"),
					resource.TestCheckResourceAttr(resourceName, "max_consume_count", "80"),
				),
			},
		},
	})
}

func testAccDmsV1Queue_basic(queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_queue_v1" "queue_1" {
  name = "%s"
}`, queueName)
}

func testAccDmsV1Queue_FIFOmode(queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_queue_v1" "queue_1" {
  name              = "%s"
  description       = "test create dms queue"
  queue_mode        = "FIFO"
  redrive_policy    = "enable"
  max_consume_count = 80
}`, queueName)
}

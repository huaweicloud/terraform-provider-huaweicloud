package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKafkaMessageProduce_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaMessageProduce_basic(rName),
			},
		},
	})
}

func testAccKafkaMessageProduce_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_message_produce" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.topic]

  instance_id = huaweicloud_dms_kafka_instance.test.id
  topic       = huaweicloud_dms_kafka_topic.topic.name
  body        = "test"

  property_list {
    name  = "KEY"
    value = "testKey"
  }

  property_list {
    name  = "PARTITION"
    value = "1"
  }
}`, testAccDmsKafkaTopic_basic(rName))
}

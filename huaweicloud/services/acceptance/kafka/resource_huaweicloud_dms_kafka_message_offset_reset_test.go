package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKafkaMessageOffsetReset_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaMessageOffsetReset_basic(),
			},
		},
	})
}

func testAccKafkaMessageOffsetReset_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  topic       = ""
  partition   = -1
  timestamp   = 0
}`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}

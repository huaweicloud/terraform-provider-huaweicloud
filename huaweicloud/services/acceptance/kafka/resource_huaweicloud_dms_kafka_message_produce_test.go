package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMessageProduce_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccMessageProduce_basic(rName),
			},
		},
	})
}

func testAccMessageProduce_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  partitions  = 3
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccMessageProduce_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_message_produce" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id = huaweicloud_dms_kafka_topic.test.instance_id
  topic       = huaweicloud_dms_kafka_topic.test.name
  body        = "test"

  property_list {
    name  = "KEY"
    value = "testKey"
  }

  property_list {
    name  = "PARTITION"
    value = "1"
  }
}`, testAccMessageProduce_base(name))
}

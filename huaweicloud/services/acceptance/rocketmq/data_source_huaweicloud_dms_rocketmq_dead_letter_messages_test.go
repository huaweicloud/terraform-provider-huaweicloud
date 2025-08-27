package rocketmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, ensure that the consumer group is online and that there are dead letter messages generated.
func TestAccDataDeadLetterMessages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rocketmq_dead_letter_messages.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQInstanceID(t)
			acceptance.TestAccPreCheckDMSRocketMQGroupName(t)
			acceptance.TestAccPreCheckDMSRocketMQDeadLetterMessageIDs(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDeadLetterMessages_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config:      testAccDeadLetterMessages_topicNotFound(),
				ExpectError: regexp.MustCompile(`Query topic failed. No topic route info in name server for topic`),
			},
			{
				Config: testDataSourceDeadLetterMessages_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "messages.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.msg_id"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.topic"),
					resource.TestMatchResourceAttr(dataSource, "messages.0.store_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "messages.0.born_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.reconsume_times"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.body_crc"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.store_size"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.born_host"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.store_host"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.queue_id"),
					resource.TestCheckResourceAttrSet(dataSource, "messages.0.queue_offset"),
					resource.TestMatchResourceAttr(dataSource, "messages.0.property_list.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDeadLetterMessages_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_dead_letter_messages" "test" {
  instance_id = "%[1]s"
  topic       = "%[2]s"
  msg_id_list = split(",", "%[3]s")
}
`, randomId, acceptance.HW_DMS_ROCKETMQ_GROUP_NAME, acceptance.HW_DMS_ROCKETMQ_DEAD_LETTER_MESSAGE_IDs)
}

func testAccDeadLetterMessages_topicNotFound() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_dead_letter_messages" "test" {
  instance_id = "%[1]s"
  topic       = "not_found"
  msg_id_list = split(",", "%[2]s")
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID, acceptance.HW_DMS_ROCKETMQ_DEAD_LETTER_MESSAGE_IDs)
}

func testDataSourceDeadLetterMessages_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rocketmq_dead_letter_messages" "test" {
  instance_id = "%[1]s"
  topic       = "%%DLQ%%%[2]s"
  msg_id_list = split(",", "%[3]s")
}
`, acceptance.HW_DMS_ROCKETMQ_INSTANCE_ID,
		acceptance.HW_DMS_ROCKETMQ_GROUP_NAME,
		acceptance.HW_DMS_ROCKETMQ_DEAD_LETTER_MESSAGE_IDs,
	)
}

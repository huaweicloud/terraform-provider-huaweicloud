package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkaTopicProducers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_topic_producers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsKafkaTopicProducers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "producers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "producers.0.producer_address"),
					resource.TestCheckResourceAttrSet(dataSource, "producers.0.broker_address"),
					resource.TestCheckResourceAttrSet(dataSource, "producers.0.join_time"),
				),
			},
		},
	})
}

func testDataSourceDmsKafkaTopicProducers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_topic_producers" "test" {
  instance_id = "%[1]s"
  topic       = "%[2]s"
}`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

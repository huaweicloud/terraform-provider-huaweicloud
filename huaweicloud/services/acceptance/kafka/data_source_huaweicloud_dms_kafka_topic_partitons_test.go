package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkaTopicPartitions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_topic_partitions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsKafkaTopicPartitions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "partitions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "partitions.0.partition"),
				),
			},
		},
	})
}

func testDataSourceDmsKafkaTopicPartitions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_topic_partitions" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.topic]

  instance_id = huaweicloud_dms_kafka_instance.test.id
  topic       = huaweicloud_dms_kafka_topic.topic.name
}
`, testAccDmsKafkaTopic_basic(name))
}

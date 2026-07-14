package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTopicPartitions_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_dms_kafka_topic_partitions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTopicPartitions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "partitions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet(dataSource, "partitions.0.partition"),
				),
			},
		},
	})
}

func testAccDataSourcePartitions_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  partitions  = 3
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testDataSourceTopicPartitions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_kafka_topic_partitions" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id = huaweicloud_dms_kafka_topic.test.instance_id
  topic       = huaweicloud_dms_kafka_topic.test.name
}
`, testAccDataSourcePartitions_basic(name))
}

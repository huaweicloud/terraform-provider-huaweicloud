package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkaTopics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_topics.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaTopics_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "topics.#"),
					resource.TestCheckOutput("name_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkaTopics_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_topics" "all" {
  depends_on = [huaweicloud_dms_kafka_topic.topic]

  instance_id = huaweicloud_dms_kafka_instance.test.id
}

// filter
data "huaweicloud_dms_kafka_topics" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.topic]

  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = huaweicloud_dms_kafka_topic.topic.name
}

locals {
  filter_results = data.huaweicloud_dms_kafka_topics.test
}

output "name_validation" {
  value = alltrue([for v in local.filter_results.topics[*].name : strcontains(v, huaweicloud_dms_kafka_topic.topic.name)])
}
`, testAccDmsKafkaTopic_basic(name))
}

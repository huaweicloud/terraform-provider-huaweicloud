package rocketmq

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsRocketMQTopics_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_dms_rocketmq_topics.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsRocketMQSearchTopic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.0.total_read_queue_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.0.total_write_queue_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "topics.0.permission"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("total_read_queue_num_filter_is_useful", "true"),
					resource.TestCheckOutput("total_write_queue_num_filter_is_useful", "true"),
					resource.TestCheckOutput("permission_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDmsRocketMQSearchTopic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_rocketmq_topics" "test" {
  depends_on  = [huaweicloud_dms_rocketmq_topic.test]
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
}

data "huaweicloud_dms_rocketmq_topics" "name_filter" {
  depends_on  = [huaweicloud_dms_rocketmq_topic.test]
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = "%[2]s"
}
  
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_topics.name_filter.topics) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_topics.name_filter.topics[*].name : v == "%[2]s"]
  )  
}

data "huaweicloud_dms_rocketmq_topics" "total_read_queue_num_filter" {
  depends_on           = [huaweicloud_dms_rocketmq_topic.test]
  instance_id          = huaweicloud_dms_rocketmq_instance.test.id
  total_read_queue_num = huaweicloud_dms_rocketmq_topic.test.total_read_queue_num
}

locals {
  total_read_queue_num = huaweicloud_dms_rocketmq_topic.test.total_read_queue_num
}
	
output "total_read_queue_num_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_topics.total_read_queue_num_filter.topics) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_topics.total_read_queue_num_filter.topics[*].total_read_queue_num : v == local.total_read_queue_num]
  )  
}

data "huaweicloud_dms_rocketmq_topics" "total_write_queue_num_filter" {
  depends_on            = [huaweicloud_dms_rocketmq_topic.test]
  instance_id           = huaweicloud_dms_rocketmq_instance.test.id
  total_write_queue_num = huaweicloud_dms_rocketmq_topic.test.total_write_queue_num
}

locals {
  total_write_queue_num = huaweicloud_dms_rocketmq_topic.test.total_write_queue_num
}
	
output "total_write_queue_num_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_topics.total_write_queue_num_filter.topics) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_topics.total_write_queue_num_filter.topics[*].total_write_queue_num : v == local.total_write_queue_num]
  )
}

data "huaweicloud_dms_rocketmq_topics" "permission_filter" {
  depends_on  = [huaweicloud_dms_rocketmq_topic.test]
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  permission  = huaweicloud_dms_rocketmq_topic.test.permission
}
  
locals {
  permission = huaweicloud_dms_rocketmq_topic.test.permission
}
	  
output "permission_filter_is_useful" {
  value = length(data.huaweicloud_dms_rocketmq_topics.permission_filter.topics) > 0 && alltrue(
  [for v in data.huaweicloud_dms_rocketmq_topics.permission_filter.topics[*].permission : v == local.permission]
  )
}

`, testDmsRocketMQTopic_basic(name), name)
}

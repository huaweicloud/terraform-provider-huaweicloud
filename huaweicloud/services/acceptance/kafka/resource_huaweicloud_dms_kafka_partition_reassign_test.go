package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKafkaPartitionReassign_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_kafka_partition_reassign.test"

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaPartitionReassign_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "task_id"),
				),
			},
			{
				Config: testAccKafkaPartitionReassign_automatical(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "task_id"),
				),
			},
			{
				Config: testAccKafkaPartitionReassign_automatical(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "reassignment_time"),
				),
			},
		},
	})
}

func testAccKafkaPartitionReassign_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%s"
  partitions  = 2
  replicas    = 3
}

resource "huaweicloud_dms_kafka_partition_reassign" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id = huaweicloud_dms_kafka_instance.test.id
  
  reassignments {
    topic = huaweicloud_dms_kafka_topic.test.name

    assignment {
      partition         = 0
      partition_brokers = [0,1,2]
    }

    assignment {
      partition         = 1
      partition_brokers = [2,0,1]
    }
  }
}`, testAccKafkaInstance_newFormat(rName), rName)
}

func testAccKafkaPartitionReassign_automatical(rName string, timeEstimate bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%s"
  partitions  = 2
  replicas    = 3
}

resource "huaweicloud_dms_kafka_partition_reassign" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id   = huaweicloud_dms_kafka_instance.test.id
  throttle      = -1
  time_estimate = %t
  
  reassignments {
    topic              = huaweicloud_dms_kafka_topic.test.name
    brokers            = [0,1,2]
    replication_factor = 3
  }
}`, testAccKafkaInstance_newFormat(rName), rName, timeEstimate)
}

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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
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

func testAccKafkaPartitionReassign_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  partitions  = 2
  replicas    = 3
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccKafkaPartitionReassign_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_partition_reassign" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id = "%[2]s"
  
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
}`, testAccKafkaPartitionReassign_base(rName), acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

func testAccKafkaPartitionReassign_automatical(rName string, timeEstimate bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_partition_reassign" "test" {
  depends_on = [huaweicloud_dms_kafka_topic.test]

  instance_id   = "%[2]s"
  throttle      = -1
  time_estimate = %[3]t
  
  reassignments {
    topic              = huaweicloud_dms_kafka_topic.test.name
    brokers            = [0,1,2]
    replication_factor = 3
  }
}`, testAccKafkaPartitionReassign_base(rName), acceptance.HW_DMS_KAFKA_INSTANCE_ID, timeEstimate)
}

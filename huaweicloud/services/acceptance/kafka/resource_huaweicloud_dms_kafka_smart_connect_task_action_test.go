package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKafkaSmartConnectTaskAction_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	// Avoid CheckDestroy
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaSmartConnectTaskAction_basic(rName, "start"),
			},
			{
				Config: testAccKafkaSmartConnectTaskAction_basic(rName, "pause"),
			},
			{
				Config: testAccKafkaSmartConnectTaskAction_basic(rName, "resume"),
			},
			{
				Config: testAccKafkaSmartConnectTaskAction_basic(rName, "restart"),
			},
			{
				Config: testAccKafkaSmartConnectTaskAction_basic(rName, "pause"),
			},
			{
				Config: testAccKafkaSmartConnectTaskAction_basic(rName, "restart"),
			},
		},
	})
}

func testAccKafkaSmartConnectTaskAction_basic(rName, action string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_smart_connect_task_action" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  task_id     = huaweicloud_dms_kafkav2_smart_connect_task.test.id
  action      = "%[2]s"
}`, testDmsKafkaSmartConnectTaskActionBase(rName), action)
}

func testDmsKafkaSmartConnectTaskActionBase(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[2]s"
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_dms_kafka_smart_connect" "test" {
  instance_id       = huaweicloud_dms_kafka_instance.test.id
  storage_spec_code = "dms.physical.storage.high.v2"
  node_count        = 2
  bandwidth         = "100MB"
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[2]s"
  partitions  = 10
  aging_time  = 36
}

resource "huaweicloud_dms_kafkav2_smart_connect_task" "test" {
  depends_on = [huaweicloud_dms_kafka_smart_connect.test, huaweicloud_dms_kafka_topic.test]

  instance_id      = huaweicloud_dms_kafka_instance.test.id
  task_name        = "%[2]s"
  destination_type = "OBS_SINK"
  topics           = [huaweicloud_dms_kafka_topic.test.name]
  start_later      = true

  destination_task {
    consumer_strategy     = "latest"
    destination_file_type = "TEXT"
    access_key            = "%[3]s"
    secret_key            = "%[4]s"
    obs_bucket_name       = huaweicloud_obs_bucket.test.bucket
    partition_format      = "yyyy/MM/dd/HH/mm"
    record_delimiter      = ";"
    deliver_time_interval = 300
  }
}`, testAccKafkaInstance_newFormat(rName), rName, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

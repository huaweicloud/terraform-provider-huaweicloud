package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRegisterKafkaInstance_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsKafkaInstanceIds(t)
			acceptance.TestAccPreCheckLtsKafkaInstancePsw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: resourceRegisterKafkaInstance_basic(name),
			},
		},
	})
}

func resourceRegisterKafkaInstance_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}
`, name)
}

func resourceRegisterKafkaInstance_verity(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = local.kafka_instance.id
  name        = "%[1]s"
  partitions  = 3
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.test.id
  }

  log_transfer_info {
    log_transfer_type   = "DMS"
    log_transfer_mode   = "realTime"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      kafka_id    = local.kafka_instance.id
      kafka_topic = huaweicloud_dms_kafka_topic.test.id
    }
  }

  depends_on = [huaweicloud_lts_register_kafka_instance.test]
}
`, name)
}

func resourceRegisterKafkaInstance_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_kafka_instances" test {
  instance_id = element(split(",", "%[2]s"), 0)
}

locals {
  kafka_instance = data.huaweicloud_dms_kafka_instances.test.instances[0]
}

resource "huaweicloud_lts_register_kafka_instance" "test" {
  instance_id = local.kafka_instance.id
  kafka_name  = local.kafka_instance.name

  connect_info {
    user_name = local.kafka_instance.access_user
    pwd       = "%[3]s"
  }
}

# Verify that the Kafka instance accessed by the encrypted data has been successfully registered to LTS.
%[4]s
`, resourceRegisterKafkaInstance_base(name),
		acceptance.HW_LTS_KAFKA_INSTANCE_IDS,
		acceptance.HW_LTS_KAFKA_INSTANCE_PASSWORD,
		resourceRegisterKafkaInstance_verity(name))
}

func TestAccRegisterKafkaInstance_notSSL(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsKafkaInstanceIds(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: resourceRegisterKafkaInstance_notSSL(name),
			},
		},
	})
}

func resourceRegisterKafkaInstance_notSSL(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_kafka_instances" test {
  instance_id = element(split(",", "%[2]s"), 1)
}

locals {
  kafka_instance = data.huaweicloud_dms_kafka_instances.test.instances[0]
}

resource "huaweicloud_lts_register_kafka_instance" "test" {
  instance_id = local.kafka_instance.id
  kafka_name  = local.kafka_instance.name
}

# Verify that the Kafka instance by plaintext access has been successfully registered to LTS.
%[3]s
`, resourceRegisterKafkaInstance_base(name),
		acceptance.HW_LTS_KAFKA_INSTANCE_IDS,
		resourceRegisterKafkaInstance_verity(name))
}

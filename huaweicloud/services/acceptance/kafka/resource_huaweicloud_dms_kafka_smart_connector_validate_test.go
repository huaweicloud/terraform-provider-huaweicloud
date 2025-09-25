package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccSmartConnectorValidate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccSmartConnectorValidate_basic(),
			},
		},
	})
}

func testAccSmartConnectorValidate_base(name string) string {
	return fmt.Sprintf(`
%[1]s

# Defalut rule cannot be deleted. Otherwise, Kafka instance network will be disconnected.
resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "cluster"
}

locals {
  flavor = try(data.huaweicloud_dms_kafka_flavors.test.flavors[0], {})
}

resource "huaweicloud_dms_kafka_instance" "test" {
  count = 2

  name              = "%[2]s${count.index}"
  flavor_id         = local.flavor.id
  engine_version    = "2.7"
  storage_space     = local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node
  storage_spec_code = local.flavor.ios[0].storage_spec_code
  broker_num        = 3

  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  access_user        = "%[2]s"
  password           = "%[3]s"
  enabled_mechanisms = ["SCRAM-SHA-512"]

  port_protocol {
    private_sasl_ssl_enable = true
  }
}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test[0].id
  name        = "%[2]s"
  partitions  = 10
}

resource "huaweicloud_dms_kafka_smart_connect" "test" {
  instance_id       = huaweicloud_dms_kafka_instance.test[0].id
  storage_spec_code = local.flavor.ios[0].storage_spec_code
}
`, common.TestVpc(name), name, acceptance.RandomPassword())
}

func testAccSmartConnectorValidate_basic() string {
	name := acceptance.RandomAccResourceNameWithDash()
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_smart_connector_validate" "test" {
  instance_id = huaweicloud_dms_kafka_instance.test[0].id
  type        = "KAFKA_REPLICATOR_SOURCE"

  task {
    current_cluster_name          = "%[2]sA"
    cluster_name                  = "%[2]sB"
    user_name                     = huaweicloud_dms_kafka_instance.test[1].access_user
    password                      = huaweicloud_dms_kafka_instance.test[1].password
    sasl_mechanism                = "SCRAM-SHA-512"
    instance_id                   = huaweicloud_dms_kafka_instance.test[1].id
    security_protocol             = "SASL_SSL"
    direction                     = "push"
    sync_consumer_offsets_enabled = true
    replication_factor            = 3
    task_num                      = 2
    provenance_header_enabled     = true
    consumer_strategy             = "earliest"
    compression_type              = "none"
    topics_mapping                = "${huaweicloud_dms_kafka_topic.test.name}:%[2]s-test"
  }

  depends_on = [
    huaweicloud_dms_kafka_smart_connect.test,
    huaweicloud_dms_kafka_topic.test,
  ]
}
`, testAccSmartConnectorValidate_base(name), name)
}

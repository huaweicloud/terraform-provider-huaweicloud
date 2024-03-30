package eg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/eg/v1/subscriptions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getEventStreamFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	httpUrl := "v1/{project_id}/eventstreamings/{eventstreaming_id}"
	client, err := cfg.NewServiceClient("eg", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{eventstreaming_id}", state.Primary.ID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving event stream: %s", err)
	}

	getConnectionRespBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Connection: %s", err)
	}

	return getConnectionRespBody, nil
}

func TestAccEventStream_basic(t *testing.T) {
	var (
		obj subscriptions.Subscription

		rName = "huaweicloud_eg_event_stream.test"
		name  = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(rName, &obj, getEventStreamFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgAgencyName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEventStream_kafkaSource_step1(name, "START"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "action", "START"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					// Source check
					resource.TestCheckResourceAttr(rName, "source.#", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.name", "HC.Kafka"),
					// Target check
					resource.TestCheckResourceAttr(rName, "sink.#", "1"),
					resource.TestCheckResourceAttr(rName, "sink.0.name", "HC.FunctionGraph"),
					// Rules check
					resource.TestCheckResourceAttr(rName, "rule_config.0.transform.0.type", "ORIGINAL"),
					// Option check
					resource.TestCheckResourceAttr(rName, "option.0.thread_num", "2"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.count", "5"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.time", "3"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.interval", "2"),
				),
			},
			{
				Config: testAccEventStream_kafkaSource_step2(name, "START"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "action", "START"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					// Source check
					resource.TestCheckResourceAttr(rName, "source.#", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.name", "HC.Kafka"),
					// Target check
					resource.TestCheckResourceAttr(rName, "sink.#", "1"),
					resource.TestCheckResourceAttr(rName, "sink.0.name", "HC.Kafka"),
					// Rules check
					resource.TestCheckResourceAttr(rName, "rule_config.0.transform.0.type", "ORIGINAL"),
					// Option check
					resource.TestCheckResourceAttr(rName, "option.0.thread_num", "3"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.count", "10"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.time", "5"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.interval", "3"),
				),
			},
			{
				Config: testAccEventStream_kafkaSource_step2(name, "PAUSE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "STOPPED"),
					resource.TestCheckResourceAttr(rName, "action", "PAUSE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source",
					"sink",
				},
			},
		},
	})
}

func testAccEventStream_kafkaSource_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type              = "cluster"
  arch_type         = "X86"
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 3)
}

// When EG communicates with the source instance (Kafka type), it will register 1 to 3 backend instances to the
// subnet and does not allow them to be deleted within one hour after the communication is disconnected.
// This results for the subnet means the creation through resources is no longer possible, it's released normally when
// the test case is completed, so the subnet can only be obtained through the data source.
data "huaweicloud_vpc_subnets" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "9011,9092-9095"
  remote_ip_prefix  = "0.0.0.0/0"
}

locals {
  flavor          = [for o in data.huaweicloud_dms_kafka_flavors.test.flavors : o if !strcontains(o.properties[0].flavor_alias, "beta")][0]
  connect_address = join(",", [for o in try(huaweicloud_dms_kafka_instance.test.cross_vpc_accesses, []): format("%%s:%%d",
                        o.advertised_ip, huaweicloud_dms_kafka_instance.test.port)])
}

resource "huaweicloud_dms_kafka_instance" "test" {
  name              = "%[1]s"
  vpc_id            = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, "")
  network_id        = try(data.huaweicloud_vpc_subnets.test.subnets[0].id, "")
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor_id          = try(local.flavor.id, "")
  storage_spec_code  = try(local.flavor.ios[0].storage_spec_code, "")
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 3)

  engine_version = "2.7"
  storage_space  = try(local.flavor.properties[0].min_broker, 0) * try(local.flavor.properties[0].min_storage_per_node, 0)
  broker_num     = 3

  security_protocol  = "PLAINTEXT"
  enabled_mechanisms = ["SCRAM-SHA-512"]
}

resource "huaweicloud_dms_kafka_topic" "test" {
  count = 2

  instance_id = huaweicloud_dms_kafka_instance.test.id
  name        = "%[1]s_${count.index}"
  partitions  = 1
}

variable "request_resp_print_script_content" {
  default = <<EOT
exports.handler = async (event, context) => {
    const result =
    {
        'repsonse_code': 200,
        'headers':
        {
            'Content-Type': 'application/json'
        },
        'isBase64Encoded': false,
        'body': JSON.stringify(event)
    }
    return result
}
EOT
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  agency      = "function_all_trust"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  code_type   = "inline"
  runtime     = "Node.js12.13"
  func_code   = base64encode(var.request_resp_print_script_content)
}

resource "huaweicloud_eg_connection" "test" {
  name      = "%[1]s"
  vpc_id    = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, "")
  subnet_id = try(data.huaweicloud_vpc_subnets.test.subnets[0].id, "")
  type      = "KAFKA"

  kafka_detail {
    instance_id     = huaweicloud_dms_kafka_instance.test.id
    connect_address = local.connect_address
  }
}
`, name)
}

func testAccEventStream_kafkaSource_step1(name, action string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_event_stream" "test" {
  name        = "%[2]s"
  description = "Created by acceptance test"
  action      = "%[3]s"

  source {
    name  = "HC.Kafka"
    kafka = jsonencode({
      "addr": local.connect_address,
      "group": "%[2]s",
      "instance_name": huaweicloud_dms_kafka_instance.test.name,
      "instance_id": huaweicloud_dms_kafka_instance.test.id,
      "topic": huaweicloud_dms_kafka_topic.test[0].name,
      "seek_to": "latest",
      "security_protocol": "PLAINTEXT",
    })
  }
  sink {
    name          = "HC.FunctionGraph"
    functiongraph = jsonencode({
      "urn": huaweicloud_fgs_function.test.urn,
      "agency": "%[4]s",
      "invoke_type": "ASYNC"
    })
  }
  rule_config {
    transform {
      type = "ORIGINAL"
    }
  }
  option {
    thread_num = 2

    batch_window {
      count    = 5
      time     = 3
      interval = 2
    }
  }
}
`, testAccEventStream_kafkaSource_base(name), name, action, acceptance.HW_EG_AGENCY_NAME)
}

func testAccEventStream_kafkaSource_step2(name, action string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_event_stream" "test" {
  name   = "%[2]s_update"
  action = "%[3]s"

  source {
    name  = "HC.Kafka"
    kafka = jsonencode({
      "addr": local.connect_address,
      "group": "%[2]s",
      "instance_name": huaweicloud_dms_kafka_instance.test.name,
      "instance_id": huaweicloud_dms_kafka_instance.test.id,
      "topic": huaweicloud_dms_kafka_topic.test[1].name,
      "seek_to": "latest",
      "security_protocol": "PLAINTEXT",
    })
  }
  sink {
    name  = "HC.Kafka"
    kafka = jsonencode({
      "connection_id": huaweicloud_eg_connection.test.id,
      "topic": huaweicloud_dms_kafka_topic.test[0].name,
      "key_transform": {
        "type": "ORIGINAL"
      }
    })
  }
  rule_config {
    transform {
      type = "ORIGINAL"
    }
  }
  option {
    thread_num = 3

    batch_window {
      count    = 10
      time     = 5
      interval = 3
    }
  }
}
`, testAccEventStream_kafkaSource_base(name), name, action)
}

func TestAccEventStream_rocketMQ(t *testing.T) {
	var (
		obj subscriptions.Subscription

		rName = "huaweicloud_eg_event_stream.test"
		name  = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(rName, &obj, getEventStreamFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEG(t)
			acceptance.TestAccPreCheckEgAgencyName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEventStream_rocketMQSource_step1(name, "START"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttr(rName, "status", "RUNNING"),
					resource.TestCheckResourceAttr(rName, "action", "START"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					// Source check
					resource.TestCheckResourceAttr(rName, "source.#", "1"),
					resource.TestCheckResourceAttr(rName, "source.0.name", "HC.DMS_ROCKETMQ"),
					// Target check
					resource.TestCheckResourceAttr(rName, "sink.#", "1"),
					resource.TestCheckResourceAttr(rName, "sink.0.name", "HC.FunctionGraph"),
					// Rules check
					resource.TestCheckResourceAttr(rName, "rule_config.0.transform.0.type", "ORIGINAL"),
					// Option check
					resource.TestCheckResourceAttr(rName, "option.0.thread_num", "2"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.count", "5"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.time", "3"),
					resource.TestCheckResourceAttr(rName, "option.0.batch_window.0.interval", "2"),
				),
			},
			{
				Config: testAccEventStream_rocketMQSource_step1(name, "PAUSE"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "STOPPED"),
					resource.TestCheckResourceAttr(rName, "action", "PAUSE"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source",
					"sink",
				},
			},
		},
	})
}

func testAccEventStream_rocketMQSource_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_rocketmq_flavors" "test" {
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 3)
  arch_type          = "X86"
  charging_mode      = "prePaid"
  type               = "cluster"
}

// When EG communicates with the source instance (RocketMQ type), it will register 1 to 3 backend instances to the
// subnet and does not allow them to be deleted within one hour after the communication is disconnected.
// This results for the subnet means the creation through resources is no longer possible, it's released normally when
// the test case is completed, so the subnet can only be obtained through the data source.
data "huaweicloud_vpc_subnets" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

locals {
  flavor = try(data.huaweicloud_dms_rocketmq_flavors.test.flavors[0], "")
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "8080,8081,8100,8200,10100-10199"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name               = "%[1]s"
  flavor_id          = try(local.flavor.id, "")
  engine_version     = "4.8.0"
  broker_num         = 1
  storage_space      = 300
  vpc_id             = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, "")
  subnet_id          = try(data.huaweicloud_vpc_subnets.test.subnets[0].id, "")
  security_group_id  = huaweicloud_networking_secgroup.test.id
  storage_spec_code  = try(local.flavor.ios[0].storage_spec_code, "")
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 3)
  enable_acl         = true
}

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = "%[1]s"

  brokers {
    name = "broker-0"
  }

  lifecycle {
    ignore_changes = [
      brokers,
    ]
  }
}

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id     = huaweicloud_dms_rocketmq_instance.test.id
  name            = "%[1]s"
  brokers         = ["broker-0"]
  retry_max_times = 3

  lifecycle {
    ignore_changes = [
      brokers, retry_max_times,
    ]
  }
}

variable "request_resp_print_script_content" {
  default = <<EOT
exports.handler = async (event, context) => {
    const result =
    {
        'repsonse_code': 200,
        'headers':
        {
            'Content-Type': 'application/json'
        },
        'isBase64Encoded': false,
        'body': JSON.stringify(event)
    }
    return result
}
EOT
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  agency      = "function_all_trust"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  code_type   = "inline"
  runtime     = "Node.js12.13"
  func_code   = base64encode(var.request_resp_print_script_content)
}
`, name)
}

func testAccEventStream_rocketMQSource_step1(name, action string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_event_stream" "test" {
  name        = "%[2]s"
  description = "Created by acceptance test"
  action      = "%[3]s"

  source {
    name         = "HC.DMS_ROCKETMQ"
    dms_rocketmq = jsonencode({
      "instance_id": huaweicloud_dms_rocketmq_instance.test.id,
      "group": huaweicloud_dms_rocketmq_consumer_group.test.name,
      "topic": huaweicloud_dms_rocketmq_topic.test.name,
      "tag": "lance",
      "access_key": "user_test",
      "secret_key": "Overlord!!52259",
      "ssl_enable": false,
      "enable_acl": true,
      "message_type": "NORMAL",
      "consume_timeout": 30000,
      "consumer_thread_nums": 20,
      "consumer_batch_max_size": 2
    })
  }
  sink {
    name    = "HC.FunctionGraph"
    functiongraph = jsonencode({
      "urn": huaweicloud_fgs_function.test.urn,
      "agency": "%[4]s",
      "invoke_type": "ASYNC"
    })
  }
  rule_config {
    transform {
      type = "ORIGINAL"
    }
  }
  option {
    thread_num = 2

    batch_window {
      count    = 5
      time     = 3
      interval = 2
    }
  }
}
`, testAccEventStream_rocketMQSource_base(name), name, action, acceptance.HW_EG_AGENCY_NAME)
}

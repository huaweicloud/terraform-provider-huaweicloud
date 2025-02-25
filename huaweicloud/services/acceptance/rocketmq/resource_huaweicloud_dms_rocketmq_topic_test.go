package rocketmq

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDmsRocketMQTopicResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getRocketmqTopic: query DMS rocketmq topic
	var (
		getRocketmqTopicHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}"
		getRocketmqTopicProduct = "dmsv2"
	)
	getRocketmqTopicClient, err := cfg.NewServiceClient(getRocketmqTopicProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DmsRocketMQTopic Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<topic>")
	}
	instanceID := parts[0]
	topic := parts[1]
	getRocketmqTopicPath := getRocketmqTopicClient.Endpoint + getRocketmqTopicHttpUrl
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{project_id}", getRocketmqTopicClient.ProjectID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{instance_id}", instanceID)
	getRocketmqTopicPath = strings.ReplaceAll(getRocketmqTopicPath, "{topic}", topic)

	getRocketmqTopicOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqTopicResp, err := getRocketmqTopicClient.Request("GET", getRocketmqTopicPath, &getRocketmqTopicOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQTopic: %s", err)
	}
	return utils.FlattenResponse(getRocketmqTopicResp)
}

func TestAccDmsRocketMQTopic_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dms_rocketmq_topic.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDmsRocketMQTopicResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQTopic_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "total_read_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "total_write_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "3"),
				),
			},
			{
				Config: testDmsRocketMQTopic_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "permission", "sub"),
					resource.TestCheckResourceAttr(rName, "total_read_queue_num", "4"),
					resource.TestCheckResourceAttr(rName, "total_write_queue_num", "5"),
					resource.TestCheckResourceAttr(rName, "brokers.0.read_queue_num", "4"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "5"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "queue_num"},
			},
		},
	})
}

func TestAccDmsRocketMQTopic_withQueues(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dms_rocketmq_topic.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDmsRocketMQTopicResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQTopic_with_queues(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "permission", "sub"),
					resource.TestCheckResourceAttr(rName, "total_read_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "total_write_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "3"),
					resource.TestCheckResourceAttr(rName, "brokers.0.write_queue_num", "3"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "queues"},
			},
		},
	})
}

func TestAccDmsRocketMQTopic_version5(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dms_rocketmq_topic.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDmsRocketMQTopicResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQTopic_verison5(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "message_type", "NORMAL"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

func testAccDmsRocketmqTopic_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "4.8.0"
  storage_space     = 600
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccDmsRocketmqTopic_verison5(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%s"
  engine_version    = "5.x"
  storage_space     = 200
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  flavor_id         = "rocketmq.b1.large.1"
  storage_spec_code = "dms.physical.storage.high.v2"
}
`, common.TestBaseNetwork(rName), rName)
}

func testDmsRocketMQTopic_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = "%s"
  queue_num   = 3

  brokers {
    name = "broker-0"
  }
}
`, testAccDmsRocketmqTopic_Base(name), name)
}

func testDmsRocketMQTopic_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id           = huaweicloud_dms_rocketmq_instance.test.id
  name                  = "%s"
  permission            = "sub"
  total_read_queue_num  = "4"
  total_write_queue_num = "5"
}
`, testAccDmsRocketmqTopic_Base(name), name)
}

func testDmsRocketMQTopic_with_queues(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = "%s"
  permission  = "sub"

  queues {
    broker    = "broker-0"
    queue_num = 3
  }
}
`, testAccDmsRocketmqTopic_Base(name), name)
}

func testDmsRocketMQTopic_verison5(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id  = huaweicloud_dms_rocketmq_instance.test.id
  name         = "%s"
  message_type = "NORMAL"
}
`, testAccDmsRocketmqTopic_verison5(name), name)
}

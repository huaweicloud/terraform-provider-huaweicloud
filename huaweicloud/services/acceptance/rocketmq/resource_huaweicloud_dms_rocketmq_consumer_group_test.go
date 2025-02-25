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

func getDmsRocketMQConsumerGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getRocketmqConsumerGroup: query DMS rocketmq consumer group
	var (
		getRocketmqConsumerGroupHttpUrl = "v2/{project_id}/instances/{instance_id}/groups/{group}"
		getRocketmqConsumerGroupProduct = "dmsv2"
	)
	getRocketmqConsumerGroupClient, err := cfg.NewServiceClient(getRocketmqConsumerGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and group from resource id
	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<consumerGroup>")
	}
	instanceID := parts[0]
	name := parts[1]
	getRocketmqConsumerGroupPath := getRocketmqConsumerGroupClient.Endpoint + getRocketmqConsumerGroupHttpUrl
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{project_id}",
		getRocketmqConsumerGroupClient.ProjectID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{instance_id}", instanceID)
	getRocketmqConsumerGroupPath = strings.ReplaceAll(getRocketmqConsumerGroupPath, "{group}", name)

	getRocketmqConsumerGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getRocketmqConsumerGroupResp, err := getRocketmqConsumerGroupClient.Request("GET", getRocketmqConsumerGroupPath,
		&getRocketmqConsumerGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DmsRocketMQConsumerGroup: %s", err)
	}
	return utils.FlattenResponse(getRocketmqConsumerGroupResp)
}

func TestAccDmsRocketMQConsumerGroup_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_consumer_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQConsumerGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQConsumerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "true"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", "add description."),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testDmsRocketMQConsumerGroup_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "false"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "5"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDmsRocketMQConsumerGroup_version5(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dms_rocketmq_consumer_group.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDmsRocketMQConsumerGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDmsRocketMQConsumerGroup_version5(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "broadcast", "true"),
					resource.TestCheckResourceAttr(resourceName, "retry_max_times", "3"),
					resource.TestCheckResourceAttr(resourceName, "description", "add description."),
					resource.TestCheckResourceAttr(resourceName, "consume_orderly", "true"),
				),
			},
			{
				Config: testDmsRocketMQConsumerGroup_version5(rName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "consume_orderly", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDmsRocketmqConsumerGroup_version4(rName string) string {
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
  broker_num        = 2
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccDmsRocketmqConsumerGroup_version5(rName string) string {
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

func testDmsRocketMQConsumerGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  broadcast   = true

  brokers = [
    "broker-0",
    "broker-1"
  ]

  name            = "%s"
  retry_max_times = "3"
  description     = "add description."
}
`, testAccDmsRocketmqConsumerGroup_version4(name), name)
}

func testDmsRocketMQConsumerGroup_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  broadcast   = false

  brokers = [
    "broker-0",
    "broker-1"
  ]

  name            = "%s"
  retry_max_times = "5"
  enabled         = false
  description     = ""
}
`, testAccDmsRocketmqConsumerGroup_version4(name), name)
}

func testDmsRocketMQConsumerGroup_version5(name string, orderly bool) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id     = huaweicloud_dms_rocketmq_instance.test.id
  broadcast       = true
  name            = "%s"
  retry_max_times = "3"
  description     = "add description."
  consume_orderly = %v
}
`, testAccDmsRocketmqConsumerGroup_version5(name), name, orderly)
}

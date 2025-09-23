package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/eg/v1/source/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getCustomEventSourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.EgV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG v1 client: %s", err)
	}

	return custom.Get(client, state.Primary.ID)
}

func TestAccCustomEventSource_basic(t *testing.T) {
	var (
		obj custom.Source

		rName = "huaweicloud_eg_custom_event_source.test"
		name  = acceptance.RandomAccResourceName()
		rc    = acceptance.InitResourceCheck(rName, &obj, getCustomEventSourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgChannelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventSource_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "channel_id", acceptance.HW_EG_CHANNEL_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "APPLICATION"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccCustomEventSource_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "channel_id", acceptance.HW_EG_CHANNEL_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "APPLICATION"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCustomEventSource_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id  = "%[1]s"
  name        = "%[2]s"
  description = "Created by acceptance test"
}
`, acceptance.HW_EG_CHANNEL_ID, name)
}

func testAccCustomEventSource_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = "%[1]s"
  name       = "%[2]s"
}
`, acceptance.HW_EG_CHANNEL_ID, name)
}

func TestAccCustomEventSource_rocketMQ(t *testing.T) {
	var (
		obj custom.Source

		rName = "huaweicloud_eg_custom_event_source.test"
		name  = acceptance.RandomAccResourceName()
		rc    = acceptance.InitResourceCheck(rName, &obj, getCustomEventSourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventSource_rocketMQ_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "channel_id", "huaweicloud_eg_custom_event_channel.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ROCKETMQ"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccCustomEventSource_rocketMQ_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "channel_id", "huaweicloud_eg_custom_event_channel.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "ROCKETMQ"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCustomEventSource_rocketMQ_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

# The corresponding elastic network card cannot be deleted within one hour after the EG resources are deleted.
data "huaweicloud_vpc_subnets" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "with_icmp_ingress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "with_rocketmq_ingress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "8100,10100-10103"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "with_rocketmq_egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_dms_rocketmq_instance" "test" {
  name              = "%[1]s"
  engine_version    = "4.8.0"
  storage_space     = 300
  vpc_id            = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, "")
  subnet_id         = try(data.huaweicloud_vpc_subnets.test.subnets[0].id, "")
  security_group_id = huaweicloud_networking_secgroup.test.id

  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 2)

  flavor_id         = "c6.4u8g.cluster"
  storage_spec_code = "dms.physical.storage.high.v2"
  broker_num        = 1
}

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  count = 2

  instance_id     = huaweicloud_dms_rocketmq_instance.test.id
  name            = format("%[1]s-%%d", count.index)
  enabled         = true
  broadcast       = true
  brokers         = ["broker-0"]
  retry_max_times = 3
}

resource "huaweicloud_dms_rocketmq_topic" "test" {
  count = 2

  instance_id = huaweicloud_dms_rocketmq_instance.test.id
  name        = format("%[1]s-%%d", count.index)
  queue_num   = 3
  permission  = "all"

  brokers {
    name = "broker-0"
  }
}

resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

resource "huaweicloud_eg_endpoint" "test" {
  name      = "%[1]s"
  vpc_id    = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, "")
  subnet_id = try(data.huaweicloud_vpc_subnets.test.subnets[0].id, "")
}
`, name)
}

func testAccCustomEventSource_rocketMQ_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = "%[2]s"
  type       = "ROCKETMQ"
  detail     = jsonencode({
    instance_id     = huaweicloud_dms_rocketmq_instance.test.id
    group           = try(split("/", huaweicloud_dms_rocketmq_consumer_group.test[0].id)[0], "")
    topic           = try(split("/", huaweicloud_dms_rocketmq_topic.test[0].id)[0], "")
    enable_acl      = false
    name            = "%[2]s"
    namesrv_address = huaweicloud_dms_rocketmq_instance.test.namesrv_address
    ssl_enable      = false
  })
}
`, testAccCustomEventSource_rocketMQ_base(name), name)
}

func testAccCustomEventSource_rocketMQ_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
  name       = "%[2]s"
  type       = "ROCKETMQ"
  detail     = jsonencode({
    instance_id     = huaweicloud_dms_rocketmq_instance.test.id
    group           = try(split("/", huaweicloud_dms_rocketmq_consumer_group.test[1].id)[0], "")
    topic           = try(split("/", huaweicloud_dms_rocketmq_topic.test[1].id)[0], "")
    enable_acl      = false
    name            = "%[2]s"
    namesrv_address = huaweicloud_dms_rocketmq_instance.test.namesrv_address
    ssl_enable      = false
  })
}
`, testAccCustomEventSource_rocketMQ_base(name), name)
}

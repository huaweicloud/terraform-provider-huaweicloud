package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/apig"
)

func getChannelMemberFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		instanceId   = state.Primary.Attributes["instance_id"]
		vpcChannelId = state.Primary.Attributes["vpc_channel_id"]
		memberId     = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("apig", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating APIG client: %s", err)
	}

	return apig.GetChannelMemberById(client, instanceId, vpcChannelId, memberId)
}

func TestAccChannelMember_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		ipMember     interface{}
		ipMemberName = "huaweicloud_apig_channel_member.ip_member"
		rcIpMember   = acceptance.InitResourceCheck(ipMemberName, &ipMember, getChannelMemberFunc)

		ecsMember     interface{}
		ecsMemberName = "huaweicloud_apig_channel_member.ecs_member"
		rcEcsMember   = acceptance.InitResourceCheck(ecsMemberName, &ecsMember, getChannelMemberFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSubnetId(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcIpMember.CheckResourceDestroy(),
			rcEcsMember.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccChannelMember_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rcIpMember.CheckResourceExists(),
					resource.TestCheckResourceAttr(ipMemberName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrSet(ipMemberName, "vpc_channel_id"),
					resource.TestCheckResourceAttr(ipMemberName, "weight", "20"),
					resource.TestCheckResourceAttr(ipMemberName, "port", "80"),
					resource.TestCheckResourceAttr(ipMemberName, "is_backup", "false"),
					resource.TestCheckResourceAttr(ipMemberName, "status", "1"),
					resource.TestCheckResourceAttrSet(ipMemberName, "ecs_id"),
					resource.TestCheckResourceAttrSet(ipMemberName, "ecs_name"),
					resource.TestMatchResourceAttr(ipMemberName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(ipMemberName, "member_group_id"),
					rcEcsMember.CheckResourceExists(),
					resource.TestCheckResourceAttr(ecsMemberName, "instance_id", acceptance.HW_APIG_DEDICATED_INSTANCE_ID),
					resource.TestCheckResourceAttrSet(ecsMemberName, "vpc_channel_id"),
					resource.TestCheckResourceAttr(ecsMemberName, "weight", "20"),
					resource.TestCheckResourceAttr(ecsMemberName, "port", "80"),
					resource.TestCheckResourceAttr(ecsMemberName, "is_backup", "false"),
					resource.TestCheckResourceAttr(ecsMemberName, "status", "1"),
					resource.TestCheckResourceAttrSet(ecsMemberName, "ecs_id"),
					resource.TestCheckResourceAttrSet(ecsMemberName, "ecs_name"),
					resource.TestMatchResourceAttr(ecsMemberName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(ecsMemberName, "member_group_id"),
				),
			},
			{
				ResourceName:      ipMemberName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelMemberImportStateFunc(ipMemberName),
			},
			{
				ResourceName:      ecsMemberName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccChannelMemberImportStateFunc(ecsMemberName),
			},
		},
	})
}

func testAccChannelMemberImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("Resource (%s) not found: %s", rsName, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" || rs.Primary.Attributes["vpc_channel_id"] == "" || rs.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", rs.Primary.Attributes["instance_id"],
				rs.Primary.Attributes["vpc_channel_id"], rs.Primary.ID)
		}
		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["vpc_channel_id"], rs.Primary.ID), nil
	}
}

func testAccChannelMember_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

resource "huaweicloud_apig_channel" "ip_channel" {
  instance_id      = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name             = format("%[1]s_%%s", "ip")
  port             = 80
  balance_strategy = 1
  member_type      = "ip"
  type             = "builtin"
}

resource "huaweicloud_apig_channel" "ecs_channel" {
  instance_id      = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name             = format("%[1]s_%%s", "ecs")
  port             = 80
  balance_strategy = 1
  member_type      = "ecs"
  type             = "builtin"
}

resource "huaweicloud_apig_channel_member_group" "ip_channel_member_group" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ip_channel.id
  name           = format("%[1]s_%%s", "ip")
  weight         = 20
  description    = "ip channel member group."
}

resource "huaweicloud_apig_channel_member_group" "ecs_channel_member_group" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ecs_channel.id
  name           = format("%[1]s_%%s", "ecs")
  weight         = 20
  description    = "ecs channel member group."
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  id = "%[3]s"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_instance" "test" {
  name                = "%[1]s"
  description         = "terraform test"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]
  stop_before_destroy = true
  agency_name         = "test111"
  agent_list          = "hss"

  user_data = <<EOF
#! /bin/bash
echo user_test > /home/user.txt
EOF

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }
}
`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, acceptance.HW_SUBNET_ID)
}

func testAccChannelMember_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_apig_channel_member" "ip_member" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ip_channel.id
  weight            = 20
  port              = 80
  member_group_name = huaweicloud_apig_channel_member_group.ip_channel_member_group.name
  member_ip_address = huaweicloud_compute_instance.test.access_ip_v4
}

resource "huaweicloud_apig_channel_member" "ecs_member" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ecs_channel.id
  weight            = 20
  port              = 80
  member_group_name = huaweicloud_apig_channel_member_group.ecs_channel_member_group.name
  ecs_id            = huaweicloud_compute_instance.test.id
  ecs_name          = huaweicloud_compute_instance.test.name
}
`, testAccChannelMember_base(name))
}

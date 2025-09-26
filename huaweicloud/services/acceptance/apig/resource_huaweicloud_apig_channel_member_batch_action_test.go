package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApigChannelMemberBatchAction_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		ipMemberDisableName = "data.huaweicloud_apig_channel_members.ip_after_disable_operation"
		dcIpMemberDisable   = acceptance.InitDataSourceCheck(ipMemberDisableName)

		ecsMemberDisableName = "data.huaweicloud_apig_channel_members.ecs_after_disable_operation"
		dcEcsMemberDisable   = acceptance.InitDataSourceCheck(ecsMemberDisableName)

		ipMemberEnableName = "data.huaweicloud_apig_channel_members.ip_after_enable_operation"
		dcIpMemberEnable   = acceptance.InitDataSourceCheck(ipMemberEnableName)

		ecsMemberEnableName = "data.huaweicloud_apig_channel_members.ecs_after_enable_operation"
		dcEcsMemberEnable   = acceptance.InitDataSourceCheck(ecsMemberEnableName)
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelMemberBatchAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					dcIpMemberDisable.CheckResourceExists(),
					resource.TestCheckOutput("ip_disable_is_useful", "true"),
					dcEcsMemberDisable.CheckResourceExists(),
					resource.TestCheckOutput("ecs_disable_is_useful", "true"),
				),
			},
			{
				Config: testAccChannelMemberBatchAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dcIpMemberEnable.CheckResourceExists(),
					resource.TestCheckOutput("ip_enable_is_useful", "true"),
					dcEcsMemberEnable.CheckResourceExists(),
					resource.TestCheckOutput("ecs_enable_is_useful", "true"),
				),
			},
		},
	})
}

func testAccChannelMemberBatchAction_base(name string) string {
	return fmt.Sprintf(`
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
  count = 2

  name                = format("%[1]s_%%s", count.index)
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

resource "huaweicloud_apig_channel_member" "ip_member" {
  count = 2

  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ip_channel.id
  member_group_name = huaweicloud_apig_channel_member_group.ip_channel_member_group.name
  member_ip_address = huaweicloud_compute_instance.test[count.index].access_ip_v4
  port              = 80
}

resource "huaweicloud_apig_channel_member" "ecs_member" {
  count = 2

  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ecs_channel.id
  member_group_name = huaweicloud_apig_channel_member_group.ecs_channel_member_group.name
  ecs_id            = huaweicloud_compute_instance.test[count.index].id
  ecs_name          = huaweicloud_compute_instance.test[count.index].name
  port              = 80
}`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccChannelMemberBatchAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# Disable channel members
resource "huaweicloud_apig_channel_member_batch_action" "disable_ip_members" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ip_channel.id
  action         = "disable"
  
  member_ids = [
    huaweicloud_apig_channel_member.ip_member[0].id,
    huaweicloud_apig_channel_member.ip_member[1].id,
  ]
}

resource "huaweicloud_apig_channel_member_batch_action" "disable_ecs_members" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ecs_channel.id
  action         = "disable"

  member_ids = [
    huaweicloud_apig_channel_member.ecs_member[0].id,
    huaweicloud_apig_channel_member.ecs_member[1].id,
  ]
}

# Whether channel members are disable
locals {
  group_name_prefix = "%[2]s"
}

data "huaweicloud_apig_channel_members" "ip_after_disable_operation" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ip_channel.id
  member_group_name = local.group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member_batch_action.disable_ip_members,
  ]
}

locals {
  ip_disable_result = [
    for v in data.huaweicloud_apig_channel_members.ip_after_disable_operation.members[*].status
    : v == 2
  ]
}

output "ip_disable_is_useful" {
  value = length(local.ip_disable_result) >= 2 && alltrue(local.ip_disable_result)
}

data "huaweicloud_apig_channel_members" "ecs_after_disable_operation" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ecs_channel.id
  member_group_name = local.group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member_batch_action.disable_ecs_members,
  ]
}

locals {
  ecs_disable_result = [
    for v in data.huaweicloud_apig_channel_members.ecs_after_disable_operation.members[*].status
    : v == 2
  ]
}

output "ecs_disable_is_useful" {
  value = length(local.ecs_disable_result) >= 2 && alltrue(local.ecs_disable_result)
}
`, testAccChannelMemberBatchAction_base(name), name)
}

func testAccChannelMemberBatchAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

# Enable channel members
resource "huaweicloud_apig_channel_member_batch_action" "enable_ip_members" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ip_channel.id
  action         = "enable"

  member_ids = [
    huaweicloud_apig_channel_member.ip_member[0].id,
    huaweicloud_apig_channel_member.ip_member[1].id,
  ]

  depends_on = [
    data.huaweicloud_apig_channel_members.ip_after_disable_operation,
  ]
}

resource "huaweicloud_apig_channel_member_batch_action" "enable_ecs_members" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ecs_channel.id
  action         = "enable"

  member_ids = [
    huaweicloud_apig_channel_member.ecs_member[0].id,
    huaweicloud_apig_channel_member.ecs_member[1].id,
  ]

  depends_on = [
    data.huaweicloud_apig_channel_members.ecs_after_disable_operation,
  ]
}

# Whether channel members are enable
data "huaweicloud_apig_channel_members" "ip_after_enable_operation" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ip_channel.id
  member_group_name = local.group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member_batch_action.enable_ip_members,
  ]
}

locals {
  ip_enable_result = [
    for v in data.huaweicloud_apig_channel_members.ip_after_enable_operation.members[*].status
    : v == 1
  ]
}

output "ip_enable_is_useful" {
  value = length(local.ip_enable_result) >= 2 && alltrue(local.ip_enable_result)
}

data "huaweicloud_apig_channel_members" "ecs_after_enable_operation" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ecs_channel.id
  member_group_name = local.group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member_batch_action.enable_ecs_members,
  ]
}

locals {
  ecs_enable_result = [
    for v in data.huaweicloud_apig_channel_members.ecs_after_enable_operation.members[*].status
    : v == 1
  ]
}

output "ecs_enable_is_useful" {
  value = length(local.ecs_enable_result) >= 2 && alltrue(local.ecs_enable_result)
}
`, testAccChannelMemberBatchAction_basic_step1(name), name)
}

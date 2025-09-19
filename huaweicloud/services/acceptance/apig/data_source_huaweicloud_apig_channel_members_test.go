package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceChannelMembers_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_apig_channel_members.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_apig_channel_members.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byMemberGroupName   = "data.huaweicloud_apig_channel_members.filter_by_member_group_name"
		dcByMemberGroupName = acceptance.InitDataSourceCheck(byMemberGroupName)

		byMemberGroupId   = "data.huaweicloud_apig_channel_members.filter_by_member_group_id"
		dcByMemberGroupId = acceptance.InitDataSourceCheck(byMemberGroupId)

		byPreciseSearch   = "data.huaweicloud_apig_channel_members.filter_by_precise_search"
		dcByPreciseSearch = acceptance.InitDataSourceCheck(byPreciseSearch)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceChannelMembers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "members.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.weight"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.is_backup"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.vpc_channel_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.member_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.member_group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.member_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.ecs_id"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.ecs_name"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "members.0.health_status"),
					resource.TestMatchResourceAttr(dataSource, "members.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByMemberGroupName.CheckResourceExists(),
					resource.TestCheckOutput("member_group_name_filter_is_useful", "true"),
					dcByMemberGroupId.CheckResourceExists(),
					resource.TestCheckOutput("member_group_id_filter_is_useful", "true"),
					dcByPreciseSearch.CheckResourceExists(),
					resource.TestCheckOutput("precise_search_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceChannelMembers_base(name string) string {
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
  count = 2

  name                = format("%[1]s_%%d", count.index)
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

resource "huaweicloud_apig_channel_member" "ip_member" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ip_channel.id
  weight            = 20
  port              = 80
  member_group_name = huaweicloud_apig_channel_member_group.ip_channel_member_group.name
  member_ip_address = huaweicloud_compute_instance.test[0].access_ip_v4
}

resource "huaweicloud_apig_channel_member" "ecs_member" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ecs_channel.id
  weight            = 20
  port              = 80
  member_group_name = huaweicloud_apig_channel_member_group.ecs_channel_member_group.name
  ecs_id            = huaweicloud_compute_instance.test[1].id
  ecs_name          = huaweicloud_compute_instance.test[1].name
}
`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, acceptance.HW_SUBNET_ID)
}

func testAccDataSourceChannelMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all
data "huaweicloud_apig_channel_members" "test" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
	huaweicloud_apig_channel_member.ecs_member,
  ]
}

# Filter by name (fuzzy search)
locals {
  member_name_prefix = "%[2]s"
}

data "huaweicloud_apig_channel_members" "filter_by_name" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id
  name           = local.member_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
	huaweicloud_apig_channel_member.ecs_member,
  ]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_channel_members.filter_by_name.members[*].name
	: strcontains(v, local.member_name_prefix)
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) >= 1 && alltrue(local.name_filter_result)
}

# Filter by group name (fuzzy search)
locals {
  member_group_name_prefix = "%[2]s"
}

data "huaweicloud_apig_channel_members" "filter_by_member_group_name" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.test.id
  member_group_name = local.member_group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
	huaweicloud_apig_channel_member.ecs_member,
  ]
}

locals {
  member_group_name_filter_result = [
    for v in data.huaweicloud_apig_channel_members.filter_by_member_group_name.members[*].member_group_name
	: strcontains(v, local.member_group_name_prefix)
  ]
}

output "member_group_name_filter_is_useful" {
  value = length(local.member_group_name_filter_result) >= 2 && alltrue(local.member_group_name_filter_result)
}

# Filter by group id
locals {
  member_group_id = huaweicloud_apig_channel_member_group.ip_channel_member_group.id
}

data "huaweicloud_apig_channel_members" "filter_by_member_group_id" {
  instance_id     = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id  = huaweicloud_apig_channel.test.id
  member_group_id = local.member_group_id

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
	huaweicloud_apig_channel_member.ecs_member,
  ]
}

locals {
  member_group_id_filter_result = [
    for v in data.huaweicloud_apig_channel_members.filter_by_member_group_id.members[*].member_group_id
	: v == local.member_group_id
  ]
}

output "member_group_id_filter_is_useful" {
  value = length(local.member_group_id_filter_result) >= 1 && alltrue(local.member_group_id_filter_result)
}

# Filter by group name (exact search)
locals {
  member_group_name = huaweicloud_apig_channel_member_group.ip_channel_member_group.name
}

data "huaweicloud_apig_channel_members" "filter_by_precise_search" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.test.id
  member_group_name = local.member_group_name
  precise_search    = "member_group_name"

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
	huaweicloud_apig_channel_member.ecs_member,
  ]
}

locals {
  precise_search_filter_result = [
    for v in data.huaweicloud_apig_channel_members.filter_by_precise_search.members[*].member_group_name
	: v == local.member_group_name
  ]
}

output "precise_search_filter_is_useful" {
  value = length(local.precise_search_filter_result) >= 1 && alltrue(local.precise_search_filter_result)
}
`, testAccDataSourceChannelMembers_base(name), name)
}

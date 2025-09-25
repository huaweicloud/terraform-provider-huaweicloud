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

		ipAllName   = "data.huaweicloud_apig_channel_members.ip_all"
		dcIpAllName = acceptance.InitDataSourceCheck(ipAllName)

		ecsAllName   = "data.huaweicloud_apig_channel_members.ecs_all"
		dcEcsAllName = acceptance.InitDataSourceCheck(ecsAllName)

		ecsFilterByName   = "data.huaweicloud_apig_channel_members.ecs_filter_by_name"
		dcEcsFilterByName = acceptance.InitDataSourceCheck(ecsFilterByName)

		ecsFilterByGroupName   = "data.huaweicloud_apig_channel_members.ecs_filter_by_group_name"
		dcEcsFilterByGroupName = acceptance.InitDataSourceCheck(ecsFilterByGroupName)

		ipFilterByGroupName   = "data.huaweicloud_apig_channel_members.ip_filter_by_group_name"
		dcIpFilterByGroupName = acceptance.InitDataSourceCheck(ipFilterByGroupName)

		ecsFilterByGroupId   = "data.huaweicloud_apig_channel_members.ecs_filter_by_group_id"
		dcEcsFilterByGroupId = acceptance.InitDataSourceCheck(ecsFilterByGroupId)

		ipFilterByGroupId   = "data.huaweicloud_apig_channel_members.ip_filter_by_group_id"
		dcIpFilterByGroupId = acceptance.InitDataSourceCheck(ipFilterByGroupId)

		filterByPreciseSearch   = "data.huaweicloud_apig_channel_members.filter_by_precise_search"
		dcFilterByPreciseSearch = acceptance.InitDataSourceCheck(filterByPreciseSearch)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceChannelMembers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcIpAllName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.id"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.vpc_channel_id"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.member_group_name"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.member_group_id"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.member_ip_address"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.ecs_id"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.ecs_name"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.port"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.is_backup"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.status"),
					resource.TestCheckResourceAttrSet(ipAllName, "members.0.weight"),
					resource.TestMatchResourceAttr(ipAllName, "members.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcEcsAllName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.id"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.vpc_channel_id"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.member_group_name"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.member_group_id"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.member_ip_address"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.ecs_id"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.ecs_name"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.port"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.is_backup"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.status"),
					resource.TestCheckResourceAttrSet(ecsAllName, "members.0.weight"),
					resource.TestMatchResourceAttr(ecsAllName, "members.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcEcsAllName.CheckResourceExists(),
					dcEcsFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("ecs_name_filter_is_useful", "true"),
					dcEcsFilterByGroupName.CheckResourceExists(),
					resource.TestCheckOutput("ecs_group_name_filter_is_useful", "true"),
					dcIpFilterByGroupName.CheckResourceExists(),
					resource.TestCheckOutput("ip_group_name_filter_is_useful", "true"),
					dcEcsFilterByGroupId.CheckResourceExists(),
					resource.TestCheckOutput("ecs_group_id_filter_is_useful", "true"),
					dcIpFilterByGroupId.CheckResourceExists(),
					resource.TestCheckOutput("ip_group_id_filter_is_useful", "true"),
					dcFilterByPreciseSearch.CheckResourceExists(),
					resource.TestCheckOutput("precise_search_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceChannelMembers_Compute_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  id = "%[2]s"
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
}`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID)
}

func testAccDataSourceChannelMembers_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

resource "huaweicloud_apig_channel" "ip_channel" {
  instance_id      = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name             = format("%[3]s_%%s", "ip")
  port             = 80
  balance_strategy = 1
  member_type      = "ip"
  type             = "builtin"
}

resource "huaweicloud_apig_channel" "ecs_channel" {
  instance_id      = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name             = format("%[3]s_%%s", "ecs")
  port             = 80
  balance_strategy = 1
  member_type      = "ecs"
  type             = "builtin"
}

resource "huaweicloud_apig_channel_member_group" "ip_channel_member_group" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ip_channel.id
  name           = format("%[3]s_%%s", "ip")
  weight         = 20
  description    = "ip channel member group."
}

resource "huaweicloud_apig_channel_member_group" "ecs_channel_member_group" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ecs_channel.id
  name           = format("%[3]s_%%s", "ecs")
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
  status            = 2
  is_backup         = true
}

resource "huaweicloud_apig_channel_member" "ecs_member" {
  count = 2

  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ecs_channel.id
  member_group_name = huaweicloud_apig_channel_member_group.ecs_channel_member_group.name
  ecs_id            = huaweicloud_compute_instance.test[count.index].id
  ecs_name          = huaweicloud_compute_instance.test[count.index].name
  port              = 80
}
`, testAccDataSourceChannelMembers_Compute_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceChannelMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all
data "huaweicloud_apig_channel_members" "ip_all" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ip_channel.id

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
  ]
}

data "huaweicloud_apig_channel_members" "ecs_all" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ecs_channel.id

  depends_on = [
	  huaweicloud_apig_channel_member.ecs_member,
  ]
}

# Filter by name (fuzzy search)
locals {
  member_name_prefix = "%[2]s"
}

data "huaweicloud_apig_channel_members" "ecs_filter_by_name" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.ecs_channel.id
  name           = local.member_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member.ecs_member,
  ]
}

locals {
  ecs_name_filter_result = [
    for v in data.huaweicloud_apig_channel_members.ecs_filter_by_name.members[*].ecs_name
    : strcontains(v, local.member_name_prefix)
  ]
}

output "ecs_name_filter_is_useful" {
  value = length(local.ecs_name_filter_result) >= 2 && alltrue(local.ecs_name_filter_result)
}

# Filter by group name (fuzzy search)
locals {
  group_name_prefix = "%[2]s"
}

data "huaweicloud_apig_channel_members" "ecs_filter_by_group_name" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ecs_channel.id
  member_group_name = local.group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member.ecs_member,
  ]
}

locals {
  ecs_group_name_filter_result = [
    for v in data.huaweicloud_apig_channel_members.ecs_filter_by_group_name.members[*].member_group_name 
    : strcontains(v, local.group_name_prefix)
  ]
}

output "ecs_group_name_filter_is_useful" {
  value = length(local.ecs_group_name_filter_result) >= 2 && alltrue(local.ecs_group_name_filter_result)
}

data "huaweicloud_apig_channel_members" "ip_filter_by_group_name" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ip_channel.id
  member_group_name = local.group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
  ]
}

locals {
  ip_group_name_filter_result = [
    for v in data.huaweicloud_apig_channel_members.ip_filter_by_group_name.members[*].member_group_name 
    : strcontains(v, local.group_name_prefix)
  ]
}

output "ip_group_name_filter_is_useful" {
  value = length(local.ip_group_name_filter_result) >= 2 && alltrue(local.ip_group_name_filter_result)
}

# Filter by group ID
locals {
  ip_group_id  = huaweicloud_apig_channel_member_group.ip_channel_member_group.id
  ecs_group_id = huaweicloud_apig_channel_member_group.ecs_channel_member_group.id
}

data "huaweicloud_apig_channel_members" "ecs_filter_by_group_id" {
  instance_id     = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id  = huaweicloud_apig_channel.ecs_channel.id
  member_group_id = local.ecs_group_id

  depends_on = [
    huaweicloud_apig_channel_member.ecs_member,
  ]
}

locals {
  ecs_group_id_filter_result = [
    for v in data.huaweicloud_apig_channel_members.ecs_filter_by_group_id.members[*].member_group_id
    : v == local.ecs_group_id
  ]
}

output "ecs_group_id_filter_is_useful" {
  value = length(local.ecs_group_id_filter_result) >= 2 && alltrue(local.ecs_group_id_filter_result)
}

data "huaweicloud_apig_channel_members" "ip_filter_by_group_id" {
  instance_id     = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id  = huaweicloud_apig_channel.ip_channel.id
  member_group_id = local.ip_group_id

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
  ]
}

locals {
  ip_group_id_filter_result = [
    for v in data.huaweicloud_apig_channel_members.ip_filter_by_group_id.members[*].member_group_id
    : v == local.ip_group_id
  ]
}

output "ip_group_id_filter_is_useful" {
  value = length(local.ip_group_id_filter_result) >= 2 && alltrue(local.ip_group_id_filter_result)
}

# Filter by group name (exact search)
locals {
  exact_group_name = huaweicloud_apig_channel_member_group.ip_channel_member_group.name
}

data "huaweicloud_apig_channel_members" "filter_by_precise_search" {
  instance_id       = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id    = huaweicloud_apig_channel.ip_channel.id
  member_group_name = local.exact_group_name
  precise_search    = "member_group_name"

  depends_on = [
    huaweicloud_apig_channel_member.ip_member,
  ]
}

locals {
  precise_search_filter_result = [
    for v in data.huaweicloud_apig_channel_members.filter_by_precise_search.members[*].member_group_name
    : v == local.exact_group_name
  ]
}

output "precise_search_filter_is_useful" {
  value = length(local.precise_search_filter_result) > 1 && alltrue(local.precise_search_filter_result)
}
`, testAccDataSourceChannelMembers_base(name), name)
}

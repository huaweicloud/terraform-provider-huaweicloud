package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceChannels_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_channels.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		rName      = acceptance.RandomAccResourceName()

		byId   = "data.huaweicloud_apig_channels.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_channels.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byExactName   = "data.huaweicloud_apig_channels.filter_by_exact_name"
		dcByExactName = acceptance.InitDataSourceCheck(byExactName)

		byNotFoundName   = "data.huaweicloud_apig_channels.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byMemberGroupId   = "data.huaweicloud_apig_channels.filter_by_member_group_id"
		dcByMemberGroupId = acceptance.InitDataSourceCheck(byMemberGroupId)

		byMemberGroupName   = "data.huaweicloud_apig_channels.filter_by_member_group_name"
		dcByMemberGroupName = acceptance.InitDataSourceCheck(byMemberGroupName)
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
				Config: testAccDataSourceChannels_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "vpc_channels.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttr(byId, "vpc_channels.#", "1"),
					resource.TestCheckResourceAttrSet(byId, "vpc_channels.0.id"),
					resource.TestCheckResourceAttrSet(byId, "vpc_channels.0.name"),
					resource.TestCheckResourceAttr(byId, "vpc_channels.0.member_group.#", "1"),
					resource.TestCheckResourceAttrSet(byId, "vpc_channels.0.member_group.0.id"),
					resource.TestCheckResourceAttrSet(byId, "vpc_channels.0.member_group.0.name"),
					resource.TestMatchResourceAttr(byId, "vpc_channels.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("channel_id_filter_is_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByExactName.CheckResourceExists(),
					resource.TestCheckOutput("exact_name_filter_is_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("not_found_name_filter_is_useful", "true"),
					dcByMemberGroupId.CheckResourceExists(),
					resource.TestCheckOutput("member_group_id_filter_is_useful", "true"),
					dcByMemberGroupName.CheckResourceExists(),
					resource.TestCheckOutput("member_group_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccChannel_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  count = 1

  name               = format("%[2]s-%%d", count.index)
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  system_disk_type   = "SSD"

  network {
    uuid = "%[3]s"
  }
}

resource "huaweicloud_apig_channel" "test" {
  instance_id        = local.instance_id
  name               = "%[2]s"
  port               = 80
  balance_strategy   = 1
  member_type        = "ecs"
  type               = 2

  health_check {
    protocol           = "TCP"
    threshold_normal   = 1 # minimum value
    threshold_abnormal = 1 # minimum value
    interval           = 1 # minimum value
    timeout            = 1 # minimum value
  }

  member_group {
    name        ="%[4]s"
    description = "Created by script"
    weight      = 1
  }

  dynamic "member" {
    for_each = huaweicloud_compute_instance.test[*]

    content {
      id   = member.value.id
      name = member.value.name
    }
  }
}
`, testAccChannel_base(rName), rName, acceptance.HW_APIG_DEDICATED_INSTANCE_USED_SUBNET_ID,
		acceptance.RandomAccResourceName())
}

func testAccDataSourceChannels_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_channels" "test" {
  depends_on = [
    huaweicloud_apig_channel.test
  ]

  instance_id = local.instance_id
}

# Filter by ID
locals {
  channel_id = huaweicloud_apig_channel.test.id
}

data "huaweicloud_apig_channels" "filter_by_id" {
  depends_on = [
    huaweicloud_apig_channel.test
  ]

  instance_id = local.instance_id
  channel_id  = local.channel_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_channels.filter_by_id.vpc_channels[*].id : v == local.channel_id
  ]
}

output "channel_id_filter_is_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name (fuzzy search)
locals {
  name = huaweicloud_apig_channel.test.name
}

data "huaweicloud_apig_channels" "filter_by_name" {
  depends_on = [
    huaweicloud_apig_channel.test
  ]

  instance_id = local.instance_id
  name        = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_channels.filter_by_name.vpc_channels[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (exact search)
data "huaweicloud_apig_channels" "filter_by_exact_name" {
  depends_on = [
    huaweicloud_apig_channel.test
  ]

  instance_id    = local.instance_id
  name           = local.name
  precise_search = "name,member_group_name"
}

output "exact_name_filter_is_useful" {
  value = length(data.huaweicloud_apig_channels.filter_by_exact_name.vpc_channels) == 1
}

# Filter by name (not found)
locals {
  not_found_name = "not_found"
}

data "huaweicloud_apig_channels" "filter_by_not_found_name" {
  depends_on = [
    huaweicloud_apig_channel.test
  ]

  instance_id    = local.instance_id
  name           = local.not_found_name
  precise_search = "name"
}

output "not_found_name_filter_is_useful" {
  value = length(data.huaweicloud_apig_channels.filter_by_not_found_name.vpc_channels) == 0
}

# Filter by member_group_id
# The member_group ID does not exist in the resource. Obtaining the member_group ID from the resource ID filtering result
# to query can ensure the correct filtering results. 
locals {
  member_group_id = data.huaweicloud_apig_channels.filter_by_id.vpc_channels[0].member_group[0].id
}

data "huaweicloud_apig_channels" "filter_by_member_group_id" {
  depends_on = [
    huaweicloud_apig_channel.test
  ]

  instance_id     = local.instance_id
  member_group_id = local.member_group_id
}

locals {
  member_group_id_filter_result = [
    for v in flatten(data.huaweicloud_apig_channels.filter_by_member_group_id.vpc_channels[*].member_group[*].id) : v == local.member_group_id
  ]
}

output "member_group_id_filter_is_useful" {
  value = length(local.member_group_id_filter_result) > 0 && alltrue(local.member_group_id_filter_result)
}

# Filter by member_group_name
locals {
  member_group_name = huaweicloud_apig_channel.test.member_group[0].name
}

data "huaweicloud_apig_channels" "filter_by_member_group_name" {
  depends_on = [
    huaweicloud_apig_channel.test
  ]

  instance_id       = local.instance_id
  member_group_name = local.member_group_name
}

locals {
  member_group_name_filter_result = [
    for v in data.huaweicloud_apig_channels.filter_by_member_group_name.vpc_channels[0].member_group[*].name : v == local.member_group_name
  ]
}

output "member_group_name_filter_is_useful" {
  value = length(local.member_group_name_filter_result) > 0 && alltrue(local.member_group_name_filter_result)
}
`, testAccChannel_basic(rName))
}

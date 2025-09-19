package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceChannelMemberGroups_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_apig_channel_member_groups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_apig_channel_member_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byExactName   = "data.huaweicloud_apig_channel_member_groups.filter_by_exact_name"
		dcByExactName = acceptance.InitDataSourceCheck(byExactName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceChannelMemberGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "member_groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "member_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "member_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "member_groups.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "member_groups.0.weight"),
					resource.TestMatchResourceAttr(dataSource, "member_groups.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "member_groups.0.update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByExactName.CheckResourceExists(),
					resource.TestCheckOutput("exact_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceChannelMemberGroups_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[2]s"
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  name             = "%[1]s"
  port             = 80
  balance_strategy = 1
}

resource "huaweicloud_apig_channel_member_group" "test" {
  count = 2

  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id
  name           = format("%[1]s_%%d", count.index)
  description    = "terraform script test."
  weight         = 50
}`, name, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccDataSourceChannelMemberGroups_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

# Query all
data "huaweicloud_apig_channel_member_groups" "test" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id

  depends_on = [
    huaweicloud_apig_channel_member_group.test,
  ]
}

# Filter by name (fuzzy search)
locals {
  member_group_name_prefix = "%[2]s"
}

data "huaweicloud_apig_channel_member_groups" "filter_by_name" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id
  name           = local.member_group_name_prefix

  depends_on = [
    huaweicloud_apig_channel_member_group.test,
  ]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_channel_member_groups.filter_by_name.member_groups[*].name : strcontains(v, "%[2]s")
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) >= 2 && alltrue(local.name_filter_result)
}

# Filter by name (exact search)
locals {
  member_group_name = try(data.huaweicloud_apig_channel_member_groups.test.member_groups[0].name, null)
}

data "huaweicloud_apig_channel_member_groups" "filter_by_exact_name" {
  instance_id    = try(data.huaweicloud_apig_instances.test.instances[0].id, "NOT_FOUND")
  vpc_channel_id = huaweicloud_apig_channel.test.id
  name           = local.member_group_name
  precise_search = "member_group_name"

  depends_on = [
    huaweicloud_apig_channel_member_group.test,
  ]
}

output "exact_name_filter_is_useful" {
  value = length(data.huaweicloud_apig_channel_member_groups.filter_by_exact_name.member_groups) == 1
}
`, testAccDataSourceChannelMemberGroups_base(rName), rName)
}

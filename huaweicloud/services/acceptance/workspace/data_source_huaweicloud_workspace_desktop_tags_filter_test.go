package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopTagsFilterDataSource_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_desktop_tags_filter.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterByWithoutTags   = "data.huaweicloud_workspace_desktop_tags_filter.filter_by_without_tags"
		dcFilterByWithoutTags = acceptance.InitDataSourceCheck(filterByWithoutTags)

		filterByAllTags   = "data.huaweicloud_workspace_desktop_tags_filter.filter_by_all_tags"
		dcFilterByAllTags = acceptance.InitDataSourceCheck(filterByAllTags)

		filterByAnyTags   = "data.huaweicloud_workspace_desktop_tags_filter.filter_by_any_tags"
		dcFilterByAnyTags = acceptance.InitDataSourceCheck(filterByAnyTags)

		filterByWithoutAllTags   = "data.huaweicloud_workspace_desktop_tags_filter.filter_by_without_all_tags"
		dcFilterByWithoutAllTags = acceptance.InitDataSourceCheck(filterByWithoutAllTags)

		filterWithoutAnyTags   = "data.huaweicloud_workspace_desktop_tags_filter.filter_by_without_any_tags"
		dcFilterWithoutAnyTags = acceptance.InitDataSourceCheck(filterWithoutAnyTags)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopTagsFilterDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "desktops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "desktops.0.resource_id"),
					resource.TestCheckResourceAttrSet(dcName, "desktops.0.resource_name"),
					resource.TestMatchResourceAttr(dcName, "desktops.0.tags.#", regexp.MustCompile(`(0|^[1-9]([0-9]*)?$)`)),
					// without tags
					dcFilterByWithoutTags.CheckResourceExists(),
					resource.TestCheckOutput("is_without_tags_result_useful", "true"),
					// with all tags
					dcFilterByAllTags.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByAllTags, "desktops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_all_tags_result_useful", "true"),
					// with any tags
					dcFilterByAnyTags.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByAnyTags, "desktops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_any_tags_result_useful", "true"),
					// without all tags
					dcFilterByWithoutAllTags.CheckResourceExists(),
					resource.TestCheckOutput("is_without_all_tags_result_useful", "true"),
					// without any tags
					dcFilterWithoutAnyTags.CheckResourceExists(),
					resource.TestCheckOutput("is_without_any_tags_result_useful", "true"),
				),
			},
		},
	})
}

func testAccDesktopTagsFilterDataSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  tag_key   = "owner"
  tag_value = "terraform"
}

data "huaweicloud_workspace_desktop_tags_filter" "test" {
  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

# filter by without tags
data "huaweicloud_workspace_desktop_tags_filter" "filter_by_without_tags" {
  without_any_tag = true

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  without_tags_result = [
    for v in data.huaweicloud_workspace_desktop_tags_filter.filter_by_without_tags.desktops : length(v.tags) == 0
  ]
}

output "is_without_tags_result_useful" {
  value = alltrue(local.without_tags_result)
}

# filter by all tags 
data "huaweicloud_workspace_desktop_tags_filter" "filter_by_all_tags" {
  tags {
    key    = local.tag_key
    values = [local.tag_value]
  }

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  all_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.filter_by_all_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]
  ])
}

output "is_all_tags_result_useful" {
  value = contains(local.all_tag_keys, local.tag_key)
}

# filter by any tags
data "huaweicloud_workspace_desktop_tags_filter" "filter_by_any_tags" {
  tags_any {
    key    = local.tag_key
    values = [local.tag_value]
  }

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  any_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.filter_by_any_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]
  ])
}

output "is_any_tags_result_useful" {
  value = contains(local.any_tag_keys, local.tag_key)
}

# filter by without all tags
data "huaweicloud_workspace_desktop_tags_filter" "filter_by_without_all_tags" {
  not_tags {
    key    = local.tag_key
    values = [local.tag_value]
  }

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  without_all_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.filter_by_without_all_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]
  ])
}

output "is_without_all_tags_result_useful" {
  value = !contains(local.without_all_tag_keys, local.tag_key)
}

# filter by without any tags
data "huaweicloud_workspace_desktop_tags_filter" "filter_by_without_any_tags" {
  not_tags_any {
    key    = local.tag_key
    values = [local.tag_value]
  }

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  without_any_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.filter_by_without_any_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]
  ])
}

output "is_without_any_tags_result_useful" {
  value = !contains(local.without_any_tag_keys, local.tag_key)
}
`, testAccDesktopTagsFilterDataSource_base(name))
}

func testAccDesktopTagsFilterDataSource_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_vpc_subnets" "test" {
  vpc_id = data.huaweicloud_workspace_service.test.vpc_id
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  name_regex = "WORKSPACE"
  visibility = "market"
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = data.huaweicloud_workspace_flavors.test.flavors[0].id
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "")
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    data.huaweicloud_workspace_service.test.desktop_security_group[0].id,
    data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id,
  ]
  nic {
    network_id = data.huaweicloud_workspace_service.test.network_ids[0]
  }

  name       = "%[1]s"
  user_name  = "%[1]s-user"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }

  email_notification = true
  delete_user        = true
}
`, name)
}

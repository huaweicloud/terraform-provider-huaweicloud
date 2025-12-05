package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktopTagsFilter_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_desktop_tags_filter.all"
		dc  = acceptance.InitDataSourceCheck(all)

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
				Config: testAccDataDesktopTagsFilter_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "desktops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "desktops.0.resource_id"),
					resource.TestCheckResourceAttrSet(all, "desktops.0.resource_name"),
					resource.TestMatchResourceAttr(all, "desktops.0.tags.#", regexp.MustCompile(`(0|^[1-9]([0-9]*)?$)`)),
					// Filter by 'without_any_tag' parameter.
					dcFilterByWithoutTags.CheckResourceExists(),
					resource.TestCheckOutput("is_without_tags_result_useful", "true"),
					// Filter by 'tags' parameter.
					dcFilterByAllTags.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByAllTags, "desktops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_all_tags_result_useful", "true"),
					// Filter by 'tags_any' parameter.
					dcFilterByAnyTags.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByAnyTags, "desktops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_any_tags_result_useful", "true"),
					// Filter by 'not_tags' parameter.
					dcFilterByWithoutAllTags.CheckResourceExists(),
					resource.TestCheckOutput("is_without_all_tags_result_useful", "true"),
					// Filter by 'not_tags_any' parameter.
					dcFilterWithoutAnyTags.CheckResourceExists(),
					resource.TestCheckOutput("is_without_any_tags_result_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDesktopTagsFilter_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  tag_key   = "owner"
  tag_value = "terraform"
}

# Without any filter parameter.
data "huaweicloud_workspace_desktop_tags_filter" "all" {
  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

# Filter by 'without_any_tag' parameter.
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

# Filter by 'tags' parameter.
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

# Filter by 'tags_any' parameter.
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

# Filter by 'not_tags' parameter.
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

# Filter by 'not_tags_any' parameter.
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
`, testAccDataDesktopTags_base(name))
}

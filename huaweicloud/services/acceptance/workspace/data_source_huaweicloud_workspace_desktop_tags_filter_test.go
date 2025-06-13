package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkspaceDesktopTagsFilterDataSource_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_desktop_tags_filter.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		withoutTagsName = "data.huaweicloud_workspace_desktop_tags_filter.test_without_tags"
		withoutTagsDc   = acceptance.InitDataSourceCheck(withoutTagsName)

		withAllTagsName = "data.huaweicloud_workspace_desktop_tags_filter.test_with_all_tags"
		withAllTagsDc   = acceptance.InitDataSourceCheck(withAllTagsName)

		withAnyTagsName = "data.huaweicloud_workspace_desktop_tags_filter.test_with_any_tags"
		withAnyTagsDc   = acceptance.InitDataSourceCheck(withAnyTagsName)

		withoutAllTagsName = "data.huaweicloud_workspace_desktop_tags_filter.test_without_all_tags"
		withoutAllTagsDc   = acceptance.InitDataSourceCheck(withoutAllTagsName)

		withoutAnyTagsName = "data.huaweicloud_workspace_desktop_tags_filter.test_without_any_tags"
		withoutAnyTagsDc   = acceptance.InitDataSourceCheck(withoutAnyTagsName)

		tagKey   = "terraform_test"
		tagValue = "for_terraform"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceDesktopTagsFilterDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "desktops.#"),
					resource.TestCheckResourceAttrSet(dcName, "desktops.0.resource_id"),
					resource.TestCheckResourceAttrSet(dcName, "desktops.0.resource_name"),
					resource.TestCheckResourceAttrSet(dcName, "desktops.0.tags.#"),
				),
			},
			{
				Config: testAccWorkspaceDesktopTagsFilterDataSource_withoutTags(),
				Check: resource.ComposeTestCheckFunc(
					withoutTagsDc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(withoutTagsName, "desktops.#"),
					resource.TestCheckOutput("is_without_tags_result_useful", "true"),
				),
			},
			{
				Config: testAccWorkspaceDesktopTagsFilterDataSource_withAllTags(tagKey, tagValue),
				Check: resource.ComposeTestCheckFunc(
					withAllTagsDc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(withAllTagsName, "desktops.#"),
					resource.TestCheckOutput("is_all_tags_result_useful", "true"),
				),
			},
			{
				Config: testAccWorkspaceDesktopTagsFilterDataSource_withAnyTags(tagKey, tagValue),
				Check: resource.ComposeTestCheckFunc(
					withAnyTagsDc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(withAnyTagsName, "desktops.#"),
					resource.TestCheckOutput("is_any_tags_result_useful", "true"),
				),
			},
			{
				Config: testAccWorkspaceDesktopTagsFilterDataSource_withoutAllTags(tagKey, tagValue),
				Check: resource.ComposeTestCheckFunc(
					withoutAllTagsDc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(withoutAllTagsName, "desktops.#"),
					resource.TestCheckOutput("is_without_all_tags_result_useful", "true"),
				),
			},
			{
				Config: testAccWorkspaceDesktopTagsFilterDataSource_withoutAnyTags(tagKey, tagValue),
				Check: resource.ComposeTestCheckFunc(
					withoutAnyTagsDc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(withoutAnyTagsName, "desktops.#"),
					resource.TestCheckOutput("is_without_any_tags_result_useful", "true"),
				),
			},
		},
	})
}

func testAccWorkspaceDesktopTagsFilterDataSource_basic() string {
	return `data "huaweicloud_workspace_desktop_tags_filter" "test" {}`
}

func testAccWorkspaceDesktopTagsFilterDataSource_withoutTags() string {
	return `
data "huaweicloud_workspace_desktop_tags_filter" "test_without_tags" {
  without_any_tag = true
}

locals {
  without_tags_result = [
    for v in data.huaweicloud_workspace_desktop_tags_filter.test_without_tags.desktops : length(v.tags) == 0
  ]
}

output "is_without_tags_result_useful" {
  value = alltrue(local.without_tags_result)
}
`
}

func testAccWorkspaceDesktopTagsFilterDataSource_withAllTags(tagKey string, tagValue string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_desktop_tags_filter" "test_with_all_tags" {
  tags {
    key = "%[1]s"
    values = ["%[2]s"]
  }
}

locals {
  # 获取所有桌面的标签key
  all_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.test_with_all_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]
  ])
  # 检查是否包含我们指定的key
  has_env_key = contains(local.all_tag_keys, "%[1]s")
}

output "is_all_tags_result_useful" {
  value = local.has_env_key
}
`, tagKey, tagValue)
}

func testAccWorkspaceDesktopTagsFilterDataSource_withAnyTags(tagKey string, tagValue string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_desktop_tags_filter" "test_with_any_tags" {
  tags_any {
    key = "%[1]s"
    values = ["%[2]s"]
  }
}	

locals {
  # 获取所有桌面的标签key
  all_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.test_with_any_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]	
  ])
  # 检查是否包含我们指定的key
  has_env_key = contains(local.all_tag_keys, "%[1]s")
}

output "is_any_tags_result_useful" {
  value = local.has_env_key
}
`, tagKey, tagValue)
}

func testAccWorkspaceDesktopTagsFilterDataSource_withoutAllTags(tagKey string, tagValue string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_desktop_tags_filter" "test_without_all_tags" {
  not_tags {
    key = "%[1]s"
    values = ["%[2]s"]
  }
}

locals {
  # 获取所有桌面的标签key
  all_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.test_without_all_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]
  ])
  # 检查是否包含我们指定的key
  has_env_key = contains(local.all_tag_keys, "%[1]s")
}

output "is_without_all_tags_result_useful" {
  value = !local.has_env_key
}
`, tagKey, tagValue)
}

func testAccWorkspaceDesktopTagsFilterDataSource_withoutAnyTags(tagKey string, tagValue string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_desktop_tags_filter" "test_without_any_tags" {
  not_tags_any {
    key = "%[1]s"
    values = ["%[2]s"]
  }
}

locals {
  # 获取所有桌面的标签key
  all_tag_keys = flatten([
    for desktop in data.huaweicloud_workspace_desktop_tags_filter.test_without_any_tags.desktops : [
      for tag in desktop.tags : tag.key
    ]
  ])
  # 检查是否包含我们指定的key
  has_env_key = contains(local.all_tag_keys, "%[1]s")
}

output "is_without_any_tags_result_useful" {
  value = !local.has_env_key
}
`, tagKey, tagValue)
}

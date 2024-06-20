package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHostGroups_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_lts_host_groups.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byId   = "data.huaweicloud_lts_host_groups.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_lts_host_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_lts_host_groups.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byType   = "data.huaweicloud_lts_host_groups.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byTags   = "data.huaweicloud_lts_host_groups.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHostGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSource, "groups.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "groups.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttr(byId, "groups.0.host_ids.#", "2"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_not_found_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceHostGroups_base(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lts_host_group" "test" {
  depends_on = [null_resource.test]

  name = "%s"
  type = "linux"

  host_ids = huaweicloud_compute_instance.test[*].id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testHostGroup_base(name), name)
}

func testDataSourceHostGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_host_groups" "test" {
  depends_on = [
    huaweicloud_lts_host_group.test
  ]
}

# Filter by ID
locals {
  group_id = huaweicloud_lts_host_group.test.id
}

data "huaweicloud_lts_host_groups" "filter_by_id" {
  host_group_id = local.group_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_lts_host_groups.filter_by_id.groups[*].id : v == local.group_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  group_name = data.huaweicloud_lts_host_groups.test.groups[0].name
}

data "huaweicloud_lts_host_groups" "filter_by_name" {
  name = local.group_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_lts_host_groups.filter_by_name.groups[*].name : v == local.group_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by not exist name
locals {
  not_found_name = "not_found_name"
}

data "huaweicloud_lts_host_groups" "filter_by_not_found_name" {
  name = local.not_found_name
}

locals {
  not_found_name_filter_result = [
    for v in data.huaweicloud_lts_host_groups.filter_by_not_found_name.groups[*].name : strcontains(v, local.not_found_name)
  ]
}

output "is_name_not_found_filter_useful" {
  value = length(local.not_found_name_filter_result) == 0
}

# Filter by type
locals {
  group_type = data.huaweicloud_lts_host_groups.test.groups[0].type
}

data "huaweicloud_lts_host_groups" "filter_by_type" {
  type = local.group_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_lts_host_groups.filter_by_type.groups[*].type : v == local.group_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by tags
locals {
  tags = huaweicloud_lts_host_group.test.tags
}

data "huaweicloud_lts_host_groups" "filter_by_tags" {
  tags = local.tags
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_lts_host_groups.filter_by_tags.groups[*].tags : v == local.tags
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}
`, testDataSourceHostGroups_base(name))
}

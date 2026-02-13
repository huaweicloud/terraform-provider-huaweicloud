package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataResourceInstances_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_organizations_resource_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byWithoutAnyTag   = "data.huaweicloud_organizations_resource_instances.filter_by_without_any_tag"
		dcByWithoutAnyTag = acceptance.InitDataSourceCheck(byWithoutAnyTag)

		byTags   = "data.huaweicloud_organizations_resource_instances.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)

		byMatches   = "data.huaweicloud_organizations_resource_instances.filter_by_matches"
		dcByMatches = acceptance.InitDataSourceCheck(byMatches)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataResourceInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(all, "resources.0.resource_name"),
					// Filter by 'without_any_tag' parameter.
					dcByWithoutAnyTag.CheckResourceExists(),
					resource.TestCheckOutput("is_without_any_tag_filter_useful", "true"),
					resource.TestMatchResourceAttr(byTags, "resources.0.tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byTags, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(byTags, "resources.0.tags.0.value"),
					// Filter by 'tags' parameter.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					// Filter by 'matches' parameter.
					dcByMatches.CheckResourceExists(),
					resource.TestCheckOutput("is_matches_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataResourceInstances_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

locals {
  tag_key   = "foo"
  tag_value = "bar"
}

resource "huaweicloud_organizations_organizational_unit" "test" {
  count = 2

  name      = "%[1]s${count.index}"
  parent_id = data.huaweicloud_organizations_organization.test.root_id

  tags = count.index == 0 ? null : {
    "${local.tag_key}" = local.tag_value
  }
}
`, name)
}

func testDataResourceInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_organizations_resource_instances" "test" {
  resource_type = "organizations:ous"

  depends_on = [huaweicloud_organizations_organizational_unit.test]
}

# Filter by 'without_any_tag' parameter.
data "huaweicloud_organizations_resource_instances" "filter_by_without_any_tag" {
  resource_type   = "organizations:ous"
  without_any_tag = true

  depends_on = [huaweicloud_organizations_organizational_unit.test]
}

locals {
  without_any_tag_filter_result = [for v in data.huaweicloud_organizations_resource_instances.filter_by_without_any_tag.resources[*].tags :
  length(v) == 0]
}

output "is_without_any_tag_filter_useful" {
  value = length(local.without_any_tag_filter_result) > 0 && alltrue(local.without_any_tag_filter_result)
}

# Filter by 'tags' parameter.
data "huaweicloud_organizations_resource_instances" "filter_by_tags" {
  resource_type = "organizations:ous"

  tags {
    key    = local.tag_key
    values = [local.tag_value]
  }

  depends_on = [huaweicloud_organizations_organizational_unit.test]
}

locals {
  tags_filter_result = [for v in flatten(data.huaweicloud_organizations_resource_instances.filter_by_tags.resources[*].tags) :
  v.key == local.tag_key && v.value == local.tag_value]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}

# Filter by 'matches' parameter.
locals {
  resource_name = huaweicloud_organizations_organizational_unit.test[0].name
}

data "huaweicloud_organizations_resource_instances" "filter_by_matches" {
  resource_type = "organizations:ous"

  matches {
    key   = "resource_name"
    value = local.resource_name
  }

  depends_on = [huaweicloud_organizations_organizational_unit.test]
}

locals {
  matches_filter_result = [for v in data.huaweicloud_organizations_resource_instances.filter_by_matches.resources[*].resource_name :
  v == local.resource_name]
}

output "is_matches_filter_useful" {
  value = length(local.matches_filter_result) > 0 && alltrue(local.matches_filter_result)
}
`, testAccDataResourceInstances_base(name))
}

package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceInstancesFilter_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource   = "data.huaweicloud_apig_instances_filter.all"
		dcDataSource = acceptance.InitDataSourceCheck(dataSource)

		byWithoutAnyTag   = "data.huaweicloud_apig_instances_filter.filter_by_without_any_tag"
		dcByWithoutAnyTag = acceptance.InitDataSourceCheck(byWithoutAnyTag)

		byTag   = "data.huaweicloud_apig_instances_filter.filter_by_tags"
		dcByTag = acceptance.InitDataSourceCheck(byTag)

		byMatchesFuzzy   = "data.huaweicloud_apig_instances_filter.filter_by_matches_fuzzy"
		dcByMatchesFuzzy = acceptance.InitDataSourceCheck(byMatchesFuzzy)

		byMatches   = "data.huaweicloud_apig_instances_filter.filter_by_matches"
		dcByMatches = acceptance.InitDataSourceCheck(byMatches)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstancesFilter_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcDataSource.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.resource_name"),
					dcByWithoutAnyTag.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByTag.CheckResourceExists(),
					resource.TestCheckOutput("is_without_any_tag_filter_useful", "true"),
					dcByMatchesFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_matches_fuzzy_filter_useful", "true"),
					dcByMatches.CheckResourceExists(),
					resource.TestCheckOutput("is_matches_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInstancesFilter_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_apig_instance" "test" {
  count = 2

  availability_zones    = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  edition               = "BASIC"
  name                  = "%[2]s${count.index}"
  enterprise_project_id = "%[3]s"
  tags = count.index == 0 ? {
    owner = "terraform"
  } : null
}
`, common.TestBaseNetwork(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataSourceInstancesFilter_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_instances_filter" "all" {
  depends_on = [huaweicloud_apig_instance.test]
}

# Filter by without_any_tag.
data "huaweicloud_apig_instances_filter" "filter_by_without_any_tag" {
  without_any_tag = true

  depends_on = [huaweicloud_apig_instance.test]
}

locals {
  by_without_any_tag_filter = [
    for tag in data.huaweicloud_apig_instances_filter.filter_by_without_any_tag.instances[*].tags : length(tag) == 0
  ]
}

output "is_without_any_tag_filter_useful" {
  value = length(local.by_without_any_tag_filter) > 0 && alltrue(local.by_without_any_tag_filter)
}

# Filter by tags.
locals {
  tag_key   = "owner"
  tag_value = "terraform"
}

data "huaweicloud_apig_instances_filter" "filter_by_tags" {
  tags {
    key    = local.tag_key
    values = [local.tag_value]
  }

  depends_on = [huaweicloud_apig_instance.test]
}

locals {
  by_tags_result = [
    for tag in flatten(data.huaweicloud_apig_instances_filter.filter_by_tags.instances[*].tags) :
    true if tag.key == local.tag_key && tag.value == local.tag_value
  ]
}

output "is_tags_filter_useful" {
  value = length(local.by_tags_result) > 0 && alltrue(local.by_tags_result)
}

# Filter by matches.
locals {
  instance_name        = huaweicloud_apig_instance.test[0].name
  instance_name_prefix = "tf_test"
}

data "huaweicloud_apig_instances_filter" "filter_by_matches_fuzzy" {
  matches {
    key   = "resource_name"
    value = local.instance_name_prefix
  }

  depends_on = [huaweicloud_apig_instance.test]
}

locals {
  by_matches_fuzzy_result = [for v in data.huaweicloud_apig_instances_filter.filter_by_matches_fuzzy.instances[*].resource_name :
    strcontains(v, local.instance_name_prefix)
  ]
}

output "is_matches_fuzzy_filter_useful" {
  value = length(local.by_matches_fuzzy_result) > 0 && alltrue(local.by_matches_fuzzy_result)
}

# Filter by matches.
data "huaweicloud_apig_instances_filter" "filter_by_matches" {
  matches {
    key   = "resource_name"
    value = local.instance_name
  }

  depends_on = [huaweicloud_apig_instance.test]
}

locals {
  by_matches_result = [for v in data.huaweicloud_apig_instances_filter.filter_by_matches.instances[*].resource_name :
    v == local.instance_name
  ]
}

output "is_matches_filter_useful" {
  value = length(local.by_matches_result) > 0 && alltrue(local.by_matches_result)
}
`, testAccDataSourceInstancesFilter_base(name))
}

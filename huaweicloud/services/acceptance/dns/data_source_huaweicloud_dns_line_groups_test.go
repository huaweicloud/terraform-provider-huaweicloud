package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataLineGroups_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_dns_line_group.test"

		dataSource = "data.huaweicloud_dns_line_groups.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLineId           = "data.huaweicloud_dns_line_groups.filter_by_line_id"
		dcByLineId         = acceptance.InitDataSourceCheck(byLineId)
		byLineIdFuzzy      = "data.huaweicloud_dns_line_groups.filter_by_line_id_fuzzy"
		dcByLineIdFuzzy    = acceptance.InitDataSourceCheck(byLineIdFuzzy)
		byNotFoundLineId   = "data.huaweicloud_dns_line_groups.filter_by_not_found_line_id"
		dcByNotFoundLineId = acceptance.InitDataSourceCheck(byNotFoundLineId)

		byName           = "data.huaweicloud_dns_line_groups.filter_by_name"
		dcByName         = acceptance.InitDataSourceCheck(byName)
		byNameFuzzy      = "data.huaweicloud_dns_line_groups.filter_by_name_fuzzy"
		dcByNameFuzzy    = acceptance.InitDataSourceCheck(byNameFuzzy)
		byNotFoundName   = "data.huaweicloud_dns_line_groups.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLineGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "groups.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by line group ID.
					dcByLineId.CheckResourceExists(),
					resource.TestCheckOutput("is_line_id_filter_useful", "true"),
					dcByLineIdFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_line_id_fuzzy_filter_useful", "true"),
					dcByNotFoundLineId.CheckResourceExists(),
					resource.TestCheckOutput("line_id_not_found_validation_pass", "true"),
					// Exact match by line group name.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Fuzzy match by line group name.
					dcByNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_name_fuzzy_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrPair(byLineId, "groups.0.lines", rName, "lines"),
					resource.TestCheckResourceAttrPair(byLineId, "groups.0.description", rName, "description"),
					resource.TestCheckResourceAttrPair(byLineId, "groups.0.status", rName, "status"),
					// Time is not the time corresponding to the local computer.
					resource.TestCheckResourceAttrPair(byLineId, "groups.0.created_at", rName, "created_at"),
					resource.TestCheckResourceAttrPair(byLineId, "groups.0.updated_at", rName, "updated_at"),
				),
			},
		},
	})
}

func testAccDataLineGroups_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_line_group" "test" {
  name        = "%[1]s"
  description = "test description"
  lines       = ["Dianxin_Tianjin", "Dianxin_Jilin"]
}

data "huaweicloud_dns_line_groups" "test" {
  depends_on = [huaweicloud_dns_line_group.test]
}

# Filter by line group ID.
locals {
  line_id = huaweicloud_dns_line_group.test.id
}

data "huaweicloud_dns_line_groups" "filter_by_line_id" {
  line_id = local.line_id
}

locals {
  line_id_filter_result = [for v in data.huaweicloud_dns_line_groups.filter_by_line_id.groups[*].id : v == local.line_id]
}

output "is_line_id_filter_useful" {
  value = length(local.line_id_filter_result) > 0 && alltrue(local.line_id_filter_result)
}

locals {
  line_id_fuzzy = substr(huaweicloud_dns_line_group.test.id, 0, 6)
}

data "huaweicloud_dns_line_groups" "filter_by_line_id_fuzzy" {
  line_id = local.line_id_fuzzy
}

locals {
  line_id_fuzzy_filter_result = [for v in data.huaweicloud_dns_line_groups.filter_by_line_id_fuzzy.groups[*].id :
  strcontains(v, local.line_id_fuzzy)]
}

output "is_line_id_fuzzy_filter_useful" {
  value = length(local.line_id_fuzzy_filter_result) > 0 && alltrue(local.line_id_fuzzy_filter_result)
}

data "huaweicloud_dns_line_groups" "filter_by_not_found_line_id" {
  depends_on = [huaweicloud_dns_line_group.test]
  line_id    = "not_found_line_id"
}

output "line_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_line_groups.filter_by_not_found_line_id.groups) == 0
}

# Exact match by line group name.
locals {
  name = huaweicloud_dns_line_group.test.name
}

data "huaweicloud_dns_line_groups" "filter_by_name" {
  depends_on = [huaweicloud_dns_line_group.test]
  name       = local.name
}

locals {
  name_filter_result = [for v in data.huaweicloud_dns_line_groups.filter_by_name.groups[*].name : v == local.name]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Fuzzy match by line group name.
locals {
  name_prefix = "tf_test"
}

data "huaweicloud_dns_line_groups" "filter_by_name_fuzzy" {
  depends_on = [huaweicloud_dns_line_group.test]
  name       = local.name_prefix
}

locals {
  name_fuzzy_filter_result = [for v in data.huaweicloud_dns_line_groups.filter_by_name.groups[*].name :
  strcontains(v, local.name_prefix)]
}

output "is_name_fuzzy_filter_useful" {
  value = length(local.name_fuzzy_filter_result) > 0 && alltrue(local.name_fuzzy_filter_result)
}

data "huaweicloud_dns_line_groups" "filter_by_not_found_name" {
  depends_on = [huaweicloud_dns_line_group.test]
  name       = "not_found_name"
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_line_groups.filter_by_not_found_name.groups) == 0
}
`, name)
}

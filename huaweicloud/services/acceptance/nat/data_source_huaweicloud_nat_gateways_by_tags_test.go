package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNatGatewaysByTags_basic(t *testing.T) {
	var (
		filterByTagsAll    = "data.huaweicloud_nat_gateways_by_tags.filter_by_tags"
		dcTags             = acceptance.InitDataSourceCheck(filterByTagsAll)
		filterByMatchesAll = "data.huaweicloud_nat_gateways_by_tags.filter_by_matches"
		dcMatches          = acceptance.InitDataSourceCheck(filterByMatchesAll)
		testName           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNatGatewaysByTags_basic(testName),
				Check: resource.ComposeTestCheckFunc(
					dcTags.CheckResourceExists(),
					dcMatches.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(filterByTagsAll, "resources.#"),
					resource.TestCheckResourceAttrSet(filterByTagsAll, "resources.0.resource_id"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.resource_name", testName),
					resource.TestCheckResourceAttrSet(filterByTagsAll, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.0.key", "foo0"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.0.value", "bar0"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.1.key", "foo1"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.1.value", "bar1"),

					resource.TestCheckResourceAttrSet(filterByMatchesAll, "resources.#"),
					resource.TestCheckResourceAttrSet(filterByMatchesAll, "resources.0.resource_id"),
					resource.TestCheckResourceAttr(filterByMatchesAll, "resources.0.resource_name", testName),
					resource.TestCheckResourceAttrSet(filterByMatchesAll, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(filterByMatchesAll, "resources.0.tags.0.key", "foo0"),
					resource.TestCheckResourceAttr(filterByMatchesAll, "resources.0.tags.0.value", "bar0"),
					resource.TestCheckResourceAttr(filterByMatchesAll, "resources.0.tags.1.key", "foo1"),
					resource.TestCheckResourceAttr(filterByMatchesAll, "resources.0.tags.1.value", "bar1"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
				),
			},
		},
	},
	)
}

func testDataSourceNatGatewaysByTags_base(name string) string {
	return fmt.Sprintf(`

data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_nat_gateway" "test" {
  name      = "%[1]s"
  spec      = "1"
  vpc_id    = data.huaweicloud_vpc.myvpc.id
  subnet_id = data.huaweicloud_vpc_subnet.test.id

  tags = {
    "foo0" = "bar0"
    "foo1" = "bar1"
  }
}

`, name)
}

func testDataSourceNatGatewaysByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_gateways_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_nat_gateway.test]
  action     = "filter"

  tags {
    key    = "foo0"
    values = ["bar0"]
  }

  tags {
    key    = "foo1"
    values = ["bar1"]
  }
}

locals {
  tag_key   = "foo0"
  tag_value = "bar0"
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_nat_gateways_by_tags.filter_by_tags.resources) > 0 && alltrue(
    [for v in data.huaweicloud_nat_gateways_by_tags.filter_by_tags.resources[*].tags : anytrue(
    [for vv in v[*].key : vv == local.tag_key]) && anytrue([for vv in v[*].value : vv == local.tag_value])]
  )
}

data "huaweicloud_nat_gateways_by_tags" "filter_by_matches" {
  depends_on = [huaweicloud_nat_gateway.test]
  action     = "filter"

  tags {
    key    = "foo0"
    values = ["bar0"]
  }

  tags {
    key    = "foo1"
    values = ["bar1"]
  }

  matches {
    key   = "resource_name"
    value = "%[2]s"
  }
}

locals {
  match_key   = "resource_name"
  match_value = "%[2]s"
}
	
output "matches_filter_is_useful" {
  value = length(data.huaweicloud_nat_gateways_by_tags.filter_by_matches.resources) > 0
}

`, testDataSourceNatGatewaysByTags_base(name), name)
}

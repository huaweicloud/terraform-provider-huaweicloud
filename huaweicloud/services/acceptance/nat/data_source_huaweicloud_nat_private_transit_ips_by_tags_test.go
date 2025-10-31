package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNatPrivateTransitIpsByTags_basic(t *testing.T) {
	var (
		filterByTagsAll    = "data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_tags"
		dcTags             = acceptance.InitDataSourceCheck(filterByTagsAll)
		filterByMatchesAll = "data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_matches"
		dcMatches          = acceptance.InitDataSourceCheck(filterByMatchesAll)
		testName           = acceptance.RandomAccResourceName()
		ipAddress          = "192.168.0.251"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceNatPrivateTransitIpsByTags_basic(testName, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					dcTags.CheckResourceExists(),
					dcMatches.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(filterByTagsAll, "resources.#"),
					resource.TestCheckResourceAttrSet(filterByTagsAll, "resources.0.resource_id"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.resource_name", ipAddress),
					resource.TestCheckResourceAttrSet(filterByTagsAll, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.0.key", "foo0"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.0.value", "bar0"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.1.key", "foo1"),
					resource.TestCheckResourceAttr(filterByTagsAll, "resources.0.tags.1.value", "bar1"),

					resource.TestCheckResourceAttrSet(filterByMatchesAll, "resources.#"),
					resource.TestCheckResourceAttrSet(filterByMatchesAll, "resources.0.resource_id"),
					resource.TestCheckResourceAttr(filterByMatchesAll, "resources.0.resource_name", ipAddress),
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

func testDataSourceNatPrivateTransitIpsByTags_base(name string, ipAddress string) string {
	return fmt.Sprintf(`

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  cidr       = "192.168.0.0/20"
  vpc_id     = resource.huaweicloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_nat_private_transit_ip" "test" {
  subnet_id             = resource.huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"
  ip_address            = "%[2]s"

  tags = {
    "foo0" = "bar0"
    "foo1" = "bar1"
  }
}

`, name, ipAddress)
}

func testDataSourceNatPrivateTransitIpsByTags_basic(name string, ipAddress string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_nat_private_transit_ip.test]
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
  value = length(data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_tags.resources) > 0 && alltrue(
    [for v in data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_tags.resources[*].tags : anytrue(
    [for vv in v[*].key : vv == local.tag_key]) && anytrue([for vv in v[*].value : vv == local.tag_value])]
  )
}

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_matches" {
  depends_on = [huaweicloud_nat_private_transit_ip.test]
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
    key    = "resource_name"
    value  = "%[2]s"
  }
}

locals {
  match_key   = "resource_name"
  match_value = "%[2]s"
}
	
output "matches_filter_is_useful" {
  value = length(data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_matches.resources) > 0
}

`, testDataSourceNatPrivateTransitIpsByTags_base(name, ipAddress), ipAddress)
}

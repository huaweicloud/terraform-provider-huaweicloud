package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourcePrivateTransitIpsByTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_nat_private_transit_ips_by_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePrivateTransitIpsByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(dataSourceName, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.tags.#"),

					resource.TestCheckOutput("count_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_any_filter_is_useful", "true"),
				),
			},
		},
	},
	)
}

func testDataSourcePrivateTransitIpsByTags_base(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_nat_private_transit_ip" "test" {
  subnet_id  = huaweicloud_vpc_subnet.test.id
  ip_address = "192.168.0.68"

  tags = {
    foo = "bar"
    key = "value"
    tf  = "test"
    acc = "abc"
  }
}
`, common.TestVpc(name))
}

func testDataSourcePrivateTransitIpsByTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_nat_private_transit_ips_by_tags" "test" {
  depends_on = [huaweicloud_nat_private_transit_ip.test]
  action     = "filter"
}

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_count" {
  depends_on = [huaweicloud_nat_private_transit_ip.test]
  action     = "count"
}

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_matches" {
  action = "filter"

  matches {
    key   = "resource_name"
    value = data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.resource_name
  }
}

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_tags" {
  action = "filter"

  tags {
    key    = data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.0.key
    values = [data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.0.value]
  }
}

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_not_tags" {
  action = "filter"

  not_tags {
    key    = data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.1.key
    values = [data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.1.value]
  }
}

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_tags_any" {
  action = "filter"

  tags_any {
    key    = data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.2.key
    values = [data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.2.value]
  }
}

data "huaweicloud_nat_private_transit_ips_by_tags" "filter_by_not_tags_any" {
  action = "filter"

  not_tags_any {
    key    = data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.3.key
    values = [data.huaweicloud_nat_private_transit_ips_by_tags.test.resources.0.tags.3.value]
  }
}

output "count_filter_is_useful" {
  value = data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_count.total_count == 1
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_matches.resources) == 1
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_tags.resources) == 1
}

output "not_tags_filter_is_useful" {
  value = length(data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_not_tags.resources) == 0
}

output "tags_any_filter_is_useful" {
  value = length(data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_tags_any.resources) == 1
}

output "not_tags_any_filter_is_useful" {
  value = length(data.huaweicloud_nat_private_transit_ips_by_tags.filter_by_not_tags_any.resources) == 0
}
`, testDataSourcePrivateTransitIpsByTags_base(name))
}

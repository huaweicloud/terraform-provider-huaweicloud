package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceDcResourcesByTags_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_dc_resources_by_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcResourcesByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.%"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.sys_tags.%"),

					resource.TestCheckResourceAttrSet("data.huaweicloud_dc_resources_by_tags.total_count_test",
						"total_count"),

					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_any_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcResourcesByTags_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id      = huaweicloud_vpc.test.id
  name        = "%[2]s"
  description = "Created by acc test"

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(name), name)
}

func testDataSourceDcResourcesByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dc_resources_by_tags" "test" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
  action        = "filter"
}

data "huaweicloud_dc_resources_by_tags" "matches_filter" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
  action        = "filter"

  matches {
    key   = "resource_name"
    value = "%[2]s"
  }
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_dc_resources_by_tags.matches_filter.resources) > 0
}

data "huaweicloud_dc_resources_by_tags" "tags_filter" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
  action        = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_dc_resources_by_tags.tags_filter.resources) > 0
}

data "huaweicloud_dc_resources_by_tags" "not_tags_filter" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
  action        = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}

output "not_tags_filter_is_useful" {
  value = length(data.huaweicloud_dc_resources_by_tags.not_tags_filter.resources) > 0
}

data "huaweicloud_dc_resources_by_tags" "tags_any_filter" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
  action        = "filter"

  tags_any {
    key    = "foo"
    values = ["bar"]
  }
}

output "tags_any_filter_is_useful" {
  value = length(data.huaweicloud_dc_resources_by_tags.tags_any_filter.resources) > 0
}

data "huaweicloud_dc_resources_by_tags" "not_tags_any_filter" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
  action        = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}

output "not_tags_any_filter_is_useful" {
  value = length(data.huaweicloud_dc_resources_by_tags.not_tags_any_filter.resources) > 0
}

data "huaweicloud_dc_resources_by_tags" "total_count_test" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
  action        = "count"
}
`, testDataSourceDcResourcesByTags_base(name), name)
}

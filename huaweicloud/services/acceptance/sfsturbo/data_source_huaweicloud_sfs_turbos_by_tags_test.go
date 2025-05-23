package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceSfsTurbosByTags_basic(t *testing.T) {
	var (
		dataSource  = "data.huaweicloud_sfs_turbos_by_tags.test"
		dataSource1 = "data.huaweicloud_sfs_turbos_by_tags.filter_by_count"
		dataSource2 = "data.huaweicloud_sfs_turbos_by_tags.filter_by_tags"
		dataSource3 = "data.huaweicloud_sfs_turbos_by_tags.filter_by_matches"

		name = acceptance.RandomAccResourceName()
		dc   = acceptance.InitDataSourceCheck(dataSource)
		dc1  = acceptance.InitDataSourceCheck(dataSource1)
		dc2  = acceptance.InitDataSourceCheck(dataSource2)
		dc3  = acceptance.InitDataSourceCheck(dataSource3)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSfsTurbosByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),

					dc1.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource1, "total_count"),
					resource.TestCheckOutput("results_is_not_empty", "true"),

					dc2.CheckResourceExists(),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),

					dc3.CheckResourceExists(),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccSfsTurbosByTags_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "%[2]s"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceSfsTurbosByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_sfs_turbos_by_tags" "test" {
  action = "filter"

  depends_on = [huaweicloud_sfs_turbo.test]
}

data "huaweicloud_sfs_turbos_by_tags" "filter_by_count" {
  action = "count"

  depends_on = [huaweicloud_sfs_turbo.test]
}

data "huaweicloud_sfs_turbos_by_tags" "filter_by_tags" {
  action = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }

  depends_on = [huaweicloud_sfs_turbo.test]
}

data "huaweicloud_sfs_turbos_by_tags" "filter_by_matches" {
  action = "filter"

  matches {
    key   = "resource_name"
    value = huaweicloud_sfs_turbo.test.name
  }

  depends_on = [huaweicloud_sfs_turbo.test]
}

output "results_is_not_empty" {
  value = data.huaweicloud_sfs_turbos_by_tags.filter_by_count.total_count > 0
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_sfs_turbos_by_tags.filter_by_tags.resources) == 1
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_sfs_turbos_by_tags.filter_by_matches.resources) == 1
}

`, testAccSfsTurbosByTags_base(name))
}

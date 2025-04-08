package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNetworkingSecGroupsByTags_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_networking_secgroups_by_tags.basic"
	dataSource2 := "data.huaweicloud_networking_secgroups_by_tags.filter_by_tags"
	dataSource3 := "data.huaweicloud_networking_secgroups_by_tags.filter_by_matches"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceNetworkingSecGroupsByTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					resource.TestCheckOutput("is_matches_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceNetworkingSecGroupsByTags_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test1" {
  name        = "%[1]s-1"
  description = "created by acc test"

  tags = {
    foo = "%[1]s"
    key = "%[1]s_1"
  }
}

resource "huaweicloud_networking_secgroup" "test2" {
  name        = "%[1]s-2"
  description = "created by acc test"

  tags = {
    foo = "%[1]s"
    key = "%[1]s_2"
  }
}

data "huaweicloud_networking_secgroups_by_tags" "basic" {
  action = "filter"

  depends_on = [huaweicloud_networking_secgroup.test1, huaweicloud_networking_secgroup.test2]
}

data "huaweicloud_networking_secgroups_by_tags" "filter_by_tags" {
  action = "filter"

  tags {
    key    = "foo"
    values = ["%[1]s"]
  }

  tags {
    key    = "key"
    values = ["%[1]s_1", "%[1]s_2"]
  }

  depends_on = [huaweicloud_networking_secgroup.test1, huaweicloud_networking_secgroup.test2]
}

data "huaweicloud_networking_secgroups_by_tags" "filter_by_matches" {
  action = "filter"

  matches {
    key   = "resource_name"
    value = "%[1]s-1"
  }

  depends_on = [huaweicloud_networking_secgroup.test1, huaweicloud_networking_secgroup.test2]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_networking_secgroups_by_tags.basic.resources) > 0
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_networking_secgroups_by_tags.filter_by_tags.resources) == 2
}

output "is_matches_filter_useful" {
  value = length(data.huaweicloud_networking_secgroups_by_tags.filter_by_matches.resources) == 1
}
`, name)
}

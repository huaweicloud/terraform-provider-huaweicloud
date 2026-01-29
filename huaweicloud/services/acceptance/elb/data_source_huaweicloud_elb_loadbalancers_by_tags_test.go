package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceElbLoadbalancersByTags_basic(t *testing.T) {
	var (
		datasourceName = "data.huaweicloud_elb_loadbalancers_by_tags.test"
		dc             = acceptance.InitDataSourceCheck(datasourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbLoadbalancersByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(datasourceName, "resources.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "resources.0.tags.#"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("count_filter_is_useful", "true"),
				),
			},
		},
	},
	)
}

func testDataSourceElbLoadbalancersByTags_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
	
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }

  lifecycle {
    ignore_changes = [
      l4_flavor_id, l7_flavor_id
    ]
  }
}
`, common.TestVpc(name), name)
}

func testDataSourceElbLoadbalancersByTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_elb_loadbalancers_by_tags" "test" {
  action = "filter"

  depends_on = [
    huaweicloud_elb_loadbalancer.test
  ]
}

data "huaweicloud_elb_loadbalancers_by_tags" "filter_by_tags" {
  action = "filter"

  tags {
    key    = data.huaweicloud_elb_loadbalancers_by_tags.test.resources.0.tags.0.key
    values = [data.huaweicloud_elb_loadbalancers_by_tags.test.resources.0.tags.0.value]
  }

  depends_on = [
    huaweicloud_elb_loadbalancer.test
  ]
}

data "huaweicloud_elb_loadbalancers_by_tags" "filter_by_matches" {
  action = "filter"

  matches {
    key   = "resource_name"
    value = data.huaweicloud_elb_loadbalancers_by_tags.test.resources.0.resource_name
  }

  depends_on = [
    huaweicloud_elb_loadbalancer.test
  ]
}

data "huaweicloud_elb_loadbalancers_by_tags" "filter_by_count" {
  action = "count"

  depends_on = [
    huaweicloud_elb_loadbalancer.test
  ]
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers_by_tags.filter_by_tags.resources) > 0
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers_by_tags.filter_by_matches.resources) == 1
}

output "count_filter_is_useful" {
  value = data.huaweicloud_elb_loadbalancers_by_tags.filter_by_count.total_count > 0
}
`, testDataSourceElbLoadbalancersByTags_base(name))
}

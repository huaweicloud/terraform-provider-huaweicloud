package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssElbLoadbalancers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_elb_loadbalancers.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssElbLoadbalancers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.l7_flavor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "loadbalancers.0.ip_target_enable"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_id_filter_is_useful", "true"),
					resource.TestCheckOutput("is_cross_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCssElbLoadbalancers_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_elb_loadbalancers" "test" {
  depends_on = [
    huaweicloud_elb_loadbalancer.test1,
    huaweicloud_elb_loadbalancer.test2,
  ]
  
  cluster_id = huaweicloud_css_cluster.test.id
}
  
locals {
  loadbalancer_id = data.huaweicloud_css_elb_loadbalancers.test.loadbalancers[0].id
  name            = data.huaweicloud_css_elb_loadbalancers.test.loadbalancers[0].name
  protocol_id     = data.huaweicloud_css_elb_loadbalancers.test.loadbalancers[0].l7_flavor_id
  is_cross        = data.huaweicloud_css_elb_loadbalancers.test.loadbalancers[0].ip_target_enable
}
	
data "huaweicloud_css_elb_loadbalancers" "filter_by_id" {
  cluster_id      = huaweicloud_css_cluster.test.id
  loadbalancer_id = local.loadbalancer_id
}

data "huaweicloud_css_elb_loadbalancers" "filter_by_name" {
  cluster_id = huaweicloud_css_cluster.test.id
  name       = local.name
}

data "huaweicloud_css_elb_loadbalancers" "filter_by_protocol_id" {
  cluster_id  = huaweicloud_css_cluster.test.id
  protocol_id = local.protocol_id
}

data "huaweicloud_css_elb_loadbalancers" "filter_by_is_cross" {
  cluster_id = huaweicloud_css_cluster.test.id
  is_cross   = local.is_cross
}
  
locals {
  list_by_id          = data.huaweicloud_css_elb_loadbalancers.filter_by_id.loadbalancers
  list_by_name        = data.huaweicloud_css_elb_loadbalancers.filter_by_name.loadbalancers
  list_by_protocol_id = data.huaweicloud_css_elb_loadbalancers.filter_by_protocol_id.loadbalancers
  list_by_is_cross    = data.huaweicloud_css_elb_loadbalancers.filter_by_is_cross.loadbalancers
}

output "id_filter_is_useful" {
  value = length(local.list_by_id) > 0 && alltrue(
    [for v in local.list_by_id[*].id : v == local.loadbalancer_id]
  )
}
	
output "name_filter_is_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}

output "protocol_id_filter_is_useful" {
  value = length(local.list_by_protocol_id) > 0 && alltrue(
    [for v in local.list_by_protocol_id[*].l7_flavor_id : v == local.protocol_id]
  )
}

output "is_cross_filter_is_useful" {
  value = length(local.list_by_is_cross) > 0 && alltrue(
    [for v in local.list_by_is_cross[*].ip_target_enable : v == local.is_cross]
  )
}
`, testDataSourceCssElbLoadbalancers_data_basic(name))
}

func testDataSourceCssElbLoadbalancers_data_basic(name string) string {
	clusterString := testAccCssCluster_basic(name, "Test@passw0rd", 7, "bar")
	return fmt.Sprintf(`
  %[1]s
  
resource "huaweicloud_elb_loadbalancer" "test1" {
  name              = "%[2]s_1"
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[1]
  ]
}
  
resource "huaweicloud_elb_loadbalancer" "test2" {
  name              = "%[2]s_2"
  cross_vpc_backend = true
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[1]
  ]
}
`, clusterString, name)
}

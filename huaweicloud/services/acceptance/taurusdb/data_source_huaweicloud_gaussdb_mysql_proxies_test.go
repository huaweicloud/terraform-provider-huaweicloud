package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbMysqlProxies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_proxies.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussdbMysqlProxies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.delay_threshold_in_seconds"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.node_num"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.connection_pool_type"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.switch_connection_pool_type_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.elb_vip"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.transaction_split"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.balance_route_mode_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.route_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.consistence_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.ssl_option"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.new_node_auto_add_status"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.new_node_weight"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.master_node_weight.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.master_node_weight.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.master_node_weight.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.master_node_weight.0.weight"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.readonly_nodes_weight.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.readonly_nodes_weight.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.readonly_nodes_weight.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.readonly_nodes_weight.0.weight"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.nodes.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.nodes.0.frozen_flag"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussdbMysqlProxies_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_mysql_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = "multi"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                     = "%s"
  password                 = "Test@12345678"
  flavor                   = data.huaweicloud_gaussdb_mysql_flavors.test.flavors[0].name
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  enterprise_project_id    = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
  read_replicas            = 4
}

data "huaweicloud_gaussdb_mysql_proxy_flavors" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}

locals{
  sort_nodes = tolist(values({ for node in huaweicloud_gaussdb_mysql_instance.test.nodes : node.name => node }))
}

resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  flavor      = data.huaweicloud_gaussdb_mysql_proxy_flavors.test.flavor_groups[0].flavors[0].spec_code
  node_num    = 2
  proxy_name  = "%[2]s"
  proxy_mode  = "readwrite"
  route_mode  = 1
  subnet_id   = huaweicloud_vpc_subnet.test.id

  master_node_weight {
    id     = local.sort_nodes[0].id
    weight = 20
  }

  readonly_nodes_weight {
    id     = local.sort_nodes[1].id
    weight = 30
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceDataSourceGaussdbMysqlProxies_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_mysql_proxies" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_proxy.test]

  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}
`, testDataSourceDataSourceGaussdbMysqlProxies_base(name))
}

package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceRdsMysqlProxies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_mysql_proxies.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsMysqlProxies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.node_num"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.flavor_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.flavor_info.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.flavor_info.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.proxy_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.route_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.delay_threshold_in_seconds"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.memory"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.nodes.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.nodes.0.frozen_flag"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.transaction_split"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.connection_pool_type"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.pay_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.seconds_level_monitor_fun_status"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.alt_flag"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.force_read_only"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.ssl_option"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.support_balance_route_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.support_proxy_ssl"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.support_switch_connection_pool_type"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy.0.support_transaction_split"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.master_instance.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.master_instance.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.master_instance.0.weight"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.readonly_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.readonly_instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.readonly_instances.0.weight"),
					resource.TestCheckResourceAttrSet(dataSource, "proxy_list.0.proxy_security_group_check_result"),
					resource.TestCheckResourceAttrSet(dataSource, "max_proxy_num"),
					resource.TestCheckResourceAttrSet(dataSource, "max_proxy_node_num"),
					resource.TestCheckResourceAttrSet(dataSource, "support_balance_route_mode_for_favored_version"),
				),
			},
		},
	})
}

func testDataSourceRdsMysqlProxies_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "instance" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

data "huaweicloud_rds_flavors" "replica" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "replica"
  group_type    = "dedicated"
  memory        = 4
  vcpus         = 2
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.instance.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  db {
    password = "test_1234"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_read_replica_instance" "test" {
  count = 2

  name                = "%[2]s_${count.index}"
  flavor              = data.huaweicloud_rds_flavors.replica.flavors[0].name
  primary_instance_id = huaweicloud_rds_instance.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  security_group_id   = huaweicloud_networking_secgroup.test.id

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

data "huaweicloud_rds_mysql_proxy_flavors" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}

resource "huaweicloud_rds_mysql_proxy" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  flavor      = data.huaweicloud_rds_mysql_proxy_flavors.test.flavor_groups[0].flavors[0].code
  node_num    = 2
  proxy_name  = "%[2]s"
  proxy_mode  = "readwrite"
  route_mode  = 0

  master_node_weight {
    id     = huaweicloud_rds_instance.test.id
    weight = 10
  }

  readonly_nodes_weight {
    id     = huaweicloud_rds_read_replica_instance.test[0].id
    weight = 20
  }

  readonly_nodes_weight {
    id     = huaweicloud_rds_read_replica_instance.test[1].id
    weight = 30
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testDataSourceRdsMysqlProxies_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_mysql_proxies" "test" {
  depends_on = [huaweicloud_rds_mysql_proxy.test]

  instance_id = huaweicloud_rds_instance.test.id
}
`, testDataSourceRdsMysqlProxies_base(name))
}

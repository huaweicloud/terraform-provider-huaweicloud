package rds

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceProxy(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy-list"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	searchExpression := fmt.Sprintf("proxy_query_info_list[?proxy.pool_id=='%s']|[0]", state.Primary.ID)
	proxy := utils.PathSearch(searchExpression, getRespBody, nil)
	if proxy == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return proxy, nil
}

func TestAccMysqlProxy_basic(t *testing.T) {
	var proxy instances.Proxy
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_rds_mysql_proxy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&proxy,
		getResourceProxy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProxy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_rds_mysql_proxy_flavors.test", "flavor_groups.0.flavors.0.code"),
					resource.TestCheckResourceAttr(resourceName, "node_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "proxy_name", rName),
					resource.TestCheckResourceAttr(resourceName, "proxy_mode", "readwrite"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "route_mode", "0"),
					resource.TestCheckResourceAttr(resourceName, "master_node_weight.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "master_node_weight.0.id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "master_node_weight.0.weight", "10"),
					resource.TestCheckResourceAttr(resourceName, "readonly_nodes_weight.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttrSet(resourceName, "delay_threshold_in_seconds"),
					resource.TestCheckResourceAttrSet(resourceName, "vcpus"),
					resource.TestCheckResourceAttrSet(resourceName, "memory"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.role"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.frozen_flag"),
					resource.TestCheckResourceAttrSet(resourceName, "mode"),
					resource.TestCheckResourceAttrSet(resourceName, "flavor_group_type"),
					resource.TestCheckResourceAttrSet(resourceName, "transaction_split"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_pool_type"),
					resource.TestCheckResourceAttrSet(resourceName, "pay_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "seconds_level_monitor_fun_status"),
					resource.TestCheckResourceAttrSet(resourceName, "alt_flag"),
					resource.TestCheckResourceAttrSet(resourceName, "force_read_only"),
					resource.TestCheckResourceAttrSet(resourceName, "ssl_option"),
					resource.TestCheckResourceAttrSet(resourceName, "support_balance_route_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "support_proxy_ssl"),
					resource.TestCheckResourceAttrSet(resourceName, "support_switch_connection_pool_type"),
					resource.TestCheckResourceAttrSet(resourceName, "support_transaction_split"),
				),
			},
			{
				Config: testAccMysqlProxy_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "route_mode", "2"),
					resource.TestCheckResourceAttr(resourceName, "readonly_nodes_weight.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "readonly_nodes_weight.0.id",
						"huaweicloud_rds_read_replica_instance.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "readonly_nodes_weight.0.weight", "50"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testMysqlProxyResourceImportState(resourceName),
				ImportStateVerifyIgnore: []string{"readonly_nodes_weight"},
			},
		},
	})
}

func testAccMysqlProxy_base(rName string) string {
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
`, common.TestBaseNetwork(rName), rName)
}

func testAccMysqlProxy_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

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
`, testAccMysqlProxy_base(rName), rName)
}

func testAccMysqlProxy_basic_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_mysql_proxy" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  flavor      = data.huaweicloud_rds_mysql_proxy_flavors.test.flavor_groups[0].flavors[0].code
  node_num    = 2
  proxy_name  = "%[2]s"
  proxy_mode  = "readwrite"
  route_mode  = 2

  readonly_nodes_weight {
    id     = huaweicloud_rds_read_replica_instance.test[0].id
    weight = 50
  }
}
`, testAccMysqlProxy_base(rName), rName)
}

func testMysqlProxyResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}

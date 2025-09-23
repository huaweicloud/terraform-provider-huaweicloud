package taurusdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/taurusdb/v3/instances"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getResourceProxy(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxies"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	listMysqlDatabasesResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listMysqlDatabasesResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}

	searchExpression := fmt.Sprintf("proxy_list[?proxy.pool_id=='%s']|[0]", state.Primary.ID)
	proxy := utils.PathSearch(searchExpression, listRespBody, nil)
	if proxy == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return proxy, nil
}

func TestAccGaussDBMySQLProxy_basic(t *testing.T) {
	var proxy instances.Proxy
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_gaussdb_mysql_proxy.test"

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
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_mysql_proxy_flavors.test", "flavor_groups.0.flavors.0.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "node_num", "2"),
					resource.TestCheckResourceAttr(resourceName, "proxy_name", rName),
					resource.TestCheckResourceAttr(resourceName, "route_mode", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "new_node_auto_add_status", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "port", "3339"),
					resource.TestCheckResourceAttr(resourceName, "transaction_split", "ON"),
					resource.TestCheckResourceAttr(resourceName, "connection_pool_type", "SESSION"),
					resource.TestCheckResourceAttr(resourceName, "switch_connection_pool_type_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "master_node_weight.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "readonly_nodes_weight.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "multiStatementType"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "Loose"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.elem_type", "system"),
					resource.TestCheckResourceAttr(resourceName, "consistence_mode", "session"),
					resource.TestCheckResourceAttr(resourceName, "open_access_control", "true"),
					resource.TestCheckResourceAttr(resourceName, "access_control_type", "white"),
					resource.TestCheckResourceAttr(resourceName, "access_control_ip_list.0.ip", "3.3.3.3"),
					resource.TestCheckResourceAttr(resourceName, "access_control_ip_list.0.description",
						"test description"),
					resource.TestCheckResourceAttrSet(resourceName, "address"),
					resource.TestCheckResourceAttrSet(resourceName, "current_version"),
					resource.TestCheckResourceAttrSet(resourceName, "can_upgrade"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.status"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.role"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.az_code"),
					resource.TestCheckResourceAttrSet(resourceName, "nodes.0.frozen_flag"),
				),
			},
			{
				Config: testAccMysqlProxy_basic_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "flavor",
						"data.huaweicloud_gaussdb_mysql_proxy_flavors.test", "flavor_groups.0.flavors.1.spec_code"),
					resource.TestCheckResourceAttr(resourceName, "node_num", "4"),
					resource.TestCheckResourceAttr(resourceName, "proxy_name", updateName),
					resource.TestCheckResourceAttr(resourceName, "route_mode", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "new_node_auto_add_status", "ON"),
					resource.TestCheckResourceAttr(resourceName, "new_node_weight", "20"),
					resource.TestCheckResourceAttr(resourceName, "port", "3338"),
					resource.TestCheckResourceAttr(resourceName, "transaction_split", "OFF"),
					resource.TestCheckResourceAttr(resourceName, "connection_pool_type", "CLOSED"),
					resource.TestCheckResourceAttr(resourceName, "switch_connection_pool_type_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "master_node_weight.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "readonly_nodes_weight.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "looseImciApThreshold"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "6000"),
					resource.TestCheckResourceAttr(resourceName, "parameters.0.elem_type", "system"),
					resource.TestCheckResourceAttr(resourceName, "consistence_mode", "eventual"),
					resource.TestCheckResourceAttr(resourceName, "open_access_control", "false"),
					resource.TestCheckResourceAttr(resourceName, "access_control_type", "black"),
					resource.TestCheckResourceAttr(resourceName, "access_control_ip_list.0.ip", "4.4.4.4"),
					resource.TestCheckResourceAttr(resourceName, "access_control_ip_list.0.description",
						"test description update"),
				),
			},
			{
				Config: testAccMysqlProxy_basic_reduce_node(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "node_num", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testGaussDBMysqlProxyResourceImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"new_node_weight",
					"proxy_mode",
					"readonly_nodes_weight",
					"parameters",
				},
			},
		},
	})
}

func testAccMysqlProxy_base(rName string) string {
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
`, common.TestBaseNetwork(rName), rName)
}

func testAccMysqlProxy_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  instance_id              = huaweicloud_gaussdb_mysql_instance.test.id
  flavor                   = data.huaweicloud_gaussdb_mysql_proxy_flavors.test.flavor_groups[0].flavors[0].spec_code
  node_num                 = 2
  proxy_name               = "%[2]s"
  proxy_mode               = "readwrite"
  route_mode               = 1
  subnet_id                = huaweicloud_vpc_subnet.test.id
  new_node_auto_add_status = "OFF"
  new_node_weight          = 20
  port                     = 3339
  transaction_split        = "ON"
  consistence_mode         = "session"
  connection_pool_type     = "SESSION"
  open_access_control      = true
  access_control_type      = "white"

  access_control_ip_list {
    ip          = "3.3.3.3"
    description = "test description"
  }

  master_node_weight {
    id     = local.sort_nodes[0].id
    weight = 20
  }

  readonly_nodes_weight {
    id     = local.sort_nodes[1].id
    weight = 30
  }

  parameters {
    name      = "multiStatementType"
    value     = "Loose"
    elem_type = "system"
  }
}
`, testAccMysqlProxy_base(rName), rName)
}

func testAccMysqlProxy_basic_update(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  instance_id              = huaweicloud_gaussdb_mysql_instance.test.id
  flavor                   = data.huaweicloud_gaussdb_mysql_proxy_flavors.test.flavor_groups[0].flavors[1].spec_code
  node_num                 = 4
  proxy_name               = "%[2]s"
  proxy_mode               = "readwrite"
  route_mode               = 1
  subnet_id                = huaweicloud_vpc_subnet.test.id
  new_node_auto_add_status = "ON"
  new_node_weight          = 20
  port                     = 3338
  transaction_split        = "OFF"
  consistence_mode         = "eventual"
  connection_pool_type     = "CLOSED"
  open_access_control      = false
  access_control_type      = "black"

  access_control_ip_list {
    ip          = "4.4.4.4"
    description = "test description update"
  }

  master_node_weight {
    id     = local.sort_nodes[0].id
    weight = 10
  }

  readonly_nodes_weight {
    id     = local.sort_nodes[2].id
    weight = 20
  }

  readonly_nodes_weight {
    id     = local.sort_nodes[3].id
    weight = 30
  }

  parameters {
    name      = "looseImciApThreshold"
    value     = "6000"
    elem_type = "system"
  }
}
`, testAccMysqlProxy_base(rName), updateName)
}

func testAccMysqlProxy_basic_reduce_node(rName, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  instance_id              = huaweicloud_gaussdb_mysql_instance.test.id
  flavor                   = data.huaweicloud_gaussdb_mysql_proxy_flavors.test.flavor_groups[0].flavors[1].spec_code
  node_num                 = 3
  proxy_name               = "%[2]s"
  proxy_mode               = "readwrite"
  route_mode               = 1
  subnet_id                = huaweicloud_vpc_subnet.test.id
  new_node_auto_add_status = "ON"
  new_node_weight          = 20
  port                     = 3338
  transaction_split        = "OFF"
  consistence_mode         = "eventual"
  connection_pool_type     = "CLOSED"
  open_access_control      = false
  access_control_type      = "black"

  access_control_ip_list {
    ip          = "4.4.4.4"
    description = "test description update"
  }

  master_node_weight {
    id     = local.sort_nodes[0].id
    weight = 10
  }

  readonly_nodes_weight {
    id     = local.sort_nodes[2].id
    weight = 20
  }

  readonly_nodes_weight {
    id     = local.sort_nodes[3].id
    weight = 30
  }

  parameters {
    name      = "looseImciApThreshold"
    value     = "6000"
    elem_type = "system"
  }
}
`, testAccMysqlProxy_base(rName), updateName)
}

func testGaussDBMysqlProxyResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		instanceID := rs.Primary.Attributes["instance_id"]
		return fmt.Sprintf("%s/%s", instanceID, rs.Primary.ID), nil
	}
}

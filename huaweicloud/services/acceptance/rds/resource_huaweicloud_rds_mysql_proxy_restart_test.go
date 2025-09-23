package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccMysqlProxyRestart_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_rds_mysql_proxy_restart.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceProxy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlProxyRestart_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccMysqlProxyRestart_base(rName string) string {
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

func testAccMysqlProxyRestart_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_mysql_proxy_restart" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  proxy_id    = huaweicloud_rds_mysql_proxy.test.id
}`, testAccMysqlProxyRestart_base(rName))
}

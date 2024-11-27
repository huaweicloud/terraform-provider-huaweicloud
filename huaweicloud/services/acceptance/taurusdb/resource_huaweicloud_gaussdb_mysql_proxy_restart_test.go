package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGaussDBMysqlProxyRestart_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_gaussdb_mysql_proxy_restart.test"

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
				Config: testAccGaussDBMysqlProxyRestart_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccGaussDBMysqlProxyRestart_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_mysql_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = "multi"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                     = "%[2]s"
  password                 = "Test@12345678"
  flavor                   = data.huaweicloud_gaussdb_mysql_flavors.test.flavors[0].name
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  enterprise_project_id    = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
  read_replicas            = 2
}

data "huaweicloud_gaussdb_mysql_proxy_flavors" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}

resource "huaweicloud_gaussdb_mysql_proxy" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  flavor      = data.huaweicloud_gaussdb_mysql_proxy_flavors.test.flavor_groups[0].flavors[0].spec_code
  node_num    = 2
  proxy_name  = "%[2]s"
  proxy_mode  = "readwrite"
  route_mode  = 1
  subnet_id   = huaweicloud_vpc_subnet.test.id
}

resource "huaweicloud_gaussdb_mysql_proxy_restart" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  proxy_id    = huaweicloud_gaussdb_mysql_proxy.test.id
}`, common.TestBaseNetwork(rName), rName)
}

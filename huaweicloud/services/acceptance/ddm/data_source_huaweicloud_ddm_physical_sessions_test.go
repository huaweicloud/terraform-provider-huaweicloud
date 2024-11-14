package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceDdmPhysicalSessions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ddm_physical_sessions.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDdmPhysicalSessions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.0.user"),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.0.host"),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.0.db"),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.0.command"),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "physical_processes.0.trx_executed_time"),
				),
			},
		},
	})
}

func testDataSourceDdmPhysicalSessions_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_ddm_engines" test {
  version = "3.0.9"
}

data "huaweicloud_ddm_flavors" test {
  engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
  cpu_arch  = "X86"
}

resource "huaweicloud_networking_secgroup_rule" "ingress" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 3306
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_networking_secgroup_rule" "egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_ddm_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
  node_num          = 2
  engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  admin_user        = "test_user_1"
  admin_password    = "test_password_123"

  availability_zones = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mysql.n1.large.4"
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

resource "huaweicloud_ddm_schema" "test" {
  instance_id  = huaweicloud_ddm_instance.test.id
  name         = "%[2]s"
  shard_mode   = "single"
  shard_number = "1"

  data_nodes {
    id             = huaweicloud_rds_instance.test.id
    admin_user     = "root"
    admin_password = "test_1234"
  }

  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}

resource "huaweicloud_ddm_account" "test" {
  depends_on = [huaweicloud_ddm_schema.test]

  instance_id = huaweicloud_ddm_instance.test.id
  name        = "%[2]s"
  password    = "test_1234"

  permissions = [
    "SELECT"
  ]

  schemas {
    name = huaweicloud_ddm_schema.test.name
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceDdmPhysicalSessions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ddm_physical_sessions" "test" {
  depends_on = [huaweicloud_ddm_account.test]

  instance_id = huaweicloud_rds_instance.test.id
}
`, testDataSourceDdmPhysicalSessions_base(name))
}

package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccOpenGaussInstanceNodeStartup_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_gaussdb_opengauss_instance.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getOpenGaussInstanceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceNodeStartup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
		},
	})
}

func testAccOpenGaussInstanceNodeStartup_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss_egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = "gaussdb.bs.s6.xlarge.x864.ha"
  name              = "%[2]s"
  password          = "Huangwei!120521"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

data "huaweicloud_gaussdb_opengauss_instance_nodes" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}

resource "huaweicloud_gaussdb_opengauss_instance_node_stop" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  node_id     = data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].id
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccOpenGaussInstanceNodeStartup_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_instance_node_startup" "test" {
  depends_on = [huaweicloud_gaussdb_opengauss_instance_node_stop.test]

  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  node_id     = data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].id
}`, testAccOpenGaussInstanceNodeStartup_base(rName))
}

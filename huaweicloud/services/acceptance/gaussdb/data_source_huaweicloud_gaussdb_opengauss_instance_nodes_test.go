package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussInstanceNodes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_instance_nodes.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussInstanceNodes_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.availability_zone_id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.components.#"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.components.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.components.0.role"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.components.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.components.0.type"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussInstanceNodes_base(name string) string {
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

  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  flavor                = "gaussdb.bs.s6.xlarge.x864.ha"
  name                  = "%[2]s"
  password              = "Huangwei!120521"
  sharding_num          = 1
  coordinator_num       = 2
  replica_num           = 3
  enterprise_project_id = "%[3]s"

  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

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
`, common.TestBaseNetwork(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbOpengaussInstanceNodes_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_instance_nodes" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}

locals {
  component_type = "DN"
}
data "huaweicloud_gaussdb_opengauss_instance_nodes" "component_type_filter" {
  instance_id    = huaweicloud_gaussdb_opengauss_instance.test.id
  component_type = "DN"
}
output "component_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_instance_nodes.component_type_filter.nodes) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_instance_nodes.component_type_filter.nodes[*] : alltrue(
  [for vv in v.components : vv.type == local.component_type]
  )]
  )
}

`, testDataSourceGaussdbOpengaussInstanceNodes_base(name))
}

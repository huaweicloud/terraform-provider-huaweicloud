package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussSolutionTemplateSetting_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_solution_template_setting.test"
	dataSourceWithInstanceId := "data.huaweicloud_gaussdb_opengauss_solution_template_setting.test_instance_id"
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
				Config: testDataSourceGaussdbOpengaussSolutionTemplateSetting_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "shard_num"),
					resource.TestCheckResourceAttrSet(dataSource, "replica_num"),
					resource.TestCheckResourceAttrSet(dataSource, "initial_node_num"),

					resource.TestCheckResourceAttrSet(dataSourceWithInstanceId, "shard_num"),
					resource.TestCheckResourceAttrSet(dataSourceWithInstanceId, "replica_num"),
					resource.TestCheckResourceAttrSet(dataSourceWithInstanceId, "initial_node_num"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussSolutionTemplateSetting_base(name string) string {
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

func testDataSourceGaussdbOpengaussSolutionTemplateSetting_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_opengauss_solution_template_setting" "test" {
  solution = "single"
}

data "huaweicloud_gaussdb_opengauss_solution_template_setting" "test_instance_id" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}
`, testDataSourceGaussdbOpengaussSolutionTemplateSetting_base(name))
}

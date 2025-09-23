package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussTopIoTraffics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_top_io_traffics.test"
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
				Config: testDataSourceGaussdbOpengaussTopIoTraffics_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.thread_id"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.thread_type"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.disk_read_rate"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.disk_write_rate"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.session_id"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.unique_sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.client_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "top_io_infos.0.sql_start"),
					resource.TestCheckOutput("top_io_num_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_condition_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussTopIoTraffics_base(name string) string {
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

data "huaweicloud_gaussdb_opengauss_instance_nodes" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}
`, common.TestBaseNetwork(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbOpengaussTopIoTraffics_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_opengauss_top_io_traffics" "test" {
  instance_id  = huaweicloud_gaussdb_opengauss_instance.test.id
  node_id      = data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].id
  component_id = [for v in data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].components : v.id if v.type == "DN"][0]
}

data "huaweicloud_gaussdb_opengauss_top_io_traffics" "top_io_num_filter" {
  instance_id  = huaweicloud_gaussdb_opengauss_instance.test.id
  node_id      = data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].id
  component_id = [for v in data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].components : v.id if v.type == "DN"][0]
  top_io_num   = 0
}
output "top_io_num_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_top_io_traffics.top_io_num_filter.top_io_infos) > 0
}

data "huaweicloud_gaussdb_opengauss_top_io_traffics" "sort_condition_filter" {
  instance_id    = huaweicloud_gaussdb_opengauss_instance.test.id
  node_id        = data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].id
  component_id   = [for v in data.huaweicloud_gaussdb_opengauss_instance_nodes.test.nodes[0].components : v.id if v.type == "DN"][0]
  sort_condition = "read"
}
output "sort_condition_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_top_io_traffics.sort_condition_filter.top_io_infos) > 0
}
`, testDataSourceGaussdbOpengaussTopIoTraffics_base(name))
}

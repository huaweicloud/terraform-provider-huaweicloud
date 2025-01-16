package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceOpenGaussTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_tasks.test"
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
				Config: testDataSourceOpenGaussTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.ended_at"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOpenGaussTasks_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

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
  flavor                = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name                  = "%[2]s"
  password              = "Huangwei!120521"
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

func testDataSourceOpenGaussTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_tasks" "test" {
  depends_on = [huaweicloud_gaussdb_opengauss_instance.test]
}

locals {
  status = "Completed"
}
data "huaweicloud_gaussdb_opengauss_tasks" "status_filter" {
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test
  ]

  status = "Completed"
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_tasks.status_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_tasks.status_filter.tasks[*].status : v == local.status]
  )
}

locals {
  name = data.huaweicloud_gaussdb_opengauss_tasks.test.tasks[0].name
}
data "huaweicloud_gaussdb_opengauss_tasks" "name_filter" {
  depends_on = [
    data.huaweicloud_gaussdb_opengauss_tasks.test
  ]

  name = data.huaweicloud_gaussdb_opengauss_tasks.test.tasks[0].name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_tasks.name_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_tasks.name_filter.tasks[*].name : v == local.name]
  )
}

locals {
  start_time = data.huaweicloud_gaussdb_opengauss_tasks.test.tasks[0].created_at
  end_time   = data.huaweicloud_gaussdb_opengauss_tasks.test.tasks[0].ended_at
}
data "huaweicloud_gaussdb_opengauss_tasks" "time_filter" {
  depends_on = [
    data.huaweicloud_gaussdb_opengauss_tasks.test
  ]

  start_time = data.huaweicloud_gaussdb_opengauss_tasks.test.tasks[0].created_at
  end_time   = data.huaweicloud_gaussdb_opengauss_tasks.test.tasks[0].ended_at
}
output "time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_tasks.time_filter.tasks) > 0
}
`, testDataSourceOpenGaussTasks_base(name))
}

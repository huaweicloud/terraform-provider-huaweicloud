package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbMysqlScheduledTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_scheduled_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlScheduledTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.job_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.job_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.datastore_type"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("job_id_filter_is_useful", "true"),
					resource.TestCheckOutput("job_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlScheduledTasks_base(name string) string {
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

resource "huaweicloud_gaussdb_mysql_instance_restart" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  node_id     = huaweicloud_gaussdb_mysql_instance.test.nodes[0].id
  delay       = true
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceGaussdbMysqlScheduledTasks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_mysql_scheduled_tasks" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_instance_restart.test]
}

locals {
  status = "Pending"
}
data "huaweicloud_gaussdb_mysql_scheduled_tasks" "status_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_instance_restart.test]

  status = "Pending"
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_scheduled_tasks.status_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_scheduled_tasks.status_filter.tasks[*].job_status : v == local.status]
  )
}

locals {
  job_id = data.huaweicloud_gaussdb_mysql_scheduled_tasks.test.tasks[0].job_id
}
data "huaweicloud_gaussdb_mysql_scheduled_tasks" "job_id_filter" {
  depends_on = [
    huaweicloud_gaussdb_mysql_instance_restart.test,
    data.huaweicloud_gaussdb_mysql_scheduled_tasks.test
  ]

  job_id = data.huaweicloud_gaussdb_mysql_scheduled_tasks.test.tasks[0].job_id
}
output "job_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_scheduled_tasks.job_id_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_scheduled_tasks.job_id_filter.tasks[*].job_id : v == local.job_id]
  )
}

locals {
  job_name = "REBOOT_INSTANCE"
}
data "huaweicloud_gaussdb_mysql_scheduled_tasks" "job_name_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_instance_restart.test]

  job_name = "REBOOT_INSTANCE"
}
output "job_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_scheduled_tasks.job_name_filter.tasks) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_scheduled_tasks.job_name_filter.tasks[*].job_name : v == local.job_name]
  )
}
`, testDataSourceGaussdbMysqlScheduledTasks_base(name))
}

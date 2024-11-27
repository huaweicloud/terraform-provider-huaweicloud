package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussDBMysqlInstantTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_instant_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussDBMysqlInstantTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.ended_time"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("job_id_filter_is_useful", "true"),
					resource.TestCheckOutput("job_name_filter_is_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussDBMysqlInstantTasks_base(name string) string {
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
  read_replicas            = 4
  port                     = 8888
  ssl_option               = "false"
  description              = "test_description"
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceDataSourceGaussDBMysqlInstantTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_mysql_instant_tasks" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_instance.test]
}

locals {
  status = "Completed"
}
data "huaweicloud_gaussdb_mysql_instant_tasks" "status_filter" {
  depends_on = [
    huaweicloud_gaussdb_mysql_instance.test
  ]

  status = "Completed"
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_instant_tasks.status_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_instant_tasks.status_filter.jobs[*].status : v == local.status]
  )
}

locals {
  job_id = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].job_id
}
data "huaweicloud_gaussdb_mysql_instant_tasks" "job_id_filter" {
  depends_on = [
    huaweicloud_gaussdb_mysql_instance.test,
    data.huaweicloud_gaussdb_mysql_instant_tasks.test
  ]

  job_id = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].job_id
}
output "job_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_instant_tasks.job_id_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_instant_tasks.job_id_filter.jobs[*].job_id : v == local.job_id]
  )
}

locals {
  job_name = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].job_name
}
data "huaweicloud_gaussdb_mysql_instant_tasks" "job_name_filter" {
  depends_on = [
    huaweicloud_gaussdb_mysql_instance.test,
    data.huaweicloud_gaussdb_mysql_instant_tasks.test
  ]

  job_name = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].job_name
}
output "job_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_instant_tasks.job_name_filter.jobs) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_instant_tasks.job_name_filter.jobs[*].job_name : v == local.job_name]
  )
}

locals {
  start_time = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].created_time
  end_time   = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].ended_time
}
data "huaweicloud_gaussdb_mysql_instant_tasks" "time_filter" {
  depends_on = [
    huaweicloud_gaussdb_mysql_instance.test,
    data.huaweicloud_gaussdb_mysql_instant_tasks.test
  ]

  start_time = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].created_time
  end_time   = data.huaweicloud_gaussdb_mysql_instant_tasks.test.jobs[0].ended_time
}
output "time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_instant_tasks.time_filter.jobs) > 0
}
`, testDataSourceDataSourceGaussDBMysqlInstantTasks_base(name))
}

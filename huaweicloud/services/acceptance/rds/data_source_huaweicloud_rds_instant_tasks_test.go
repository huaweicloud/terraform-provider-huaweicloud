package rds

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsInstantTasks_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_rds_instant_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsInstantTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "actions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tasks.0.end_time"),
					resource.TestCheckOutput("task_id_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsInstantTasks_base(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceRdsInstantTasks_basic(name string) string {
	now := time.Now()
	startTime := now.Add(-100 * time.Hour).UnixMilli()
	endTime := now.Add(time.Hour).UnixMilli()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_instant_tasks" "test" {
  depends_on = [huaweicloud_rds_instance.test] 
}

data "huaweicloud_rds_instant_tasks" "task_id_filter" {
  depends_on = [huaweicloud_rds_instance.test] 

  task_id = data.huaweicloud_rds_instant_tasks.test.tasks.0.id
}
locals {
  task_id = data.huaweicloud_rds_instant_tasks.test.tasks.0.id
}
output "task_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_instant_tasks.task_id_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_rds_instant_tasks.task_id_filter.tasks[*].id : v == local.task_id]
  )
}

data "huaweicloud_rds_instant_tasks" "instance_id_filter" {
  depends_on = [huaweicloud_rds_instance.test] 

  instance_id = data.huaweicloud_rds_instant_tasks.test.tasks.0.instance_id
}
locals {
  instance_id = data.huaweicloud_rds_instant_tasks.test.tasks.0.instance_id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_instant_tasks.instance_id_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_rds_instant_tasks.instance_id_filter.tasks[*].instance_id : v == local.instance_id]
  )
}

data "huaweicloud_rds_instant_tasks" "name_filter" {
  depends_on = [huaweicloud_rds_instance.test] 

  name = data.huaweicloud_rds_instant_tasks.test.tasks.0.name
}
locals {
  name = data.huaweicloud_rds_instant_tasks.test.tasks.0.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_rds_instant_tasks.name_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_rds_instant_tasks.name_filter.tasks[*].name : v == local.name]
  )
}

data "huaweicloud_rds_instant_tasks" "status_filter" {
  depends_on = [huaweicloud_rds_instance.test] 

  status = data.huaweicloud_rds_instant_tasks.test.tasks.0.status
}
locals {
  status = data.huaweicloud_rds_instant_tasks.test.tasks.0.status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_rds_instant_tasks.status_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_rds_instant_tasks.status_filter.tasks[*].status : v == local.status]
  )
}

data "huaweicloud_rds_instant_tasks" "time_filter" {
  depends_on = [huaweicloud_rds_instance.test] 

  start_time = "%[2]d"
  end_time   = "%[3]d"
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_rds_instant_tasks.time_filter.tasks) > 0
}
`, testDataSourceRdsInstantTasks_base(name), startTime, endTime)
}

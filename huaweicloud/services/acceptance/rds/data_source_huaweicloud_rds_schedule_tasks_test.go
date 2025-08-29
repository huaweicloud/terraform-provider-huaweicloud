package rds

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsScheduleTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_schedule_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsScheduleTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.datastore_type"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.order"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.target_config.#"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_tasks.0.end_time"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_name_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsScheduleTasks_base(name string) string {
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

resource "huaweicloud_rds_instance_restart" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  delay       = true
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceDataSourceRdsScheduleTasks_basic(name string) string {
	now := time.Now()
	startTime := now.Add(-100 * time.Hour).UnixMilli()
	endTime := now.Add(time.Hour).UnixMilli()
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_schedule_tasks" "test" {
  depends_on = [huaweicloud_rds_instance_restart.test]
}

data "huaweicloud_rds_schedule_tasks" "instance_id_filter" {
  depends_on = [huaweicloud_rds_instance_restart.test]

  instance_id = huaweicloud_rds_instance.test.id
}
locals {
  instance_id = huaweicloud_rds_instance.test.id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_schedule_tasks.instance_id_filter.schedule_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_rds_schedule_tasks.instance_id_filter.schedule_tasks : v.instance_id == local.instance_id]
  )
}

data "huaweicloud_rds_schedule_tasks" "instance_name_filter" {
  depends_on = [huaweicloud_rds_instance_restart.test]

  instance_name = "%[2]s"
}
output "instance_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_schedule_tasks.instance_name_filter.schedule_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_rds_schedule_tasks.instance_name_filter.schedule_tasks : v.instance_name == "%[2]s"]
  )
}

data "huaweicloud_rds_schedule_tasks" "status_filter" {
  depends_on = [huaweicloud_rds_instance_restart.test]

  status = data.huaweicloud_rds_schedule_tasks.test.schedule_tasks[0].status
}
locals {
  status = data.huaweicloud_rds_schedule_tasks.test.schedule_tasks[0].status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_rds_schedule_tasks.status_filter.schedule_tasks) > 0 && alltrue(
  [for v in data.huaweicloud_rds_schedule_tasks.status_filter.schedule_tasks : v.status == local.status]
  )
}

data "huaweicloud_rds_instant_tasks" "time_filter" {
  depends_on = [huaweicloud_rds_instance_restart.test]

  start_time = "%[3]d"
  end_time   = "%[4]d"
}
output "time_filter_is_useful" {
  value = length(data.huaweicloud_rds_instant_tasks.time_filter.tasks) > 0
}
`, testDataSourceDataSourceRdsScheduleTasks_base(name), name, startTime, endTime)
}

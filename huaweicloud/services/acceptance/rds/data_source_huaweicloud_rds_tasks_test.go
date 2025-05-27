package rds

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	startTime := time.Now().AddDate(0, 0, -10).UnixMilli()
	endTime := time.Now().AddDate(0, 0, 10).UnixMilli()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsTasks_basic(rName, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance.0.name"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsTasks_base(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup" "secgroup_test" {
  name = "secgroup_test"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = "rds.mssql.spec.se.s6.large.2"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
  charging_mode     = "postPaid"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type     = "SQLServer"
    version  = "2022_SE"
    password = "Terraform145!"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
`, testAccRdsInstance_base(name), name)
}

func testDataSourceRdsTasks_basic(name string, startTime, endTime int64) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_tasks" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  start_time  = "%[2]d"
}

data "huaweicloud_rds_tasks" "end_time_filter" {
  instance_id = huaweicloud_rds_instance.test.id
  start_time  = "%[2]d"
  end_time    = "%[3]d"
}

output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_rds_tasks.end_time_filter.jobs) > 0
}
`, testDataSourceRdsTasks_base(name), startTime, endTime)
}

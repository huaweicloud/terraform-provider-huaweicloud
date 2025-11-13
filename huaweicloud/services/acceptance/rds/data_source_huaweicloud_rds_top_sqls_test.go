package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsTopSqls_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_top_sqls.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsTopSqls_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_cpu_time_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_cpu_time_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_cpu_time_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_cpu_time_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_cpu_time_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_cpu_time_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_cpu_time_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_cpu_time_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_logical_reads_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_logical_reads_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_logical_reads_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_logical_reads_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_duration_time_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_duration_time_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_duration_time_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_duration_time_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_rows_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_rows_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_rows_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_rows_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_logical_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_logical_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_logical_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"avg_logical_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_duration_time_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_duration_time_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_duration_time_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_duration_time_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_rows_top_values.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_rows_top_values.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_rows_top_values.0.data_type"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"total_rows_top_values.0.value"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.#"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.id"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_logical_reads"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.query"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.execution_count"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_duration_time_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_rows"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_physical_reads"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_duration_time_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_rows_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_logical_reads_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_logical_reads_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_logical_write_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.last_execution_time"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_cpu_time"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_rows_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_logical_write"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_physical_reads"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_cpu_time"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_physical_reads_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_cpu_time_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.statement"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_duration_time"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_rows"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_physical_reads_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_duration_time"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.avg_logical_write"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_logical_write_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.db_name"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.execution_count_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_cpu_time_percent"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_rds_top_sqls.test",
						"list.0.total_logical_reads"),
				),
			},
		},
	})
}

func testDataSourceRdsTopSqls_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[1]s"
  flavor            = "rds.mssql.spec.x1.se.xlarge.4"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, name)
}

func testDataSourceRdsTopSqls_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_top_sqls" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}

data "huaweicloud_rds_top_sqls" "sort_asc_filter" {
  instance_id = huaweicloud_rds_instance.test.id
  sort_key    = "total_cpu_time"
  sort_dir    = "asc"
}
locals {
  sort_asc_filter = data.huaweicloud_rds_top_sqls.sort_asc_filter
}
output "sort_asc_filter_is_useful" {
  value = length(local.sort_asc_filter) > 0 && (local.sort_asc_filter.list[0].total_cpu_time <=
    local.sort_asc_filter.list[1].total_cpu_time)
}

data "huaweicloud_rds_top_sqls" "sort_desc_filter" {
  instance_id = huaweicloud_rds_instance.test.id
  sort_key    = "total_cpu_time"
  sort_dir    = "desc"
}
locals {
  sort_desc_filter = data.huaweicloud_rds_top_sqls.sort_desc_filter
}
output "sort_desc_filter_is_useful" {
  value = length(local.sort_desc_filter) > 0 && (local.sort_desc_filter.list[0].total_cpu_time >=
    local.sort_desc_filter.list[1].total_cpu_time)
}

locals {
  limit = 5
}
data "huaweicloud_rds_top_sqls" "limit_filter" {
  instance_id = huaweicloud_rds_instance.test.id
  limit       = 5
}
output "limit_filter_is_useful" {
  value = length(data.huaweicloud_rds_top_sqls.limit_filter) <= 5
}

data "huaweicloud_rds_top_sqls" "statement_filter" {
  instance_id = huaweicloud_rds_instance.test.id
  statement   = "select"
}
output "statement_filter_is_useful" {
  value = length(data.huaweicloud_rds_top_sqls.statement_filter) > 0
}
`, testDataSourceRdsTopSqls_base(name))
}

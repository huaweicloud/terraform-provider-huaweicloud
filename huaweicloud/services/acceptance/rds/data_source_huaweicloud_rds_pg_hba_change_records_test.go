package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourcePgHbaChangeRecords_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_pg_hba_change_records.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgHbaChangeRecords_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.#"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.status"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.time"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.before_confs.#"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.before_confs.0.type"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.before_confs.0.database"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.before_confs.0.user"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.before_confs.0.address"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.before_confs.0.method"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.before_confs.0.priority"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.after_confs.#"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.after_confs.0.type"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.after_confs.0.database"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.after_confs.0.user"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.after_confs.0.address"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.after_confs.0.method"),
					resource.TestCheckResourceAttrSet(rName, "pg_hba_change_records.0.after_confs.0.priority"),
					resource.TestCheckOutput("start_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePgHbaChangeRecords_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "PostgreSQL"
  db_version    = "12"
  instance_mode = "single"
  group_type    = "dedicated"
}
resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  db {
    type    = "PostgreSQL"
    version = "12"
  }
    
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_pg_hba" "test" {
  instance_id = huaweicloud_rds_instance.test.id

  host_based_authentications {
    type     = "host"
    database = "all"
    user     = "all"
    address  = "0.0.0.0/0"
    method   = "md5"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourcePgHbaChangeRecords_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_pg_hba_change_records" "test" {
  depends_on = [huaweicloud_rds_pg_hba.test]

  instance_id = huaweicloud_rds_instance.test.id
}

data "huaweicloud_rds_pg_hba_change_records" "start_time_filter" {
  depends_on = [huaweicloud_rds_pg_hba.test]

  instance_id = huaweicloud_rds_instance.test.id
  start_time  = data.huaweicloud_rds_pg_hba_change_records.test.pg_hba_change_records[0].time
}

locals {
  start_time = data.huaweicloud_rds_pg_hba_change_records.test.pg_hba_change_records[0].time
}

output "start_time_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_hba_change_records.start_time_filter.pg_hba_change_records) > 0
}

data "huaweicloud_rds_pg_hba_change_records" "end_time_filter" {
  instance_id = huaweicloud_rds_instance.test.id
  end_time    = data.huaweicloud_rds_pg_hba_change_records.test.pg_hba_change_records[0].time
}

locals {
  end_time = data.huaweicloud_rds_pg_hba_change_records.test.pg_hba_change_records[0].time
}

output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_hba_change_records.end_time_filter.pg_hba_change_records) == 0
}
`, testAccDatasourcePgHbaChangeRecords_base(name))
}

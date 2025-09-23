package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussBackups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_backups.test"
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
				Config: testDataSourceGaussdbOpengaussBackups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "backups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "backups.0.datastore.0.version"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("backup_id_filter_is_useful", "true"),
					resource.TestCheckOutput("backup_type_filter_is_useful", "true"),
					resource.TestCheckOutput("begin_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussBackups_base(name string) string {
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

resource "huaweicloud_gaussdb_opengauss_backup" "backup_1" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  name        = "%[2]s_backup_1"
  description = "test description"
}

resource "huaweicloud_gaussdb_opengauss_backup" "backup_2" {
  depends_on = [huaweicloud_gaussdb_opengauss_backup.backup_1]

  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  name        = "%[2]s_backup_2"
  description = "test description"
}
`, common.TestBaseNetwork(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbOpengaussBackups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_backups" "test" {
  depends_on = [huaweicloud_gaussdb_opengauss_backup.backup_2]
}

locals {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}
data "huaweicloud_gaussdb_opengauss_backups" "instance_id_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_backup.backup_2]

  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_backups.instance_id_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_backups.instance_id_filter.backups[*].instance_id : v == local.instance_id]
  )
}

locals {
  backup_id = huaweicloud_gaussdb_opengauss_backup.backup_2.id
}
data "huaweicloud_gaussdb_opengauss_backups" "backup_id_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_backup.backup_2]

  backup_id  = huaweicloud_gaussdb_opengauss_backup.backup_2.id
}
output "backup_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_backups.backup_id_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_backups.backup_id_filter.backups[*].id : v == local.backup_id]
  )
}

locals {
  backup_type = "manual"
}
data "huaweicloud_gaussdb_opengauss_backups" "backup_type_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_backup.backup_2]

  backup_type = "manual"
}
output "backup_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_backups.backup_type_filter.backups) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_backups.backup_type_filter.backups[*].type : v == local.backup_type]
  )
}

locals {
  begin_time = huaweicloud_gaussdb_opengauss_backup.backup_1.begin_time
}
data "huaweicloud_gaussdb_opengauss_backups" "begin_time_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_backup.backup_2]

  begin_time = huaweicloud_gaussdb_opengauss_backup.backup_1.begin_time
}
output "begin_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_backups.begin_time_filter.backups) > 0
}

locals {
  end_time = huaweicloud_gaussdb_opengauss_backup.backup_2.end_time
}
data "huaweicloud_gaussdb_opengauss_backups" "end_time_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_backup.backup_2]

  end_time = huaweicloud_gaussdb_opengauss_backup.backup_2.end_time
}
output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_backups.end_time_filter.backups) > 0
}
`, testDataSourceGaussdbOpengaussBackups_base(name))
}

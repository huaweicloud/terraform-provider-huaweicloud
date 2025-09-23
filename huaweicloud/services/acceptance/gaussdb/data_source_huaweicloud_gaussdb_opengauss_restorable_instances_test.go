package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussRestorableInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_restorable_instances.test"
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
				Config: testDataSourceGaussdbOpengaussRestorableInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.data_volume_size"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckOutput("restore_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussRestorableInstances_base(name string) string {
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

  count = 2

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name              = "%[2]s_${count.index}"
  password          = "Huangwei!120521"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

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

resource "huaweicloud_gaussdb_opengauss_backup" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test[0].id
  name        = "%[2]s_backup"
}

data "huaweicloud_gaussdb_opengauss_restore_time_ranges" "test" {
  instance_id = huaweicloud_gaussdb_opengauss_instance.test[0].id
  date        = split("T", huaweicloud_gaussdb_opengauss_backup.test.end_time)[0]
}
`, common.TestBaseNetwork(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbOpengaussRestorableInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_opengauss_restorable_instances" "test" {
  source_instance_id = huaweicloud_gaussdb_opengauss_instance.test[0].id
  backup_id          = huaweicloud_gaussdb_opengauss_backup.test.id
}

locals {
  restore_time = data.huaweicloud_gaussdb_opengauss_restore_time_ranges.test.restore_time[0].end_time
}
data "huaweicloud_gaussdb_opengauss_restorable_instances" "restore_time_filter" {
  source_instance_id = huaweicloud_gaussdb_opengauss_instance.test[0].id
  backup_id          = huaweicloud_gaussdb_opengauss_backup.test.id
  restore_time       = data.huaweicloud_gaussdb_opengauss_restore_time_ranges.test.restore_time[0].end_time
}
output "restore_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_restorable_instances.restore_time_filter.instances) > 0
}
`, testDataSourceGaussdbOpengaussRestorableInstances_base(name))
}

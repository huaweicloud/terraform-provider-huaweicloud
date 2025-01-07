package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbOpengaussRecyclingInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_recycling_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussRecyclingInstances_create_instance(rName),
			},
			{
				Config: testDataSourceGaussdbOpengaussRecyclingInstances_delete_instance(),
			},
			{
				Config: testDataSourceGaussdbOpengaussRecyclingInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.ha_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine_version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.pay_model"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.recycle_backup_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.data_vip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.recycle_status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.deleted_at"),
					resource.TestCheckOutput("instance_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussRecyclingInstances_create_instance(rName string) string {
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

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name              = "%[2]s"
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
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGaussdbOpengaussRecyclingInstances_delete_instance() string {
	return `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}`
}

func testDataSourceGaussdbOpengaussRecyclingInstances_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_recycling_instances" "test" {}

locals {
  instance_name = "%[1]s"
}
data "huaweicloud_gaussdb_opengauss_recycling_instances" "instance_name_filter" {
  instance_name = "%[1]s"
}
output "instance_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_recycling_instances.instance_name_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_recycling_instances.instance_name_filter.instances[*] : v.name == local.instance_name]
  )
}
`, rName)
}

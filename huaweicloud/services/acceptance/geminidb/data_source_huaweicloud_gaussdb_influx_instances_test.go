package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGaussdbInfluxInstancesDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_influx_instances.test"
	rName := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbInfluxInstances_basic(rName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.engine"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.datastore.0.patch_available"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.engine"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.db_user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.security_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.backup_strategy.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.backup_strategy.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.backup_strategy.0.keep_days"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.pay_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.maintain_begin"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.maintain_end"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.volume.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.volume.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.volume.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.groups.0.nodes.0.support_reduce"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.actions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.lb_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.lb_port"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.updated_at"),

					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("mode_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbInfluxInstances_base(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "influxdb"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_gaussdb_influx_instance" "test" {
  name        = "%[2]s"
  password    = "%[3]s"
  flavor      = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  volume_size = 100
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 3

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(name), name, password)
}

func testDataSourceGaussdbInfluxInstances_basic(name, password string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_influx_instances" "test" {
  depends_on = [huaweicloud_gaussdb_influx_instance.test]
}

data "huaweicloud_gaussdb_influx_instances" "id_filter" {
  depends_on  = [huaweicloud_gaussdb_influx_instance.test]

  instance_id = huaweicloud_gaussdb_influx_instance.test.id
}

locals {
  id = huaweicloud_gaussdb_influx_instance.test.id
}

output "id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_influx_instances.id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_influx_instances.id_filter.instances[*].id : v == local.id]
  )
}

data "huaweicloud_gaussdb_influx_instances" "name_filter" {
  depends_on  = [huaweicloud_gaussdb_influx_instance.test]

  name = huaweicloud_gaussdb_influx_instance.test.name
}

locals {
  name = huaweicloud_gaussdb_influx_instance.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_influx_instances.name_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_influx_instances.name_filter.instances[*].name : v == local.name]
  )
}

data "huaweicloud_gaussdb_influx_instances" "mode_filter" {
  depends_on  = [huaweicloud_gaussdb_influx_instance.test]

  mode = huaweicloud_gaussdb_influx_instance.test.mode
}

locals {
  mode = huaweicloud_gaussdb_influx_instance.test.mode
}

output "mode_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_influx_instances.mode_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_influx_instances.mode_filter.instances[*].mode : v == local.mode]
  )
}

data "huaweicloud_gaussdb_influx_instances" "vpc_id_filter" {
  depends_on  = [huaweicloud_gaussdb_influx_instance.test]

  vpc_id = huaweicloud_gaussdb_influx_instance.test.vpc_id
}

locals {
  vpc_id = huaweicloud_gaussdb_influx_instance.test.vpc_id
}

output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_influx_instances.vpc_id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_influx_instances.vpc_id_filter.instances[*].vpc_id : v == local.vpc_id]
  )
}

data "huaweicloud_gaussdb_influx_instances" "subnet_id_filter" {
  depends_on  = [huaweicloud_gaussdb_influx_instance.test]

  subnet_id = huaweicloud_gaussdb_influx_instance.test.subnet_id
}

locals {
  subnet_id = huaweicloud_gaussdb_influx_instance.test.subnet_id
}

output "subnet_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_influx_instances.subnet_id_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_influx_instances.subnet_id_filter.instances[*].subnet_id : v == local.subnet_id]
  )
}
`, testDataSourceGaussdbInfluxInstances_base(name, password))
}

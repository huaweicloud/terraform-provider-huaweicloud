package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccRdsInstancesDataSource_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_instances.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instances.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.region"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.type"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.status"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.created"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.security_group_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.flavor"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.time_zone"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.ssl_enable"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.tags.%"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.public_ips.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.private_ips.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.private_ips.0"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.fixed_ip"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.volume.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.volume.0.type"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.volume.0.size"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.db.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.db.0.type"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.db.0.version"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.db.0.user_name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.backup_strategy.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.backup_strategy.0.start_time"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.backup_strategy.0.start_time"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.nodes.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.nodes.0.id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.nodes.0.name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.nodes.0.role"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.nodes.0.status"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.nodes.0.availability_zone"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.availability_zone.#"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.availability_zone.0"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("datastore_type_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("group_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccRdsInstanceDataSource_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    password = "Huangwei!120521"
    type     = "PostgreSQL"
    version  = "12"
    port     = 8635
  }
  volume {
    type = "CLOUDSSD"
    size = 50
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}

resource "huaweicloud_rds_instance" "flexus_test" {
  name              = "%[2]s"
  flavor            = "rds.mysql.y1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"
  is_flexus         = true

  db {
    type    = "MySQL"
    version = "8.0"
    port    = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}

`, common.TestBaseNetwork(rName), rName)
}

func testAccRdsInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_instances" "test" {
  depends_on = [
    huaweicloud_rds_instance.test,
    huaweicloud_rds_instance.flexus_test,
  ]
}

locals {
  name = "%[2]s"
}
data "huaweicloud_rds_instances" "name_filter" {
  depends_on = [
    huaweicloud_rds_instance.test,
    huaweicloud_rds_instance.flexus_test,
  ]

  name = "%[2]s"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_rds_instances.name_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_instances.name_filter.instances[*].name : v == local.name]
  )
}

locals {
  type = "Single"
}
data "huaweicloud_rds_instances" "type_filter" {
  depends_on = [
    huaweicloud_rds_instance.test,
    huaweicloud_rds_instance.flexus_test,
  ]

  type = "Single"
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_rds_instances.type_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_instances.type_filter.instances[*].type : v == local.type]
  )
}

locals {
  datastore_type = "PostgreSQL"
}
data "huaweicloud_rds_instances" "datastore_type_filter" {
  depends_on = [
    huaweicloud_rds_instance.test,
    huaweicloud_rds_instance.flexus_test,
  ]

  datastore_type = "PostgreSQL"
}
output "datastore_type_filter_is_useful" {
  value = length(data.huaweicloud_rds_instances.datastore_type_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_instances.datastore_type_filter.instances[*].db[0].type : v == local.datastore_type]
  )
}

locals {
  vpc_id = huaweicloud_vpc.test.id
}
data "huaweicloud_rds_instances" "vpc_id_filter" {
  depends_on = [
    huaweicloud_rds_instance.test,
    huaweicloud_rds_instance.flexus_test,
  ]

  vpc_id = huaweicloud_vpc.test.id
}
output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_instances.vpc_id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_instances.vpc_id_filter.instances[*].vpc_id : v == local.vpc_id]
  )
}

locals {
  subnet_id = huaweicloud_vpc_subnet.test.id
}
data "huaweicloud_rds_instances" "subnet_id_filter" {
  depends_on = [
    huaweicloud_rds_instance.test,
    huaweicloud_rds_instance.flexus_test,
  ]

  subnet_id = huaweicloud_vpc_subnet.test.id
}
output "subnet_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_instances.subnet_id_filter.instances) > 0 && alltrue(
  [for v in data.huaweicloud_rds_instances.subnet_id_filter.instances[*].subnet_id : v == local.subnet_id]
  )
}

locals {
  group_type = "flexus"
}
data "huaweicloud_rds_instances" "group_type_filter" {
  depends_on = [
    huaweicloud_rds_instance.test,
    huaweicloud_rds_instance.flexus_test,
  ]

  group_type = "flexus"
}
output "group_type_filter_is_useful" {
  value = length(data.huaweicloud_rds_instances.group_type_filter.instances) > 0
}
`, testAccRdsInstanceDataSource_base(rName), rName)
}

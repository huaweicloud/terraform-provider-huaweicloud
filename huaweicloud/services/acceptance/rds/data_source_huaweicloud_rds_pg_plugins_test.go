package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPgPlugins_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_pg_plugins.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgPlugins_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.name"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.version"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.created"),
					resource.TestCheckResourceAttrSet(rName, "plugins.0.description"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("version_filter_is_useful", "true"),

					resource.TestCheckOutput("created_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePgPlugins_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "12"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}

resource "huaweicloud_rds_pg_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "%[2]s"
  owner         = "root"
  character_set = "UTF8"
  template      = "template1"
  lc_collate    = "en_US.UTF-8"
  lc_ctype      = "en_US.UTF-8"
}
`, testAccRdsInstance_base(), name)
}

func testAccDatasourcePgPlugins_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_pg_plugins" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  database_name = huaweicloud_rds_pg_database.test.name
}

data "huaweicloud_rds_pg_plugins" "name_filter" {
  instance_id   = huaweicloud_rds_instance.test.id
  database_name = huaweicloud_rds_pg_database.test.name
  name          = data.huaweicloud_rds_pg_plugins.test.plugins[0].name
}
locals {
  name = data.huaweicloud_rds_pg_plugins.test.plugins[0].name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_plugins.name_filter.plugins) > 0 && alltrue(
    [for v in data.huaweicloud_rds_pg_plugins.name_filter.plugins[*].name : v == local.name]
  )  
}

data "huaweicloud_rds_pg_plugins" "version_filter" {
  instance_id   = huaweicloud_rds_instance.test.id
  database_name = huaweicloud_rds_pg_database.test.name
  version       = data.huaweicloud_rds_pg_plugins.test.plugins[0].version

}
locals {
  version = data.huaweicloud_rds_pg_plugins.test.plugins[0].version
}
output "version_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_plugins.version_filter.plugins) > 0 && alltrue(
    [for v in data.huaweicloud_rds_pg_plugins.version_filter.plugins[*].version : v == local.version]
  )  
}

data "huaweicloud_rds_pg_plugins" "created_filter" {
  instance_id   = huaweicloud_rds_instance.test.id
  database_name = huaweicloud_rds_pg_database.test.name
  created       = data.huaweicloud_rds_pg_plugins.test.plugins[0].created

}
locals {
  created = data.huaweicloud_rds_pg_plugins.test.plugins[0].created
}
output "created_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_plugins.created_filter.plugins) > 0 && alltrue(
    [for v in data.huaweicloud_rds_pg_plugins.created_filter.plugins[*].created : v == local.created]
  )  
}
`, testAccDatasourcePgPlugins_base(name))
}

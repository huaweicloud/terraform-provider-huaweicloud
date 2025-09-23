package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccOpenGaussSchemasDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_schemas.test"
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
				Config: testAccOpenGaussSchemasDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "database_schemas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "database_schemas.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "database_schemas.0.owner"),
				),
			},
		},
	})
}

func testAccOpenGaussSchemasDataSource_base(rName string) string {
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

resource "huaweicloud_gaussdb_opengauss_database" "test" {
  instance_id   = huaweicloud_gaussdb_opengauss_instance.test.id
  name          = "%[2]s"
  character_set = "UTF8"
  owner         = "root"
  template      = "template0"
  lc_collate    = "C"
  lc_ctype      = "C"
}

resource "huaweicloud_gaussdb_opengauss_schema" "test" {
  depends_on = [huaweicloud_gaussdb_opengauss_database.test]

  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  db_name     = huaweicloud_gaussdb_opengauss_database.test.name
  name        = "%[2]s"
  owner       = "root"
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccOpenGaussSchemasDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_schemas" "test" {
  depends_on = [huaweicloud_gaussdb_opengauss_schema.test]

  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  db_name     = huaweicloud_gaussdb_opengauss_database.test.name
}
`, testAccOpenGaussSchemasDataSource_base(rName))
}

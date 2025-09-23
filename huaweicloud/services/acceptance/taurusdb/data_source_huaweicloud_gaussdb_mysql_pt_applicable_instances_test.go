package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbMysqlPtApplicableInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_pt_applicable_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlPtApplicableInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_name"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlPtApplicableInstances_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_mysql_flavors" "test" {
  engine                 = "gaussdb-mysql"
  version                = "8.0"
  availability_zone_mode = "multi"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                     = "%[2]s"
  password                 = "Test@12345678"
  flavor                   = data.huaweicloud_gaussdb_mysql_flavors.test.flavors[0].name
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  enterprise_project_id    = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
  read_replicas            = 2

  parameters {
    name  = "auto_increment_increment"
    value = "100"
  }

  lifecycle {
    ignore_changes = [
      parameters,
    ]
  }
}

resource "huaweicloud_gaussdb_mysql_parameter_template" "test" {
  name              = "%[2]s"
  datastore_engine  = "gaussdb-mysql"
  datastore_version = "8.0"

  parameter_values = {
    auto_increment_increment = "50"
    character_set_server     = "gbk"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceGaussdbMysqlPtApplicableInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_mysql_pt_applicable_instances" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_instance.test]

  configuration_id = huaweicloud_gaussdb_mysql_parameter_template.test.id
}
`, testDataSourceGaussdbMysqlPtApplicableInstances_base(name))
}

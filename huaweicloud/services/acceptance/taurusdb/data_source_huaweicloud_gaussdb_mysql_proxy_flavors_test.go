package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGaussdbMysqlProxyFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_proxy_flavors.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlProxyFlavors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.az_status.%"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlProxyFlavors_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_mysql_flavors" "test" {
  engine  = "gaussdb-mysql"
  version = "8.0"
}

resource "huaweicloud_gaussdb_mysql_instance" "test" {
  name                     = "%s"
  password                 = "Test@12345678"
  flavor                   = data.huaweicloud_gaussdb_mysql_flavors.test.flavors[0].name
  vpc_id                   = huaweicloud_vpc.test.id
  subnet_id                = huaweicloud_vpc_subnet.test.id
  security_group_id        = huaweicloud_networking_secgroup.test.id
  availability_zone_mode   = "multi"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  read_replicas            = 2
  enterprise_project_id    = "0"
}
`, common.TestBaseNetwork(rName), rName)
}

func testDataSourceGaussdbMysqlProxyFlavors_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_gaussdb_mysql_proxy_flavors" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}
`, testDataSourceGaussdbMysqlProxyFlavors_base(name))
}

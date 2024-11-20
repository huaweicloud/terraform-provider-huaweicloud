package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccRdsMysqlProxyFlavorsDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_mysql_proxy_flavors.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsMysqlProxyFlavors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.memory"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavor_groups.0.flavors.0.az_status.%"),
				),
			},
		},
	})
}

func testDataSourceRdsMysqlProxyFlavors_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  db {
    password = "test_1234"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testDataSourceRdsMysqlProxyFlavors_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_mysql_proxy_flavors" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testDataSourceRdsMysqlProxyFlavors_base(name))
}

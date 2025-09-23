package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsDatabases_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dbss_rds_databases.test"
		name           = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsDatabases_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.is_supported"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "databases.0.enterprise_project_id"),
				),
			},
		},
	})
}

func testDataSourceRdsDatabases_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 2
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 100
  }
}

data "huaweicloud_dbss_rds_databases" "test" {
  type = "MYSQL"

  depends_on = [huaweicloud_rds_instance.test]
}
`, name)
}

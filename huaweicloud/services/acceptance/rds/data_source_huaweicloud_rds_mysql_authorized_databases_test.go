package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsMysqlAuthorizedDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_mysql_authorized_databases.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsMysqlAuthorizedDatabases_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.readonly"),
				),
			},
		},
	})
}

func testDataSourceRdsMysqlAuthorizedDatabases_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mysql.x1.large.2"
  vpc_id            = data.huaweicloud_vpc.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  charging_mode     = "postPaid"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type     = "MySQL"
    version  = "8.0"
    password = "Terraform145!"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "test"
  character_set = "utf8"
}

resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "test"
  password    = "Terraform145!"
}

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_mysql_database.test.name

  users {
    name     = huaweicloud_rds_mysql_account.test.name
    readonly = true
  }
}
`, testAccRdsInstance_base(), name)
}

func testDataSourceRdsMysqlAuthorizedDatabases_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_mysql_authorized_databases" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  user_name   = huaweicloud_rds_mysql_account.test.name
  depends_on  = [huaweicloud_rds_mysql_database_privilege.test]
}
`, testDataSourceRdsMysqlAuthorizedDatabases_base(name))
}

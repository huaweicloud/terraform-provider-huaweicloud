package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_databases.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussdbMysqlDatabases_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.character_set"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.users.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.users.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.users.0.host"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.users.0.readonly"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("character_set_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussdbMysqlDatabases_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_mysql_account" "test" {
  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  name        = "%[2]s"
  password    = "Test@12345678"
  host        = "10.10.10.10"
}

resource "huaweicloud_gaussdb_mysql_account_privilege" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_account.test,huaweicloud_gaussdb_mysql_database.test]

  instance_id  = huaweicloud_gaussdb_mysql_instance.test.id
  account_name = huaweicloud_gaussdb_mysql_account.test.name
  host         = huaweicloud_gaussdb_mysql_account.test.host

  databases {
    name     = huaweicloud_gaussdb_mysql_database.test.name
    readonly = false
  }
}


data "huaweicloud_gaussdb_mysql_databases" "test" {
  depends_on = [huaweicloud_gaussdb_mysql_account_privilege.test]

  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
}

locals {
  name = "%[2]s"
}
data "huaweicloud_gaussdb_mysql_databases" "name_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_account_privilege.test]

  instance_id = huaweicloud_gaussdb_mysql_instance.test.id
  name        = "%[2]s"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_databases.name_filter.databases) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_databases.name_filter.databases[*].name : v == local.name]
  )
}

locals {
  character_set = huaweicloud_gaussdb_mysql_database.test.character_set
}
data "huaweicloud_gaussdb_mysql_databases" "character_set_filter" {
  depends_on = [huaweicloud_gaussdb_mysql_account_privilege.test]

  instance_id   = huaweicloud_gaussdb_mysql_instance.test.id
  character_set = huaweicloud_gaussdb_mysql_database.test.character_set
}
output "character_set_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_databases.character_set_filter.databases) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_databases.character_set_filter.databases[*].character_set : v == local.character_set]
  )
}
`, testGaussDBDatabase_basic(name), name)
}

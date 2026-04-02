package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBDatabases_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_databases.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceTaurusDBDatabases_basic(rName),
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

func testDataSourceDataSourceTaurusDBDatabases_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_taurusdb_account" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id
  name        = "%[2]s"
  password    = "Test@12345678"
  host        = "10.10.10.10"
}

resource "huaweicloud_taurusdb_account_privilege" "test" {
  depends_on = [huaweicloud_taurusdb_account.test,huaweicloud_taurusdb_database.test]

  instance_id  = huaweicloud_taurusdb_instance.test.id
  account_name = huaweicloud_taurusdb_account.test.name
  host         = huaweicloud_taurusdb_account.test.host

  databases {
    name     = huaweicloud_taurusdb_database.test.name
    readonly = false
  }
}


data "huaweicloud_taurusdb_databases" "test" {
  depends_on = [huaweicloud_taurusdb_account_privilege.test]

  instance_id = huaweicloud_taurusdb_instance.test.id
}

locals {
  name = "%[2]s"
}
data "huaweicloud_taurusdb_databases" "name_filter" {
  depends_on = [huaweicloud_taurusdb_account_privilege.test]

  instance_id = huaweicloud_taurusdb_instance.test.id
  name        = "%[2]s"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_databases.name_filter.databases) > 0 && alltrue(
  [for v in data.huaweicloud_taurusdb_databases.name_filter.databases[*].name : v == local.name]
  )
}

locals {
  character_set = huaweicloud_taurusdb_database.test.character_set
}
data "huaweicloud_taurusdb_databases" "character_set_filter" {
  depends_on = [huaweicloud_taurusdb_account_privilege.test]

  instance_id   = huaweicloud_taurusdb_instance.test.id
  character_set = huaweicloud_taurusdb_database.test.character_set
}
output "character_set_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_databases.character_set_filter.databases) > 0 && alltrue(
  [for v in data.huaweicloud_taurusdb_databases.character_set_filter.databases[*].character_set : v == local.character_set]
  )
}
`, testGaussDBDatabase_basic(name), name)
}

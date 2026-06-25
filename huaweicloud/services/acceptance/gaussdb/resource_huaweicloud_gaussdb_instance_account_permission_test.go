package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceGaussdbInstanceAccountPermission_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_instance_account_permission.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGaussdbInstanceAccountPermission_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttr(rName, "db_name", name),
					resource.TestCheckResourceAttr(rName, "users.0.name", name),
					resource.TestCheckResourceAttr(rName, "users.0.schema_name", name),
				),
			},
			{
				Config: testAccResourceGaussdbInstanceAccountPermission_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttr(rName, "db_name", name),
					resource.TestCheckResourceAttr(rName, "users.0.name", name),
					resource.TestCheckResourceAttr(rName, "users.0.schema_name", name),
				),
			},
		},
	})
}

func testAccResourceGaussdbInstanceAccountPermission_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_database" "test" {
  instance_id   = huaweicloud_gaussdb_instance.test.id
  name          = "%[2]s"
  character_set = "UTF8"
  owner         = "root"
  template      = "template0"
  lc_collate    = "C"
  lc_ctype      = "C"
}

resource "huaweicloud_gaussdb_instance_database_account" "test" {
  instance_id   = huaweicloud_gaussdb_instance.test.id
  name          = "%[2]s"
  password      = "Test@963852"
  is_login_only = "false"
}

resource "huaweicloud_gaussdb_schema" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  db_name     = huaweicloud_gaussdb_database.test.name
  name        = "%[2]s"
  owner       = "root"
}

resource "huaweicloud_gaussdb_instance_account_permission" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  db_name     = huaweicloud_gaussdb_database.test.name

  users {
    name        = "%[2]s"
    readonly    = "true"
    schema_name = huaweicloud_gaussdb_schema.test.name
  }

depends_on = [
    huaweicloud_gaussdb_database.test,
    huaweicloud_gaussdb_instance_database_account.test,
    huaweicloud_gaussdb_schema.test
  ]
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}

func testAccResourceGaussdbInstanceAccountPermission_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_database" "test" {
  instance_id   = huaweicloud_gaussdb_instance.test.id
  name          = "%[2]s"
  character_set = "UTF8"
  owner         = "root"
  template      = "template0"
  lc_collate    = "C"
  lc_ctype      = "C"
}

resource "huaweicloud_gaussdb_instance_database_account" "test" {
  instance_id   = huaweicloud_gaussdb_instance.test.id
  name          = "%[2]s"
  password      = "Test@963852"
  is_login_only = "false"
}

resource "huaweicloud_gaussdb_schema" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  db_name     = huaweicloud_gaussdb_database.test.name
  name        = "%[2]s"
  owner       = "root"
}

resource "huaweicloud_gaussdb_instance_account_permission" "test" {
  instance_id = huaweicloud_gaussdb_instance.test.id
  db_name     = huaweicloud_gaussdb_database.test.name

  users {
    name        = "%[2]s"
    readonly    = "false"
    schema_name = huaweicloud_gaussdb_schema.test.name
  }

depends_on = [
    huaweicloud_gaussdb_database.test,
    huaweicloud_gaussdb_instance_database_account.test,
    huaweicloud_gaussdb_schema.test
  ]
}
`, testDataSourceGaussdbInstanceMetrics_base(name), name)
}

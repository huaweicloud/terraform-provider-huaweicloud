package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMysqlDatabasePrivileges_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_mysql_database_privileges.test"

	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDatabasePrivileges_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.#"),
					resource.TestCheckResourceAttrSet(rName, "users.0.name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.readonly"),
					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
					resource.TestCheckOutput("readonly_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccMysqlDatabasePrivileges_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_mysql_database_privileges" "test" {
  depends_on  = [huaweicloud_rds_mysql_database_privilege.test]
  instance_id = huaweicloud_rds_mysql_database_privilege.test.instance_id
  db_name     = huaweicloud_rds_mysql_database_privilege.test.db_name
}

locals {
  user_name = huaweicloud_rds_mysql_account.test_1.name
}

data "huaweicloud_rds_mysql_database_privileges" "user_name_filter" {
  depends_on  = [huaweicloud_rds_mysql_database_privilege.test]
  instance_id = huaweicloud_rds_mysql_database_privilege.test.instance_id
  db_name     = huaweicloud_rds_mysql_database_privilege.test.db_name
  user_name   = local.user_name
}

output "user_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_mysql_database_privileges.user_name_filter.users) > 0 && alltrue(
  [for v in data.huaweicloud_rds_mysql_database_privileges.user_name_filter.users[*].name : v == local.user_name]
  )
}

locals {
  readonly = tolist(huaweicloud_rds_mysql_database_privilege.test.users)[0].readonly
}

data "huaweicloud_rds_mysql_database_privileges" "readonly_filter" {
  depends_on  = [huaweicloud_rds_mysql_database_privilege.test]
  instance_id = huaweicloud_rds_mysql_database_privilege.test.instance_id
  db_name     = huaweicloud_rds_mysql_database_privilege.test.db_name
  readonly    = local.readonly
}

output "readonly_filter_is_useful" {
  value = length(data.huaweicloud_rds_mysql_database_privileges.readonly_filter.users) > 0 && alltrue(
  [for v in data.huaweicloud_rds_mysql_database_privileges.readonly_filter.users[*].readonly : v == local.readonly]
)
}

`, testAccRdsDatabasePrivilege_basic(name))
}

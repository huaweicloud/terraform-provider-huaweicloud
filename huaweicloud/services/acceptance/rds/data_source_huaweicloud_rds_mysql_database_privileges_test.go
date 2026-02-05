package rds

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMysqlDatabasePrivileges_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_rds_mysql_database_privileges.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byUserName   = "data.huaweicloud_rds_mysql_database_privileges.filter_by_user_name"
		dcByUserName = acceptance.InitDataSourceCheck(byUserName)

		byReadonly   = "data.huaweicloud_rds_mysql_database_privileges.filter_by_readonly"
		dcByReadonly = acceptance.InitDataSourceCheck(byReadonly)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataMysqlDatabasePrivileges_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByUserName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byUserName, "users.0.name"),
					resource.TestCheckResourceAttrSet(byUserName, "users.0.readonly"),
					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
					dcByReadonly.CheckResourceExists(),
					resource.TestCheckOutput("readonly_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataMysqlDatabasePrivileges_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  # The behavior of parameter 'name' of the database resource is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [huaweicloud_rds_mysql_database.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_mysql_database.test.name
  
  dynamic "users" {
    for_each = slice(huaweicloud_rds_mysql_account.test[*].name, 0, 1)

    content {
      name     = users.value
      readonly = false
    }
  }
  dynamic "users" {
    for_each = slice(huaweicloud_rds_mysql_account.test[*].name, 1, 2)

    content {
      name     = users.value
      readonly = true
    }
  }
}

data "huaweicloud_rds_mysql_database_privileges" "test" {
  # The behavior of parameter 'instance_id' and 'db_name' of the privilege resource is 'Required', means this
  # parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_rds_mysql_database_privilege.test]

  instance_id = huaweicloud_rds_mysql_database_privilege.test.instance_id
  db_name     = huaweicloud_rds_mysql_database_privilege.test.db_name
}

locals {
  user_name = huaweicloud_rds_mysql_account.test[0].name # The first account is the one with the read-and-write permission.
}

data "huaweicloud_rds_mysql_database_privileges" "filter_by_user_name" {
  # The behavior of parameter 'instance_id' and 'db_name' of the privilege resource is 'Required', means this
  # parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_rds_mysql_database_privilege.test]

  instance_id = huaweicloud_rds_mysql_database_privilege.test.instance_id
  db_name     = huaweicloud_rds_mysql_database_privilege.test.db_name
  user_name   = local.user_name
}

output "user_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_mysql_database_privileges.filter_by_user_name.users) > 0 && alltrue(
    [for v in data.huaweicloud_rds_mysql_database_privileges.filter_by_user_name.users[*].name : v == local.user_name]
  )
}

locals {
  readonly = tolist(huaweicloud_rds_mysql_database_privilege.test.users)[0].readonly
}

data "huaweicloud_rds_mysql_database_privileges" "filter_by_readonly" {
  # The behavior of parameter 'instance_id' and 'db_name' of the privilege resource is 'Required', means this
  # parameter does not have 'Know After Apply' behavior.
  depends_on = [huaweicloud_rds_mysql_database_privilege.test]

  instance_id = huaweicloud_rds_mysql_database_privilege.test.instance_id
  db_name     = huaweicloud_rds_mysql_database_privilege.test.db_name
  readonly    = local.readonly
}

output "readonly_filter_is_useful" {
  value = length(data.huaweicloud_rds_mysql_database_privileges.filter_by_readonly.users) > 0 && alltrue(
    [for v in data.huaweicloud_rds_mysql_database_privileges.filter_by_readonly.users[*].readonly : v == local.readonly]
  )
}
`, testAccMysqlDatabasePrivilege_basic_base(name))
}

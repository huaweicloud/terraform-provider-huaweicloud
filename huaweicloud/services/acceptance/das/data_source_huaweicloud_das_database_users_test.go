package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDatabaseUsers_basic(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()

		all = "data.huaweicloud_das_database_users.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByUserId   = "data.huaweicloud_das_database_users.filter_by_user_id"
		dcFilterByUserId = acceptance.InitDataSourceCheck(filterByUserId)

		filterByUserName   = "data.huaweicloud_das_database_users.filter_by_user_name"
		dcFilterByUserName = acceptance.InitDataSourceCheck(filterByUserName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDatabaseUsers_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// filter by user ID
					dcFilterByUserId.CheckResourceExists(),
					resource.TestCheckOutput("is_user_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(filterByUserId, "users.0.id"),
					resource.TestCheckResourceAttrSet(filterByUserId, "users.0.name"),

					// filter by user name
					dcFilterByUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDatabaseUsers_basic_base(name, password string) string {
	return fmt.Sprintf(` 
resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = "%[1]s"
  name          = "%[2]s"
  character_set = "utf8"
}

resource "huaweicloud_rds_mysql_account" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database.test,
  ]

  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database.test,
    huaweicloud_rds_mysql_account.test,
  ]

  instance_id = "%[1]s"
  db_name     = huaweicloud_rds_mysql_database.test.name

  users {
    name     = huaweicloud_rds_mysql_account.test.name
    readonly = false
  }
}

resource "huaweicloud_das_database_user" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database_privilege.test,
  ]

  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password)
}

func testAccDataDatabaseUsers_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_database_users" "all" {
  depends_on = [
    huaweicloud_das_database_user.test,
  ]

  instance_id = "%[2]s"
}

# Filter by user ID
locals {
  user_id = huaweicloud_das_database_user.test.id
}

data "huaweicloud_das_database_users" "filter_by_user_id" {
  depends_on = [
    huaweicloud_das_database_user.test,
  ]

  instance_id = "%[2]s"
  user_id     = local.user_id
}

locals {
  user_id_filter_result = [
    for v in data.huaweicloud_das_database_users.filter_by_user_id.users : v.id == local.user_id
  ]
}

output "is_user_id_filter_useful" {
  value = length(local.user_id_filter_result) > 0 && alltrue(local.user_id_filter_result)
}

# Filter by user name
locals {
  user_name = huaweicloud_das_database_user.test.name
}

data "huaweicloud_das_database_users" "filter_by_user_name" {
  depends_on = [
    huaweicloud_das_database_user.test,
  ]

  instance_id = "%[2]s"
  user_name   = local.user_name
}

locals {
  user_name_filter_result = [
    for v in data.huaweicloud_das_database_users.filter_by_user_name.users : v.name == local.user_name
  ]
}

output "is_user_name_filter_useful" {
  value = length(local.user_name_filter_result) > 0 && alltrue(local.user_name_filter_result)
}
`, testAccDataDatabaseUsers_basic_base(name, password), acceptance.HW_RDS_INSTANCE_ID)
}

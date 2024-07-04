package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceSQLServerAccounts_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_sqlserver_accounts.test"
	dbPwd := fmt.Sprintf("%s%s%d", acctest.RandString(5),
		acctest.RandStringFromCharSet(2, "!#%^*"), acctest.RandIntRange(10, 99))
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSQLServerAccounts_basic(name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.#"),
					resource.TestCheckResourceAttrSet(rName, "users.0.name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.state"),
					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
					resource.TestCheckOutput("state_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccSQLServerAccounts_basic(name string, dbPwd string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_sqlserver_accounts" "test" {
  depends_on  = [huaweicloud_rds_sqlserver_account.test]
  instance_id = huaweicloud_rds_sqlserver_account.test.instance_id
}

data "huaweicloud_rds_sqlserver_accounts" "user_name_filter" {
  depends_on  = [huaweicloud_rds_sqlserver_account.test]
  instance_id = huaweicloud_rds_sqlserver_account.test.instance_id
  user_name   = huaweicloud_rds_sqlserver_account.test.name
}

locals {
  user_name = huaweicloud_rds_sqlserver_account.test.name
}

output "user_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_sqlserver_accounts.user_name_filter.users) > 0 && alltrue(
    [for v in data.huaweicloud_rds_sqlserver_accounts.user_name_filter.users[*].name : v == local.user_name]
  )
}

data "huaweicloud_rds_sqlserver_accounts" "state_filter" {
  depends_on  = [huaweicloud_rds_sqlserver_account.test]
  instance_id = huaweicloud_rds_sqlserver_account.test.instance_id
  state       = huaweicloud_rds_sqlserver_account.test.state
}

locals {
  state = huaweicloud_rds_sqlserver_account.test.state
}

output "state_filter_is_useful" {
  value = length(data.huaweicloud_rds_sqlserver_accounts.state_filter.users) > 0 && alltrue(
    [for v in data.huaweicloud_rds_sqlserver_accounts.state_filter.users[*].state : v == local.state]
  )
}
`, testSQLServerAccount_basic(name, dbPwd))
}

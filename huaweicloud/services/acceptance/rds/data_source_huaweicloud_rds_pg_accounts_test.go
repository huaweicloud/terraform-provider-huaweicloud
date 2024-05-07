package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePgAccounts_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_pg_accounts.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePgAccounts_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.#"),
					resource.TestCheckResourceAttrSet(rName, "users.0.name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolsuper"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolinherit"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolcreaterole"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolcreatedb"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolcanlogin"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolconnlimit"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolreplication"),
					resource.TestCheckResourceAttrSet(rName, "users.0.attributes.0.rolbypassrls"),
					resource.TestCheckResourceAttrSet(rName, "users.0.description"),
					resource.TestCheckOutput("user_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourcePgAccounts_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_pg_accounts" "test" {
  depends_on  = [huaweicloud_rds_pg_account.test]
  instance_id = huaweicloud_rds_pg_account.test.instance_id
  user_name   = "%s"
}

data "huaweicloud_rds_pg_accounts" "user_name_filter" {
  depends_on  = [huaweicloud_rds_pg_account.test]
  instance_id = huaweicloud_rds_pg_account.test.instance_id
  user_name   = huaweicloud_rds_pg_account.test.name
}

locals {
  user_name = huaweicloud_rds_pg_account.test.name
}

output "user_name_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_accounts.user_name_filter.users) > 0 && alltrue(
    [for v in data.huaweicloud_rds_pg_accounts.user_name_filter.users[*].name : v == local.user_name]
  )
}

`, testPgAccount_basic(name), name)
}

package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMysqlAccounts_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_rds_mysql_accounts.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlAccounts_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "users.#"),
					resource.TestCheckResourceAttrSet(rName, "users.0.name"),
					resource.TestCheckResourceAttrSet(rName, "users.0.hosts.0"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("host_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccMysqlAccounts_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_mysql_accounts" "test" {
  depends_on  = [huaweicloud_rds_mysql_account.test]
  instance_id = huaweicloud_rds_instance.test.id
}

data "huaweicloud_rds_mysql_accounts" "name_filter" {
  depends_on  = [huaweicloud_rds_mysql_account.test]
  instance_id = huaweicloud_rds_instance.test.id
  name        = huaweicloud_rds_mysql_account.test.name
}

locals {
  name = huaweicloud_rds_mysql_account.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_rds_mysql_accounts.name_filter.users) > 0 && alltrue(
    [for v in data.huaweicloud_rds_mysql_accounts.name_filter.users[*].name : v == local.name]
  )
}

data "huaweicloud_rds_mysql_accounts" "host_filter" {
  depends_on  = [huaweicloud_rds_mysql_account.test]
  instance_id = huaweicloud_rds_instance.test.id
  host        = huaweicloud_rds_mysql_account.test.hosts.0
}

locals {
  host = huaweicloud_rds_mysql_account.test.hosts.0
}

output "host_filter_is_useful" {
  value = length(data.huaweicloud_rds_mysql_accounts.host_filter.users) > 0 && alltrue(
    [for v in data.huaweicloud_rds_mysql_accounts.host_filter.users[*].hosts : contains(v, local.host)]
  )
}
`, testMysqlAccount_basic(name))
}

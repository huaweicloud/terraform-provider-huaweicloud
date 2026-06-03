package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSqlExecutionPlans_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_das_sql_execution_plans.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlExecutionPlans_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "plans.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSource, "plans.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "plans.0.select_type"),
					resource.TestCheckResourceAttrSet(dataSource, "plans.0.extra"),
				),
			},
		},
	})
}

func testAccSqlExecutionPlans_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}

data "huaweicloud_das_database_users" "all" {
  instance_id = local.instance_ids[0]
}

locals {
  db_user_id = try(data.huaweicloud_das_database_users.all.users[0].id, "")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccSqlExecutionPlans_basic() string {
	return fmt.Sprintf(`
%[1]s

# '__recyclebin__' is a default database
# 'SELECT 1' means to query all the plans
data "huaweicloud_das_sql_execution_plans" "test" {
  instance_id = local.instance_ids[0]
  db_user_id  = local.db_user_id
  database    = "__recyclebin__"
  sql         = "SELECT 1"
}
`, testAccSqlExecutionPlans_base())
}

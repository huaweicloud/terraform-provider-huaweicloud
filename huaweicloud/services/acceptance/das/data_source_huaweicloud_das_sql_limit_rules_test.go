package das

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSqlLimitRules_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_sql_limit_rules.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlLimitRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "rules.0.id"),
					resource.TestCheckResourceAttrSet(all, "rules.0.sql_type"),
					resource.TestCheckResourceAttrSet(all, "rules.0.max_concurrency"),
					resource.TestCheckResourceAttrSet(all, "rules.0.pattern"),
				),
			},
		},
	})
}

func testAccSqlLimitRules_basic_base() string {
	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
}
`, acceptance.HW_DAS_INSTANCE_IDS)
}

func testAccSqlLimitRules_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_das_sql_limit_rules" "all" {
  instance_id = local.instance_ids[0]
}
`, testAccSqlLimitRules_basic_base())
}

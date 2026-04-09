package dli

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSystemSQLDefendRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dli_system_sql_defend_rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSystemSQLDefendRules_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.no_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.engines.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.actions.#"),
					resource.TestMatchResourceAttr(dataSource, "rules.0.param.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.param.0.default_value"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.param.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.param.0.max"),
				),
			},
		},
	})
}

const testAccDataSystemSQLDefendRules_basic = `data "huaweicloud_dli_system_sql_defend_rules" "test" {}`

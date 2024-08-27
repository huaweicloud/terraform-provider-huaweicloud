package sfsturbo

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePermRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_sfs_turbo_perm_rules.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePermRules_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.ip_cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.rw_type"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.user_type"),
				),
			},
		},
	})
}

func testDataSourcePermRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_sfs_turbo_perm_rules" "test" {
  depends_on = [
    huaweicloud_sfs_turbo_perm_rule.test
  ]

  share_id = huaweicloud_sfs_turbo.test.id
}
`, testSFSTruboPermRuleBasic(name))
}

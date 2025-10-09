package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineCheckRuleHAB_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_baseline_check_rule_hab.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBaselineCheckRuleHAB_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_rule_num"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_num"),
					resource.TestCheckResourceAttrSet(dataSource, "host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceBaselineCheckRuleHAB_basic = `
data "huaweicloud_hss_baseline_check_rule_hab" "test" {
  action        = "ignore"
  handle_status = "unhandled"

  check_rule_list {
    check_name = "example_check"
  }
}
`

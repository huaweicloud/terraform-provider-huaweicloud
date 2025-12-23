package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineCheckRuleFailDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_baseline_check_rule_fail_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a host with host enterprise edition protection enabled.
			// Configure corresponding policy and manual check.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBaselineCheckRuleFailDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "fail_detail_list.#"),
				),
			},
		},
	})
}

const testDataSourceBaselineCheckRuleFailDetail_base = `
data "huaweicloud_hss_baseline_risk_config_check_rules" "test" {
  check_name = "SSH"
  standard   = "hw_standard"
}
`

func testDataSourceBaselineCheckRuleFailDetail_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_baseline_check_rule_fail_detail" "test" {
  check_rule_id         = data.huaweicloud_hss_baseline_risk_config_check_rules.test.data_list[0].check_rule_id
  enterprise_project_id = "0"
  check_name            = "SSH"
  standard              = "hw_standard"
}
`, testDataSourceBaselineCheckRuleFailDetail_base)
}

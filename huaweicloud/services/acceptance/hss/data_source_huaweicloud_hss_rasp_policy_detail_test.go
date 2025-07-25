package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRaspPolicyDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_rasp_policy_detail.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, prepare an HSS protection policy.
			acceptance.TestAccPreCheckHSSPolicyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceRaspPolicyDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.chk_feature_id"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.chk_feature_name"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.chk_feature_desc"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.feature_configure"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.protective_action"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.optional_protective_action"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "rule_list.0.editable"),

					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRaspPolicyDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_rasp_policy_detail" "test" {
  policy_id = "%[1]s"
}

data "huaweicloud_hss_rasp_policy_detail" "eps_filter" {
  policy_id             = "%[1]s"
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_policy_detail.eps_filter.rule_list) > 0
}
`, acceptance.HW_HSS_POLICY_ID)
}

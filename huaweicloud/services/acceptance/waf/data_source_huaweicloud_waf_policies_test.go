package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccDataSourceWafPoliciesV1_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_policies.policies_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPoliciesV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPoliciesID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "policies.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "policies.0.name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.blacklist"),
				),
			},
		},
	})
}

func testAccCheckWafPoliciesID(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmtp.Errorf("Can't find WAF policies data source: %s.", r)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("The WAF policies data source ID does not set.")
		}
		return nil
	}
}

func testAccWafPoliciesV1_conf(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_policies" "policies_1" {
  name = huaweicloud_waf_policy.policy_1.name

  depends_on = [
    huaweicloud_waf_policy.policy_1
  ]
}
`, testAccWafPolicyV1_basic(name))
}

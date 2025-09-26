package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPolicyCopy_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyCopy_basic(randName),
			},
		},
	})
}

func testAccPolicyCopy_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                  = "%[1]s"
  level                 = 1
  enterprise_project_id = "0"
}
`, name)
}

func testAccPolicyCopy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy_copy" "test" {
  src_policy_id         = huaweicloud_waf_policy.test.id
  dest_policy_name      = "%[2]s_copy"
  enterprise_project_id = "0"
}
`, testAccPolicyCopy_base(name), name)
}

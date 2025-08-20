package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePolicyBlackWhiteLists_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_aad_policy_black_white_lists.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an AAD domain name that has been configured with the policy blackWhite list and configure
			// it in the environment variable.
			acceptance.TestAccPrecheckAadDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePolicyBlackWhiteLists_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "black.#"),
					resource.TestCheckResourceAttrSet(dataSource, "white.#"),
				),
			},
		},
	})
}

func testDataSourcePolicyBlackWhiteLists_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_policy_black_white_lists" "test" {
  domain_name   = "%[1]s"
  overseas_type = 0
}
`, acceptance.HW_AAD_DOMAIN_NAME)
}

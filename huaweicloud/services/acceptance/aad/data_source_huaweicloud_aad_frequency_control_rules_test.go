package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `items`.
func TestAccDataSourceFrequencyControlRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_aad_frequency_control_rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an AAD domain name that has been configured with frequency control rules and configure it
			// in the environment variable.
			acceptance.TestAccPrecheckAadDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFrequencyControlRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
				),
			},
		},
	})
}

func testDataSourceFrequencyControlRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_frequency_control_rules" "test" {
  domain_name = "%[1]s"
}
`, acceptance.HW_AAD_DOMAIN_NAME)
}

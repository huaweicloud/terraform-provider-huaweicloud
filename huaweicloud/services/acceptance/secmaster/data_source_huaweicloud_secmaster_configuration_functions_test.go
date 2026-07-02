package secmaster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceConfigurationFunctions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_configuration_functions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceConfigurationFunctions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "support_postpaid_mgmt"),
					resource.TestCheckResourceAttrSet(dataSource, "support_large_screen_mgmt"),
					resource.TestCheckResourceAttrSet(dataSource, "support_purchase_label_mgmt"),
					resource.TestCheckResourceAttrSet(dataSource, "billing_type_mgmt"),
				),
			},
		},
	})
}

func testDataSourceConfigurationFunctions_basic() string {
	return `
data "huaweicloud_secmaster_configuration_functions" "test" {}
`
}

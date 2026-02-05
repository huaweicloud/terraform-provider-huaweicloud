package ccm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcmPrivateCaConfigConsole_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ccm_private_ca_config_console.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcmPrivateCaConfigConsole_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "is_support_eps"),
					resource.TestCheckResourceAttrSet(dataSource, "is_support_sm2"),
					resource.TestCheckResourceAttrSet(dataSource, "is_support_dhsm"),
					resource.TestCheckResourceAttrSet(dataSource, "is_support_yearly_monthly_ca"),
					resource.TestCheckResourceAttrSet(dataSource, "is_support_iam5"),
					resource.TestCheckResourceAttrSet(dataSource, "is_support_ocsp"),
				),
			},
		},
	})
}

const testDataSourceDataSourceCcmPrivateCaConfigConsole_basic = `
data "huaweicloud_ccm_private_ca_config_console" "test" {}
`

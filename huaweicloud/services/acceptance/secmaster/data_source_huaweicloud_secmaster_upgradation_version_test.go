package secmaster

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceUpgradationVersion_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_upgradation_version.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceUpgradationVersion_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "version"),
				),
			},
		},
	})
}

const testDataSourceUpgradationVersion_basic = `data "huaweicloud_secmaster_upgradation_version" "test" {}`

package ccm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcmPrivateCaQuota_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ccm_private_ca_quota.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcmPrivateCaQuota_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.quota"),
				),
			},
		},
	})
}

const testDataSourceDataSourceCcmPrivateCaQuota_basic = `
data "huaweicloud_ccm_private_ca_quota" "test" {}
`

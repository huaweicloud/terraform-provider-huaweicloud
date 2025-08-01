package sdrs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sdrs_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.quota"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.min"),
				),
			},
		},
	})
}

const testDataSourceQuotas_basic = `
data "huaweicloud_sdrs_quotas" "test" {}
`

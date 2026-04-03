package ga

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ga_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota"),
				),
			},
		},
	})
}

const testDataSourceGaQuotas_basic = `
data "huaweicloud_ga_quotas" "test" {}
`

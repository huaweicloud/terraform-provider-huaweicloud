package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcsQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.min"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.max"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.resources.0.quota"),
				),
			},
		},
	})
}

func testDataSourceDcsQuotas_basic() string {
	return `
data "huaweicloud_dcs_quotas" "test" {}
`
}

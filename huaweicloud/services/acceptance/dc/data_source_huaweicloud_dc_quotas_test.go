package dc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dc_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byType   = "data.huaweicloud_dc_quotas.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcQuotas_basic() string {
	return (`
data "huaweicloud_dc_quotas" "test" {}

locals {
  type = data.huaweicloud_dc_quotas.test.quotas[0].type
}

data "huaweicloud_dc_quotas" "filter_by_type" {
  type = [local.type]
}

output "type_filter_useful" {
  value = length(data.huaweicloud_dc_quotas.filter_by_type.quotas) > 0 && alltrue(
    [for v in data.huaweicloud_dc_quotas.filter_by_type.quotas[*].type : v == local.type]
  )
}
`)
}

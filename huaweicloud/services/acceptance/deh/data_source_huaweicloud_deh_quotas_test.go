package deh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDehQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_deh_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDehQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.resource"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.hard_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quota_set.0.used"),

					resource.TestCheckOutput("resource_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDehQuotas_basic() string {
	return `
data "huaweicloud_deh_quotas" "test" {}

locals {
  resource = data.huaweicloud_deh_quotas.test.quota_set[0].resource
}

data "huaweicloud_deh_quotas" "resource_filter" {
  resource = local.resource
}

output "resource_filter_is_useful" {
  value = length(data.huaweicloud_deh_quotas.resource_filter.quota_set) > 0 && alltrue(
    [for v in data.huaweicloud_deh_quotas.resource_filter.quota_set[*].resource :
      v == local.resource]
  )
}
`
}

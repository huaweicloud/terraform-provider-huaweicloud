package eg

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEgQuotas_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_eg_quotas.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_eg_quotas.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEgQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(all, "quotas.#"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.name"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.quota"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.max"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.min"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceEgQuotas_basic = `
data "huaweicloud_eg_quotas" "test" {}

# Filter by type
locals {
  quota_type = try(data.huaweicloud_eg_quotas.test.quotas[0].type, "")
}

data "huaweicloud_eg_quotas" "filter_by_type" {
  depends_on = [
    data.huaweicloud_eg_quotas.test
  ]

  type = local.quota_type
}

locals {
  type_filter_result = [
    for v in try(data.huaweicloud_eg_quotas.filter_by_type.quotas[*].type) : v == local.quota_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}
`

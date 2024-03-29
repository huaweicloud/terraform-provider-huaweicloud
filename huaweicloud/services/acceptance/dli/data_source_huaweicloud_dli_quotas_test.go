package dli

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliQuotas_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dli_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckOutput("no_filter_useful", "true"),
					resource.TestCheckOutput("cu_filter", "true"),
					resource.TestCheckOutput("absent_type_test", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDliQuotas_basic() string {
	return `
data "huaweicloud_dli_quotas" "test" {}

data "huaweicloud_dli_quotas" "cu_filter" {
  type = "CU"
}

data "huaweicloud_dli_quotas" "absent_type_test" {
  type = "not found"
}

output "no_filter_useful" {
  value = length(data.huaweicloud_dli_quotas.test.quotas) > 1
}

output "cu_filter" {
  value = length(data.huaweicloud_dli_quotas.cu_filter.quotas) == 1
}

output "absent_type_test" {
  value = length(data.huaweicloud_dli_quotas.absent_type_test.quotas) == 0
}
`
}

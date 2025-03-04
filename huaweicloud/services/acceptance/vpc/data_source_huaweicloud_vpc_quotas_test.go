package vpc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcQuotas_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_vpc_quotas.basic"
	dataSource2 := "data.huaweicloud_vpc_quotas.filter_by_type"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceVpcQuotas_basic = `
data "huaweicloud_vpc_quotas" "basic" {}

data "huaweicloud_vpc_quotas" "filter_by_type" {
  type = "vpc"
}

locals {
  type_filter_result = [for v in data.huaweicloud_vpc_quotas.filter_by_type.quotas[0].resources[*].type : v == "vpc"]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_quotas.basic.quotas[0].resources) > 0
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`

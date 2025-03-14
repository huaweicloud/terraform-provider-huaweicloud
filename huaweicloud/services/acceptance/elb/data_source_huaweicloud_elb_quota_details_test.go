package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceElbQuotaDetails_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_quota_details.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbQuotaDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_key"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.quota_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.used"),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.unit"),
					resource.TestCheckOutput("quota_key_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceElbQuotaDetails_basic() string {
	return `
data "huaweicloud_elb_quota_details" "test" {}

data "huaweicloud_elb_quotas" "test" {}

data "huaweicloud_elb_quota_details" "quota_key_filter" {
  quota_key = ["loadbalancer"]
}

locals {
  quota_limit = data.huaweicloud_elb_quotas.test.loadbalancer
}

output "quota_key_filter_is_useful" {
  value = length(data.huaweicloud_elb_quota_details.quota_key_filter.quotas) > 0 && alltrue(
  [for v in data.huaweicloud_elb_quota_details.quota_key_filter.quotas[*].quota_limit : v == local.quota_limit]
  )  
}
`
}

package ccm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCertificatesByTags_basic(t *testing.T) {
	var (
		datasource = "data.huaweicloud_ccm_certificates_by_tags.test"
		dc         = acceptance.InitDataSourceCheck(datasource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Prepare certificates in the runtime environment before running test cases.
			acceptance.TestAccPreCheckCCMEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCertificatesByTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.cert_id"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.cert_name"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.domain"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.cert_type"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.cert_brand"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.domain_type"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.purchase_period"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.expired_time"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.order_status"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.domain_num"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.wildcard_number"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail.0.sans"),

					resource.TestCheckOutput("non_exist_resources_is_zero", "true"),
					resource.TestCheckOutput("total_count_greater_than_zero", "true"),
				),
			},
		},
	})
}

const testAccDatasourceCertificatesByTags_basic = `
data "huaweicloud_ccm_certificates_by_tags" "test" {
  resource_instances = "resource_instances"
  action             = "filter"
}

data "huaweicloud_ccm_certificates_by_tags" "non-exist" {
  resource_instances = "resource_instances"
  action             = "filter"
  tags {
    key    = "non-exist-key"
    values = ["non-exist-value"]
  }
  matches {
    key   = "resource_name"
    value = "non-exist-value"
  }
}

output "non_exist_resources_is_zero" {
  value = length(data.huaweicloud_ccm_certificates_by_tags.non-exist.resources) == 0
}

data "huaweicloud_ccm_certificates_by_tags" "total-count" {
  resource_instances = "resource_instances"
  action             = "count"
}

output "total_count_greater_than_zero" {
  value = data.huaweicloud_ccm_certificates_by_tags.total-count.total_count > 0
}
`

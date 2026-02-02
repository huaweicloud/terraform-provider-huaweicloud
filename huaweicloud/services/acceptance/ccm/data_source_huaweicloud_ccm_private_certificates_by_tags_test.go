package ccm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePrivateCertificatesByTags_basic(t *testing.T) {
	var (
		datasource = "data.huaweicloud_ccm_private_certificates_by_tags.test"
		dc         = acceptance.InitDataSourceCheck(datasource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Prepare private certificates in the runtime environment before running test cases.
			acceptance.TestAccPreCheckCCMEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePrivateCertificatesByTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_detail"),

					resource.TestCheckOutput("non_exist_resources_is_zero", "true"),
				),
			},
		},
	})
}

const testAccDatasourcePrivateCertificatesByTags_basic = `
data "huaweicloud_ccm_private_certificates_by_tags" "test" {}

data "huaweicloud_ccm_private_certificates_by_tags" "non-exist" {
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
  value = length(data.huaweicloud_ccm_private_certificates_by_tags.non-exist.resources) == 0
}
`

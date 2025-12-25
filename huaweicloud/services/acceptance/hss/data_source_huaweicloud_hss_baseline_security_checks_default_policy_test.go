package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineSecurityChecksDefaultPolicy_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_baseline_security_checks_default_policy.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBaselineSecurityChecksDefaultPolicy_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "content"),
				),
			},
		},
	})
}

func testDataSourceBaselineSecurityChecksDefaultPolicy_basic() string {
	return `
data "huaweicloud_hss_baseline_security_checks_default_policy" "test" {
  support_os            = "Linux"
  enterprise_project_id = "0"
}
`
}

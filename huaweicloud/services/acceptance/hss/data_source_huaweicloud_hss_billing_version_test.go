package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBillingVersion_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_billing_version.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case needs to ensure the existence of a host with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBillingVersion_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "version"),
				),
			},
		},
	})
}

func testDataSourceBillingVersion_basic() string {
	return `
data "huaweicloud_hss_billing_version" "test" {
  enterprise_project_id = "0"
}
`
}

package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceUserBundle_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_bundle.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUserBundle_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "host"),
					resource.TestCheckResourceAttrSet(dataSourceName, "premium_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "premium_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "premium_host"),
				),
			},
		},
	})
}

const testAccDataSourceUserBundle_basic = `
data "huaweicloud_waf_bundle" "test" {
  enterprise_project_id = "0"
}`

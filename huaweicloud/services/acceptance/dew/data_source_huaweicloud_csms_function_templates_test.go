package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCsmsFunctionTemplates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_csms_function_templates.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCsmsFunctionTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "function_templates"),
				),
			},
		},
	})
}

const testDataSourceCsmsFunctionTemplates_basic = `
data "huaweicloud_csms_function_templates" "test" {
  secret_type     = "RDS-FG"
  secret_sub_type = "SingleUser"
  engine          = "mysql"
}`

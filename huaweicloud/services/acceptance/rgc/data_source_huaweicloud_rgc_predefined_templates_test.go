package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePredefinedTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_predefined_templates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePredefinedTemplatesConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.template_description"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.template_category"),
				),
			},
		},
	})
}

const testAccDataSourcePredefinedTemplatesConfig_basic = `
data "huaweicloud_rgc_predefined_templates" "test" {}
`

package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceMessageTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_message_templates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceMessageTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "message_templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "message_templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "message_templates.0.locale"),
					resource.TestCheckResourceAttrSet(dataSource, "message_templates.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "message_templates.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "message_templates.0.created_at"),
				),
			},
		},
	})
}

func testDataSourceMessageTemplates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_aom_message_templates" "test" {
  depends_on = [huaweicloud_aom_message_template.test]

  enterprise_project_id = huaweicloud_aom_message_template.test.enterprise_project_id
}
`, testMessageTemplate_basic(name))
}

package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceParameterTemplateNameValidation_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dds_parameter_template_name_validation.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceParameterTemplateNameValidation_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "is_existed"),
				),
			},
		},
	})
}

func testDataSourceParameterTemplateNameValidation_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dds_parameter_template_name_validation" "test" {
  name = "%s"
}
`, name)
}

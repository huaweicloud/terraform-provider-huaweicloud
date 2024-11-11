package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSParameterTemplateCopy_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testParameterTemplateCopy_basic(name),
			},
		},
	})
}

func testParameterTemplateCopy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_parameter_template_copy" "test" {
  configuration_id = huaweicloud_dds_parameter_template.test.id
  name             = "%[2]s_copy"
  description      = "test_copy"
}
`, testDdsParameterTemplate_basic(name), name)
}

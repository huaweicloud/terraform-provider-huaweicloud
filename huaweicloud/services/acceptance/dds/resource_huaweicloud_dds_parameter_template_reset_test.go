package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSParameterTemplateReset_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testParameterTemplateReset_basic(name),
			},
		},
	})
}

func testParameterTemplateReset_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dds_parameter_template_reset" "test" {
  configuration_id = huaweicloud_dds_parameter_template.test.id
}
`, testDdsParameterTemplate_basic(name))
}

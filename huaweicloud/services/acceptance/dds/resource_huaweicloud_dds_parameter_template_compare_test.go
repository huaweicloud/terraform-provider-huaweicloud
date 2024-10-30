package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDDSParameterTemplateCompare_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dds_parameter_template_compare.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testParameterTemplateCompare_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "differences.#", "1"),
					resource.TestCheckResourceAttr(rName, "differences.0.parameter_name", "connPoolMaxConnsPerHost"),
					resource.TestCheckResourceAttr(rName, "differences.0.source_value", "800"),
					resource.TestCheckResourceAttr(rName, "differences.0.target_value", "500"),
				),
			},
		},
	})
}

func testParameterTemplateCompare_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_parameter_template" "source" {
  name         = "%[1]s_source"
  node_type    = "mongos"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost = 800
  }
}

resource "huaweicloud_dds_parameter_template" "target" {
  name         = "%[1]s_target"
  node_type    = "mongos"
  node_version = "4.0"

  parameter_values = {
    connPoolMaxConnsPerHost = 500
  }
}

resource "huaweicloud_dds_parameter_template_compare" "test" {
  source_configuration_id = huaweicloud_dds_parameter_template.source.id
  target_configuration_id = huaweicloud_dds_parameter_template.target.id
}
`, name)
}

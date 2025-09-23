package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccOpenGaussParameterTemplateCompare_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_parameter_template_compare.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussParameterrTemplateCompare_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "differences.#", "1"),
					resource.TestCheckResourceAttr(rName, "differences.0.name", "cn:auto_explain_log_min_duration"),
					resource.TestCheckResourceAttr(rName, "differences.0.source_value", "1000"),
					resource.TestCheckResourceAttr(rName, "differences.0.target_value", "500"),
				),
			},
		},
	})
}

func testOpenGaussParameterrTemplateCompare_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_parameter_template" "source" {
  name           = "%[1]s_source"
  engine_version = "8.201"
  instance_mode  = "independent"

  parameters {
    name  = "cn:auto_explain_log_min_duration"
    value = "1000"
  }
}

resource "huaweicloud_gaussdb_opengauss_parameter_template" "target" {
  name           = "%[1]s_target"
  engine_version = "8.201"
  instance_mode  = "independent"

  parameters {
    name  = "cn:auto_explain_log_min_duration"
    value = "500"
  }
}

resource "huaweicloud_gaussdb_opengauss_parameter_template_compare" "test" {
  source_id = huaweicloud_gaussdb_opengauss_parameter_template.source.id
  target_id = huaweicloud_gaussdb_opengauss_parameter_template.target.id
}
`, name)
}

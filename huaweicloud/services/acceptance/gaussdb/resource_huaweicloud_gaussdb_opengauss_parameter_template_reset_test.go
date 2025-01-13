package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccOpenGaussParameterTemplateReset_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussParameterrTemplateReset_basic(name),
			},
		},
	})
}

func testOpenGaussParameterrTemplateReset_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_opengauss_parameter_template" "test" {
  name           = "%[1]s_source"
  engine_version = "8.201"
  instance_mode  = "independent"

  parameters {
    name  = "cn:auto_explain_log_min_duration"
    value = "1000"
  }

  lifecycle {
    ignore_changes = [parameters]
  }
}

resource "huaweicloud_gaussdb_opengauss_parameter_template_reset" "test" {
  config_id = huaweicloud_gaussdb_opengauss_parameter_template.test.id
}
`, name)
}

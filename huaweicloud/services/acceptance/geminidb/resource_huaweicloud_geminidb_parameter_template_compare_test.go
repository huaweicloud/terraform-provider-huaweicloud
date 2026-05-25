package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBParameterTemplateCompare_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_geminidb_parameter_template_compare.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBParameterTemplateCompare_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "source_configuration_id"),
					resource.TestCheckResourceAttrSet(resourceName, "target_configuration_id"),

					resource.TestCheckResourceAttr(resourceName, "differences.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "differences.0.parameter_name", "request_timeout_in_ms"),
					resource.TestCheckResourceAttr(resourceName, "differences.0.source_value", "20000"),
					resource.TestCheckResourceAttr(resourceName, "differences.0.target_value", "30000"),
				),
			},
		},
	})
}

func testAccGeminiDBParameterTemplateCompare_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_parameter_template" "source" {
  name = "%[1]s_source"

  datastore {
    type    = "cassandra"
    version = "3.11"
    mode    = "CloudNativeCluster"
  }

  values = {
    request_timeout_in_ms = "20000"
  }
}

resource "huaweicloud_geminidb_parameter_template" "target" {
  name = "%[1]s_target"

  datastore {
    type    = "cassandra"
    version = "3.11"
    mode    = "CloudNativeCluster"
  }

  values = {
    request_timeout_in_ms = "30000"
  }
}

resource "huaweicloud_geminidb_parameter_template_compare" "test" {
  source_configuration_id = huaweicloud_geminidb_parameter_template.source.id
  target_configuration_id = huaweicloud_geminidb_parameter_template.target.id
}
`, rName)
}

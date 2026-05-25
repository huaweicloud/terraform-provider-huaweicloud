package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBParameterTemplateReset_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBParameterTemplateReset_basic(rName),
			},
		},
	})
}

func testAccGeminiDBParameterTemplateReset_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_parameter_template" "test" {
  name = "%s"

  datastore {
    type    = "cassandra"
    version = "3.11"
    mode    = "CloudNativeCluster"
  }

  values = {
    request_timeout_in_ms = "20000"
  }
}

resource "huaweicloud_geminidb_parameter_template_reset" "test" {
  config_id = huaweicloud_geminidb_parameter_template.test.id
}
`, rName)
}

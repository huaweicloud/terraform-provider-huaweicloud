package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// lintignore:AT001
func TestAccFunctionTopping_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTopping_basic(),
			},
		},
	})
}

func testAccFunctionTopping_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "js_script_content" {
  default = <<EOT
exports.handler = async (event, context) => {
    const result =
    {
        'repsonse_code': 200,
        'headers':
        {
            'Content-Type': 'application/json'
        },
        'isBase64Encoded': false,
        'body': JSON.stringify(event)
    }
    return result
}
EOT
}

resource "huaweicloud_fgs_function" "test" {
  count = 3

  name        = format("%[1]s_%%d", count.index) 
  app         = "default"
  agency      = "%[2]s"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  code_type   = "inline"
  runtime     = "Node.js12.13"
  func_code   = base64encode(jsonencode(var.js_script_content))
}

resource "huaweicloud_fgs_function_topping" "test" {
  depends_on = [huaweicloud_fgs_function.test]

  count = 3

  function_urn = huaweicloud_fgs_function.test[count.index].urn
}
`, name, acceptance.HW_FGS_AGENCY_NAME)
}

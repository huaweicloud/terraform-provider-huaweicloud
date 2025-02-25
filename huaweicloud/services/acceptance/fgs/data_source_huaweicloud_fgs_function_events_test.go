package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFunctionEvents_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_function_events.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFunctionEvents_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "events.#", "1"),
					// Check the attributes.
					resource.TestCheckResourceAttrPair(all, "events.0.id", "huaweicloud_fgs_function_event.test", "id"),
					resource.TestCheckResourceAttrPair(all, "events.0.name", "huaweicloud_fgs_function_event.test", "name"),
					resource.TestMatchResourceAttr(all, "events.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataFunctionEvents_basic() string {
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
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  code_type   = "inline"
  runtime     = "Node.js12.13"
  func_code   = base64encode(jsonencode(var.js_script_content))
}

resource "huaweicloud_fgs_function_event" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  name         = "%[1]s"
  content      = base64encode(jsonencode({"foo": "bar"}))
}

data "huaweicloud_fgs_function_events" "test" {
  depends_on = [
	# Query function events after function event create.
    huaweicloud_fgs_function_event.test
  ]

  function_urn = huaweicloud_fgs_function.test.urn
}
`, name)
}

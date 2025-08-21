package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFunctionTags_basic(t *testing.T) {
	var (
		base = "huaweicloud_fgs_function.test"

		all          = "data.huaweicloud_fgs_function_tags.all"
		dcForAllTags = acceptance.InitDataSourceCheck(all)

		byFunctionId   = "data.huaweicloud_fgs_function_tags.filter_by_function_id"
		dcByFunctionId = acceptance.InitDataSourceCheck(byFunctionId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFunctionTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllTags.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "function_id"),
					// Filter by function ID.
					dcByFunctionId.CheckResourceExists(),
					resource.TestCheckOutput("is_function_id_filter_useful", "true"),
					// Check the attributes.
					resource.TestCheckResourceAttrPair(byFunctionId, "function_id", base, "urn"),
					resource.TestCheckResourceAttr(byFunctionId, "tags.#", "2"),
					resource.TestCheckResourceAttr(byFunctionId, "tags.0.key", "foo"),
					resource.TestCheckResourceAttr(byFunctionId, "tags.0.value", "bar"),
					resource.TestCheckResourceAttr(byFunctionId, "tags.1.key", "key"),
					resource.TestCheckResourceAttr(byFunctionId, "tags.1.value", "value"),
				),
			},
		},
	})
}

func testAccDataFunctionTags_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "function_code_content" {
  type    = string
  default = <<EOT
def main():  
    print("Hello, World!")  

if __name__ == "__main__":  
    main()
EOT
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.function_code_content)

  tags = {
    foo = "bar"
    key = "value"
  }
}

# Without any filter parameter.
data "huaweicloud_fgs_function_tags" "all" {
  function_id = huaweicloud_fgs_function.test.id
}

# Filter by function ID.
locals {
  function_id = huaweicloud_fgs_function.test.id
}

data "huaweicloud_fgs_function_tags" "filter_by_function_id" {
  function_id = local.function_id
}

locals {
  function_id_filter_result = [for v in data.huaweicloud_fgs_function_tags.filter_by_function_id.tags[*].key :
    v == "foo" || v == "key"]
}

output "is_function_id_filter_useful" {
  value = length(local.function_id_filter_result) == 2 && alltrue(local.function_id_filter_result)
}
`, name)
}

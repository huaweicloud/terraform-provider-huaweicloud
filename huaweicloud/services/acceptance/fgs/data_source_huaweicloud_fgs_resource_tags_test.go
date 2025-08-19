package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataResourceTags_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_resource_tags.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataResourceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tags.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestMatchResourceAttr(all, "sys_tags.#", regexp.MustCompile(`[0-9]*`)),
				),
			},
		},
	})
}

func testAccDataResourceTags_basic() string {
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
  
data "huaweicloud_fgs_resource_tags" "all" {
  depends_on = [
    huaweicloud_fgs_function.test,
  ]
}
`, name)
}

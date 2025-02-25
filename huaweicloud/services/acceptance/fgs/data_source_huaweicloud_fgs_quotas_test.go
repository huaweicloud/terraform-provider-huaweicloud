package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataQuotas_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_quotas.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "quotas.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Check the attributes.
					resource.TestCheckResourceAttrSet(all, "quotas.0.type"),
					resource.TestCheckResourceAttrSet(all, "quotas.0.limit"),
					resource.TestCheckOutput("is_quota_fgs_func_num_used", "true"),
					resource.TestCheckOutput("is_just_fgs_func_code_size_contains_unit", "true"),
					resource.TestCheckOutput("is_unit_MB", "true"),
				),
			},
		},
	})
}

func testAccDataQuotas_basic() string {
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

data "huaweicloud_fgs_quotas" "all" {
  depends_on = [huaweicloud_fgs_function.test]
}

locals {
  all_used_quota_types = [for _, v in data.huaweicloud_fgs_quotas.all.quotas: v.type if v.used > 0]
}

output "is_quota_fgs_func_num_used" {
  value = length(local.all_used_quota_types) > 0 && contains(local.all_used_quota_types, "fgs_func_num")
}

locals {
  all_quotas_that_contains_unit = [for _, v in data.huaweicloud_fgs_quotas.all.quotas: v if v.unit != ""]
}

output "is_just_fgs_func_code_size_contains_unit" {
  value = length(local.all_quotas_that_contains_unit) == 1 && element(local.all_quotas_that_contains_unit, 0).type == "fgs_func_code_size"
}

output "is_unit_MB" {
  value = element(local.all_quotas_that_contains_unit, 0).unit == "MB"
}
`, name)
}

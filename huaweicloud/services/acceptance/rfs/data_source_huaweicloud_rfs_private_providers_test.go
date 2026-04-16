package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateProviders_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_private_providers.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceNameWithDash()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivateProviders_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "providers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "providers.0.provider_id"),
					resource.TestCheckResourceAttrSet(dataSource, "providers.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSource, "providers.0.provider_description"),
					resource.TestCheckResourceAttrSet(dataSource, "providers.0.provider_source"),
					resource.TestCheckResourceAttrSet(dataSource, "providers.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "providers.0.update_time"),
				),
			},
		},
	})
}

func testAccDataSourcePrivateProviders_base(name string) string {
	return fmt.Sprintf(`
variable "request_resp_print_script_content" {
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
  agency      = "function_all_trust"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  code_type   = "inline"
  runtime     = "Node.js12.13"
  func_code   = base64encode(var.request_resp_print_script_content)
}

resource "huaweicloud_rfs_private_provider" "test" {
  provider_name        = "%[1]s"
  function_graph_urn   = huaweicloud_fgs_function.test.urn
  provider_description = "provider description test"
  provider_version     = "1.1.1"
  version_description  = "version description test"
}
`, name)
}

func testAccDataSourcePrivateProviders_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_private_providers" "test" {
  depends_on = [huaweicloud_rfs_private_provider.test]

  sort_key = "create_time"
  sort_dir = "desc"
}
`, testAccDataSourcePrivateProviders_base(name))
}

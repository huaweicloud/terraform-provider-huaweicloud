package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateProviderVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_private_provider_versions.test"
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
				Config: testDataSourcePrivateProviderVersions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.provider_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.provider_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.provider_version"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.function_graph_urn"),
				),
			},
		},
	})
}

func testDataSourcePrivateProviderVersions_base(name string) string {
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
  provider_description = "acc private provider for version test"
  provider_version     = "1.0.0"
  version_description  = "initial version"
}

resource "huaweicloud_rfs_private_provider_version" "test" {
  provider_name       = huaweicloud_rfs_private_provider.test.provider_name
  provider_version    = "2.0.0"
  function_graph_urn  = huaweicloud_fgs_function.test.urn
  version_description = "private provider version test"
}
`, name)
}

func testDataSourcePrivateProviderVersions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rfs_private_provider_versions" "test" {
  provider_name = huaweicloud_rfs_private_provider.test.provider_name
  provider_id   = huaweicloud_rfs_private_provider.test.provider_id
  sort_key      = "create_time"
  sort_dir      = "desc"
}
`, testDataSourcePrivateProviderVersions_base(name), name)
}

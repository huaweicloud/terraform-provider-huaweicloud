package rfs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPrivateProviderVersionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region          = acceptance.HW_REGION_NAME
		product         = "rfs"
		httpUrl         = "v1/private-providers/{provider_name}/versions/{provider_version}/metadata"
		provideName     = state.Primary.Attributes["provider_name"]
		providerVersion = state.Primary.Attributes["provider_version"]
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RFS client: %s", err)
	}

	reqUUID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate RFS request UUID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", provideName)
	requestPath = strings.ReplaceAll(requestPath, "{provider_version}", providerVersion)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Client-Request-Id": reqUUID},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccPrivateProviderVersion_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_rfs_private_provider_version.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateProviderVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPrivateProviderVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "provider_name", "huaweicloud_rfs_private_provider.test", "provider_name"),
					resource.TestCheckResourceAttr(rName, "provider_version", "2.0.0"),
					resource.TestCheckResourceAttrPair(rName, "function_graph_urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(rName, "version_description", "private provider version test"),
					resource.TestCheckResourceAttrSet(rName, "provider_id"),
					resource.TestCheckResourceAttrSet(rName, "provider_source"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPrivateProviderVersionImportState(rName),
			},
		},
	})
}

func testPrivateProviderVersion_base(name string) string {
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
`, name)
}

func testPrivateProviderVersion_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rfs_private_provider_version" "test" {
  provider_name       = huaweicloud_rfs_private_provider.test.provider_name
  provider_version    = "2.0.0"
  function_graph_urn  = huaweicloud_fgs_function.test.urn
  version_description = "private provider version test"
}
`, testPrivateProviderVersion_base(name), name)
}

func testPrivateProviderVersionImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		providerName := rs.Primary.Attributes["provider_name"]
		providerVersion := rs.Primary.Attributes["provider_version"]
		if providerName == "" || providerVersion == "" {
			return "", fmt.Errorf("the provider_name (%s) or provider_version(%s) is nil",
				providerName, providerVersion)
		}

		return providerName + "/" + providerVersion, nil
	}
}

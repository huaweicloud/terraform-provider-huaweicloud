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

func getPrivateProviderResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "rfs"
		httpUrl = "v1/private-providers/{provider_name}/metadata"
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
	requestPath = strings.ReplaceAll(requestPath, "{provider_name}", state.Primary.ID)
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

func TestAccPrivateProvider_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_rfs_private_provider.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPrivateProviderResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateProvider_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "provider_name", name),
					resource.TestCheckResourceAttr(rName, "provider_description", "provider description test"),
					resource.TestCheckResourceAttrPair(rName, "function_graph_urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(rName, "provider_version", "1.1.1"),
					resource.TestCheckResourceAttr(rName, "version_description", "version description test"),
					resource.TestCheckResourceAttrSet(rName, "provider_id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testAccPrivateProvider_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "provider_description", "provider description test update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"function_graph_urn",
					"provider_version",
					"version_description",
				},
			},
		},
	})
}

func testPrivateProvider_base(name string) string {
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
`, name)
}

func testAccPrivateProvider_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rfs_private_provider" "test" {
  provider_name        = "%[2]s"
  function_graph_urn   = huaweicloud_fgs_function.test.urn
  provider_description = "provider description test"
  provider_version     = "1.1.1"
  version_description  = "version description test"
}
`, testPrivateProvider_base(name), name)
}

func testAccPrivateProvider_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rfs_private_provider" "test" {
  provider_name        = "%[2]s"
  function_graph_urn   = huaweicloud_fgs_function.test.urn
  provider_description = "provider description test update"
  provider_version     = "1.1.1"
  version_description  = "version description test"
}
`, testPrivateProvider_base(name), name)
}

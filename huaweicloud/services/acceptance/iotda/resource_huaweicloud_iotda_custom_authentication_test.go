package iotda

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getCustomAuthenticationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-authorizers/{authorizer_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{authorizer_id}", state.Primary.ID)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IoTDA custom authentication: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccCustomAuthentication_basic(t *testing.T) {
	var (
		authObj    interface{}
		rName      = "huaweicloud_iotda_custom_authentication.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&authObj,
		getCustomAuthenticationFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
			acceptance.TestAccPreCheckIOTDASigningPublicKey(t)
			acceptance.TestAccPreCheckIOTDASigningToken(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomAuthentication_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "authorizer_name", name),
					resource.TestCheckResourceAttrPair(rName, "func_urn", "huaweicloud_fgs_function.test1", "urn"),
					resource.TestCheckResourceAttr(rName, "signing_enable", "true"),
					resource.TestCheckResourceAttr(rName, "signing_token", acceptance.HW_IOTDA_SIGNING_TOKEN),
					resource.TestCheckResourceAttr(rName, "signing_public_key", acceptance.HW_IOTDA_SIGNING_PUBLIC_KEY),
					resource.TestCheckResourceAttr(rName, "default_authorizer", "true"),
					resource.TestCheckResourceAttr(rName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(rName, "cache_enable", "true"),
					resource.TestCheckResourceAttrSet(rName, "func_name"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testAccCustomAuthentication_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "authorizer_name", updateName),
					resource.TestCheckResourceAttrPair(rName, "func_urn", "huaweicloud_fgs_function.test2", "urn"),
					resource.TestCheckResourceAttr(rName, "signing_enable", "false"),
					resource.TestCheckResourceAttr(rName, "default_authorizer", "false"),
					resource.TestCheckResourceAttr(rName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(rName, "cache_enable", "false"),
					resource.TestCheckResourceAttrSet(rName, "func_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const functionScriptVariableDefinition = `
variable "script_content" {
  type    = string
  default = <<EOT
def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
EOT
}
`

func testAccCustomAuthentication_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function" "test1" {
  name        = "%[2]s-func1"
  memory_size = 128
  runtime     = "Python2.7"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
}

resource "huaweicloud_fgs_function" "test2" {
  name        = "%[2]s-fun2"
  memory_size = 128
  runtime     = "Python3.6"
  timeout     = 3
  app         = "default"
  handler     = "index.handler"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
}
`, functionScriptVariableDefinition, name)
}

func testAccCustomAuthentication_basic(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_custom_authentication" "test" {
  authorizer_name    = "%[3]s"
  func_urn           = huaweicloud_fgs_function.test1.urn
  signing_enable     = true
  signing_token      = "%[4]s"
  signing_public_key = "%[5]s"
  default_authorizer = true
  status             = "ACTIVE"
  cache_enable       = true
}
`, buildIoTDAEndpoint(), testAccCustomAuthentication_base(name), name, acceptance.HW_IOTDA_SIGNING_TOKEN, acceptance.HW_IOTDA_SIGNING_PUBLIC_KEY)
}

func testAccCustomAuthentication_update(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_custom_authentication" "test" {
  authorizer_name    = "%[3]s"
  func_urn           = huaweicloud_fgs_function.test2.urn
  signing_enable     = false
  signing_token      = "%[4]s"
  signing_public_key = "%[5]s"
  default_authorizer = false
  status             = "INACTIVE"
  cache_enable       = false
}
`, buildIoTDAEndpoint(), testAccCustomAuthentication_base(name), name, acceptance.HW_IOTDA_SIGNING_TOKEN, acceptance.HW_IOTDA_SIGNING_PUBLIC_KEY)
}

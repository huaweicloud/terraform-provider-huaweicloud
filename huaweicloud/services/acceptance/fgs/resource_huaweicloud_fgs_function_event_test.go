package fgs

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getFunctionEventResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/fgs/functions/{function_urn}/events/{event_id}"
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{function_urn}", state.Primary.Attributes["function_urn"])
	getPath = strings.ReplaceAll(getPath, "{event_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting FunctionGraph function event: %s", err)
	}
	return utils.FlattenResponse(requestResp)
}

func TestAccFunctionEvent_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName    = "huaweicloud_fgs_function_event.test"
		name            = acceptance.RandomAccResourceName()
		eventContent    = base64.StdEncoding.EncodeToString([]byte("{\"foo\": \"bar\"}"))
		newEventContent = base64.StdEncoding.EncodeToString([]byte("{\"key\": \"value\"}"))

		rc = acceptance.InitResourceCheck(resourceName, &obj, getFunctionEventResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// There is no special requirements for the permissions of the delegate. It only needs a agency name.
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionEvent_basic(name, eventContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "content", eventContent),
					resource.TestCheckResourceAttr(resourceName, "updated_at", ""),
				),
			},
			{
				Config: testAccFunctionEvent_basic(name, newEventContent),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn", "huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "content", newEventContent),
					resource.TestMatchResourceAttr(resourceName, "updated_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionEventImportStateFunc(resourceName),
			},
		},
	})
}

func testAccFunctionEventImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var functionUrn, eventName string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of function event is not found in the tfstate", rsName)
		}
		functionUrn = rs.Primary.Attributes["function_urn"]
		eventName = rs.Primary.Attributes["name"]
		if functionUrn == "" || eventName == "" {
			return "", fmt.Errorf("the value of function URN or event name is empty")
		}
		return fmt.Sprintf("%s/%s", functionUrn, eventName), nil
	}
}

func testAccFunctionEvent_basic(name, funcCode string) string {
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
  agency      = "%[2]s"
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
  content      = "%[3]s"
}
`, name, acceptance.HW_FGS_AGENCY_NAME, funcCode)
}

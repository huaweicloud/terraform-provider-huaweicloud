package coc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPublicScriptExecuteResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	getScriptExecuteHttpUrl := "v1/job/script/orders/{id}"
	getScriptExecutePath := client.Endpoint + getScriptExecuteHttpUrl
	getScriptExecutePath = strings.ReplaceAll(getScriptExecutePath, "{id}", state.Primary.ID)

	getScriptExecuteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getScriptExecuteResp, err := client.Request("GET", getScriptExecutePath, &getScriptExecuteOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving COC script execute: %s", err)
	}

	respBody, err := utils.FlattenResponse(getScriptExecuteResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving COC script execute: %s", err)
	}

	return respBody, nil
}

func TestAccResourceCocPublicScriptExecute_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_coc_public_script_execute.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPublicScriptExecuteResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocPublicScriptExecute_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "gmt_created"),
					resource.TestCheckResourceAttrSet(resourceName, "gmt_finished"),
					resource.TestCheckResourceAttrSet(resourceName, "execute_costs"),
					resource.TestCheckResourceAttrSet(resourceName, "creator"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "script_name"),
					resource.TestCheckResourceAttrSet(resourceName, "current_execute_batch_index"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"execute_batches"},
			},
		},
	})
}

func testCocPublicScriptExecute_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_coc_public_scripts" "test" {
  name = "HWC.ECS.OSOps-switch-linux-ssh.sh"
}

locals {
  script_uuid = [for v in data.huaweicloud_coc_public_scripts.test.data[*].script_uuid : v if v != ""][0]
}

resource "huaweicloud_coc_public_script_execute" "test" {
  script_uuid  = local.script_uuid
  timeout      = 300
  success_rate = 100
  execute_user = "root"
  script_params {
    param_name  = "action"
    param_value = "stop"
  }
  execute_batches {
    batch_index = 1
    target_instances {
      resource_id        = "%s"
      region_id          = "cn-north-4"
      cloud_service_name = "ECS"
      custom_attributes {
        key   = "key"
        value = "value"
      }
    }
    rotation_strategy = "CONTINUE"
  }
}
`, acceptance.HW_COC_INSTANCE_ID)
}

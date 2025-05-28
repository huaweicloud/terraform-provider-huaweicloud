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

func getScriptExecuteResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

func TestAccScriptExecute_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_script_execute.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getScriptExecuteResourceFunc,
	)

	// lintignore:AT001
	// without CheckDestroy because the ticket ID always exits.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tesScriptExecute_basic(rName, acceptance.HW_COC_INSTANCE_ID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "script_name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "FINISHED"),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "script_id",
						"huaweicloud_coc_script.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "finished_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "parameters", "is_sync"},
			},
		},
	})
}

func TestAccScriptExecute_no_sync(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_script_execute.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getScriptExecuteResourceFunc,
	)

	// lintignore:AT001
	// without CheckDestroy because the ticket ID always exits.
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tesScriptExecute_no_sync(rName, acceptance.HW_COC_INSTANCE_ID),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "script_name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "FINISHED"),
					resource.TestCheckResourceAttr(resourceName, "parameters.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "script_id",
						"huaweicloud_coc_script.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "finished_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "parameters", "is_sync"},
			},
		},
	})
}

func tesScriptExecute_basic(name, instanceID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = "%s"
  timeout      = 600
  execute_user = "root"

  parameters {
    name  = "name"
    value = "somebody"
  }
  parameters {
    name  = "company"
    value = "HuaweiCloud"
  }
}`, tesScript_updated(name), instanceID)
}

func tesScriptExecute_no_sync(name, instanceID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = "%s"
  timeout      = 600
  execute_user = "root"
  is_sync      = false

  parameters {
    name  = "name"
    value = "somebody"
  }
  parameters {
    name  = "company"
    value = "HuaweiCloud"
  }
}`, tesScript_updated(name), instanceID)
}

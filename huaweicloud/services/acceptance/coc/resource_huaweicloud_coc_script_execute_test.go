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
	var (
		obj interface{}

		name                    = acceptance.RandomAccResourceName()
		withResourceProvider    = "huaweicloud_coc_script_execute.with_resource_provider"
		withoutResourceProvider = "huaweicloud_coc_script_execute.without_resource_provider"

		rcWithResourceProvider    = acceptance.InitResourceCheck(withResourceProvider, &obj, getScriptExecuteResourceFunc)
		rcWithoutResourceProvider = acceptance.InitResourceCheck(withoutResourceProvider, &obj, getScriptExecuteResourceFunc)
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
				Config: tesScriptExecute_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithResourceProvider.CheckResourceExists(),
					resource.TestCheckResourceAttr(withResourceProvider, "script_name", name),
					resource.TestCheckResourceAttr(withResourceProvider, "status", "FINISHED"),
					resource.TestCheckResourceAttr(withResourceProvider, "parameters.#", "2"),
					resource.TestCheckResourceAttrPair(withResourceProvider, "script_id",
						"huaweicloud_coc_script.test", "id"),
					resource.TestCheckResourceAttrSet(withResourceProvider, "created_at"),
					resource.TestCheckResourceAttrSet(withResourceProvider, "finished_at"),

					rcWithoutResourceProvider.CheckResourceExists(),
					resource.TestCheckResourceAttr(withoutResourceProvider, "script_name", name),
					resource.TestCheckResourceAttr(withoutResourceProvider, "status", "FINISHED"),
					resource.TestCheckResourceAttr(withoutResourceProvider, "parameters.#", "2"),
					resource.TestCheckResourceAttrPair(withoutResourceProvider, "script_id",
						"huaweicloud_coc_script.test", "id"),
					resource.TestCheckResourceAttrSet(withoutResourceProvider, "created_at"),
					resource.TestCheckResourceAttrSet(withoutResourceProvider, "finished_at"),
				),
			},
			{
				ResourceName:            withResourceProvider,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "parameters", "is_sync", "resource_provider", "type"},
			},
			{
				ResourceName:            withoutResourceProvider,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "parameters", "is_sync", "resource_provider", "type"},
			},
		},
	})
}

func tesScriptExecute_basic_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_coc_script" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform acceptance test"
  risk_level            = "MEDIUM"
  version               = "1.0.1"
  type                  = "SHELL"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  content = <<EOF
#! /bin/bash
echo "hello $${name}@$${company}!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
  parameters {
    name        = "company"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func tesScriptExecute_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_script_execute" "with_resource_provider" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = "%[2]s"
  timeout      = 600
  execute_user = "root"

  resource_provider = "ecs"
  type              = "cloudservers"

  parameters {
    name  = "name"
    value = "somebody"
  }
  parameters {
    name  = "company"
    value = "HuaweiCloud"
  }
}

resource "huaweicloud_coc_script_execute" "without_resource_provider" {
  depends_on = [huaweicloud_coc_script_execute.with_resource_provider]

  script_id    = huaweicloud_coc_script.test.id
  instance_id  = "%[2]s"
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
}
`, tesScriptExecute_basic_base(name), acceptance.HW_COC_INSTANCE_ID)
}

func TestAccScriptExecute_no_sync(t *testing.T) {
	var (
		obj interface{}

		name                    = acceptance.RandomAccResourceName()
		withResourceProvider    = "huaweicloud_coc_script_execute.with_resource_provider"
		withoutResourceProvider = "huaweicloud_coc_script_execute.without_resource_provider"

		rcWithResourceProvider    = acceptance.InitResourceCheck(withResourceProvider, &obj, getScriptExecuteResourceFunc)
		rcWithoutResourceProvider = acceptance.InitResourceCheck(withoutResourceProvider, &obj, getScriptExecuteResourceFunc)
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
				Config: tesScriptExecute_no_sync(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithResourceProvider.CheckResourceExists(),
					resource.TestCheckResourceAttr(withResourceProvider, "script_name", name),
					resource.TestCheckResourceAttr(withResourceProvider, "status", "FINISHED"),
					resource.TestCheckResourceAttr(withResourceProvider, "parameters.#", "2"),
					resource.TestCheckResourceAttrPair(withResourceProvider, "script_id",
						"huaweicloud_coc_script.test", "id"),
					resource.TestCheckResourceAttrSet(withResourceProvider, "created_at"),
					resource.TestCheckResourceAttrSet(withResourceProvider, "finished_at"),

					rcWithoutResourceProvider.CheckResourceExists(),
					resource.TestCheckResourceAttr(withoutResourceProvider, "script_name", name),
					resource.TestCheckResourceAttr(withoutResourceProvider, "status", "FINISHED"),
					resource.TestCheckResourceAttr(withoutResourceProvider, "parameters.#", "2"),
					resource.TestCheckResourceAttrPair(withoutResourceProvider, "script_id",
						"huaweicloud_coc_script.test", "id"),
					resource.TestCheckResourceAttrSet(withoutResourceProvider, "created_at"),
					resource.TestCheckResourceAttrSet(withoutResourceProvider, "finished_at"),
				),
			},
			{
				ResourceName:            withResourceProvider,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "parameters", "is_sync", "resource_provider", "type"},
			},
			{
				ResourceName:            withoutResourceProvider,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id", "parameters", "is_sync", "resource_provider", "type"},
			},
		},
	})
}

func tesScriptExecute_no_sync(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_script_execute" "with_resource_provider" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = "%[2]s"
  timeout      = 600
  execute_user = "root"
  is_sync      = false

  resource_provider = "ecs"
  type              = "cloudservers"

  parameters {
    name  = "name"
    value = "somebody"
  }
  parameters {
    name  = "company"
    value = "HuaweiCloud"
  }
}

resource "huaweicloud_coc_script_execute" "without_resource_provider" {
  depends_on = [huaweicloud_coc_script_execute.with_resource_provider]

  script_id    = huaweicloud_coc_script.test.id
  instance_id  = "%[2]s"
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
}
`, tesScriptExecute_basic_base(name), acceptance.HW_COC_INSTANCE_ID)
}

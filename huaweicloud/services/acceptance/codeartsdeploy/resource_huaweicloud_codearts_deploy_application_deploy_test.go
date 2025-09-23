package codeartsdeploy

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeployApplicationDeploy_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployApplicationDeploy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					// wait deployment task become failed, lest the application can not be deleted
					func(_ *terraform.State) error {
						// lintignore:R018
						time.Sleep(60 * time.Second)

						return nil
					},
				),
			},
		},
	})
}

//nolint:revive
func testAccDeployApplicationDeploy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id     = huaweicloud_codearts_project.test.id
  name           = "%[2]s"
  is_draft       = false
  create_type    = "template"
  trigger_source = "0"

  operation_list {
    name      = "Run Shell Commands"
    code      = "https://wukong-prod-cn-north-4.obs.cn-north-4.myhuaweicloud.com/extensions/devcloud2018/execute_shell_command/1.0.153/roles/main.zip"
    entrance  = "main.yml"
    version   = "1.0.153"
    module_id = "devcloud2018.execute_shell_command.execute_shell_command"
    params    = <<EOF
[
  {
    "name": "groupId",
    "label": "Environment",
    "displaySettings": {
      "DevCloud.ControlType": "DeploymentGroup"
    },
    "defaultDisplay": [{
      "displayName": "",
      "value": "",
      "host_count": 0,
      "os": "linux"
    }]
  }, 
  {
    "name": "shell_command",
    "label": "Shell Commands",
    "displaySettings": {
      "DevCloud.ControlType": "CodeText"
    },
    "defaultValue": "echo hello"
  }, 
  {
    "name": "faq_url",
    "label": "",
    "displaySettings": {
      "DevCloud.ControlType": "Hidden"
    },
    "defaultValue": "/deployman_faq_0030.html"
  }, 
  {
    "name": "controller_enabled",
    "label": "",
    "displaySettings": {
      "DevCloud.ControlType": "Hidden"
    },
    "defaultDisplay": [{
      "displayName": "Enable this action",
      "value": "1"
    }]
  }, 
  {
    "name": "controller_enabled_ignore_errors",
    "label": "",
    "displaySettings": {
      "DevCloud.ControlType": "Checkbox"
    },
    "defaultDisplay": [{
      "displayName": "Keep running on failure",
      "value": "0"
    }]
  }, 
  {
    "name": "visibleRule",
    "label": "",
    "displaySettings": {
      "DevCloud.ControlType": "Hidden"
    },
    "defaultValue": "{\"groupId\":{\"value\":true},\"shell_command\":{\"value\":true},\"faq_url\":{\"value\":true},\"controller_enabled\":{\"value\":true},\"controller_enabled_ignore_errors\":{\"value\":true},\"visibleRule\":{\"value\":true},\"controller_enabled_sudo\":{\"value\":true}}"
  }, 
  {
    "name": "controller_enabled_sudo",
    "label": "",
    "displaySettings": {
      "DevCloud.ControlType": "Checkbox"
    },
    "defaultDisplay": [{
      "displayName": "Execute this action with the sudo permission",
      "value": "0"
    }]
  }
]
EOF
  }
}

resource "huaweicloud_codearts_deploy_application_deploy" "test" {
  task_id = huaweicloud_codearts_deploy_application.test.task_id
}
`, testProject_basic(name), name)
}

func TestAccDeployApplicationDeploy_withParams(t *testing.T) {
	resourceName := "huaweicloud_codearts_deploy_application_deploy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCodeArtsDeploymentTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployApplicationDeploy_withParams(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDeployApplicationDeployImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"params",
				},
			},
		},
	})
}

func testAccDeployApplicationDeploy_withParams() string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_deploy_application_deploy" "test" {
  task_id = "%s"

  params {
    name  = "name"
    value = "value"
    type  = "enum"
  }
}
`, acceptance.HW_CODEARTS_DEPLOYMENT_TASK_ID)
}

// testDeployApplicationDeployImportState use to return an ID with format <task_id>/<id>
func testDeployApplicationDeployImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		taskId := rs.Primary.Attributes["task_id"]
		if taskId == "" {
			return "", fmt.Errorf("attribute (task_id) of resource (%s) not found: %s", name, rs)
		}

		return taskId + "/" + rs.Primary.ID, nil
	}
}

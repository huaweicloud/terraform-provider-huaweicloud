package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCodeArtsDeployApplicationCopy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application_copy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCodeArtsDeployApplicationCopy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "create_type", "template"),
					resource.TestCheckResourceAttr(rName, "steps.step1", "Download Package"),
					resource.TestCheckResourceAttr(rName, "is_disable", "false"),
					resource.TestCheckResourceAttr(rName, "permission_level", "instance"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "project_name"),
					resource.TestCheckResourceAttrSet(rName, "can_modify"),
					resource.TestCheckResourceAttrSet(rName, "can_delete"),
					resource.TestCheckResourceAttrSet(rName, "can_view"),
					resource.TestCheckResourceAttrSet(rName, "can_execute"),
					resource.TestCheckResourceAttrSet(rName, "can_copy"),
					resource.TestCheckResourceAttrSet(rName, "can_manage"),
					resource.TestCheckResourceAttrSet(rName, "can_create_env"),
					resource.TestCheckResourceAttrSet(rName, "task_id"),
					resource.TestCheckResourceAttrSet(rName, "task_name"),
				),
			},
			{
				Config: testCodeArtsDeployApplicationCopy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-copy"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "is_draft", "true"),
					resource.TestCheckResourceAttr(rName, "create_type", "template"),
					resource.TestCheckResourceAttr(rName, "steps.step1", "Download Package"),
					resource.TestCheckResourceAttr(rName, "is_disable", "false"),
					resource.TestCheckResourceAttr(rName, "trigger_source", "0"),
					resource.TestCheckResourceAttr(rName, "permission_level", "instance"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "project_name"),
					resource.TestCheckResourceAttrSet(rName, "can_modify"),
					resource.TestCheckResourceAttrSet(rName, "can_delete"),
					resource.TestCheckResourceAttrSet(rName, "can_view"),
					resource.TestCheckResourceAttrSet(rName, "can_execute"),
					resource.TestCheckResourceAttrSet(rName, "can_copy"),
					resource.TestCheckResourceAttrSet(rName, "can_manage"),
					resource.TestCheckResourceAttrSet(rName, "can_create_env"),
					resource.TestCheckResourceAttrSet(rName, "task_id"),
					resource.TestCheckResourceAttrSet(rName, "task_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source_app_id",
					"is_draft",
					"trigger_source",
					"operation_list",
					"group_id",
				},
			},
		},
	})
}

func TestAccCodeArtsDeployApplicationCopy_updateWhenCreating(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application_copy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCodeArtsDeployApplicationCopy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"source_app_id",
					"is_draft",
					"trigger_source",
					"operation_list",
					"group_id",
				},
			},
		},
	})
}

func testCodeArtsDeployApplicationCopy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_copy" "test" {
  source_app_id = huaweicloud_codearts_deploy_application.test.id
}
`, testDeployApplication_basic(name))
}

func testCodeArtsDeployApplicationCopy_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_copy" "test" {
  source_app_id = huaweicloud_codearts_deploy_application.test.id

  project_id       = huaweicloud_codearts_project.test.id
  name             = "%[2]s-copy"
  is_draft         = true
  create_type      = "template"
  trigger_source   = "0"
  is_disable       = false
  group_id         = "no_grouped"
  permission_level = "instance"

  operation_list {
    name        = "Download Package"
    description = "download package description update"
    code        = "https://example.com/xxx.zip"
    entrance    = "main.yml"
    version     = "1.1.282"
    module_id   = "devcloud2018.select_deploy_source_task.select_deploy_source_tab"
    params      = <<EOF
[
  {
    "name":"groupId",
    "label":"env",
    "displaySettings":{
      "DevCloud.ControlType":"DeploymentGroup",
      "DevCloud.ControlType.Select":[
        {
          "displayName":"",
          "value":""
        }
      ]
    },
    "defaultDisplay":[
      {
        "displayName":"$${host_group}",
        "value":"$${host_group}",
        "os":"linux"
      }
    ]
  }
]
EOF
  }
}
`, testDeployApplication_basic(name), name)
}

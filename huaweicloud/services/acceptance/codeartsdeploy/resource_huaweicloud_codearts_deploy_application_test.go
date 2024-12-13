package codeartsdeploy

import (
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

func getDeployApplicationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/applications/{app_id}/info"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{app_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CodeArts deploy application: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	errorCode := utils.PathSearch("error_code", getRespBody, "")
	if errorCode == "Deploy.00011021" {
		// 'Deploy.00011021' means the application is not exist
		return nil, golangsdk.ErrDefault404{}
	}

	if errorCode != "" {
		errorMsg := utils.PathSearch("error_msg", getRespBody, "")
		return nil, fmt.Errorf("error retrieving CodeArts deploy application: error code: %s, error message: %s",
			errorCode, errorMsg)
	}

	return getRespBody, nil
}

func TestAccDeployApplication_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployApplication_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "is_draft", "true"),
					resource.TestCheckResourceAttr(rName, "create_type", "template"),
					resource.TestCheckResourceAttr(rName, "trigger_source", "0"),
					resource.TestCheckResourceAttr(rName, "steps.step1", "Download Package"),
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
					"is_draft",
					"trigger_source",
					"operation_list",
				},
			},
		},
	})
}

func TestAccDeployApplication_resourcePoolId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCodeArtsDeployResourcePoolID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployApplication_resourcePoolId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "is_draft", "true"),
					resource.TestCheckResourceAttr(rName, "create_type", "template"),
					resource.TestCheckResourceAttr(rName, "trigger_source", "1"),
					resource.TestCheckResourceAttr(rName, "artifact_source_system", "CloudArtifact"),
					resource.TestCheckResourceAttr(rName, "artifact_type", "generic"),
					resource.TestCheckResourceAttr(rName, "resource_pool_id", acceptance.HW_CODEARTS_RESOURCE_POOL_ID),
					resource.TestCheckResourceAttr(rName, "steps.step1", "Download Package"),
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
					"is_draft",
					"trigger_source",
					"artifact_source_system",
					"artifact_type",
					"operation_list",
				},
			},
		},
	})
}

func TestAccDeployApplication_errorCheckInvalidArgument(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testDeployApplication_errorCheckInvalidArgument(name),
				ExpectError: regexp.MustCompile(`is required when application is not in draft status`),
			},
		},
	})
}

func TestAccDeployApplication_errorCheckErrorCode(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testDeployApplication_errorCheckErrorCode(name),
				ExpectError: regexp.MustCompile(`error creating CodeArts deploy application: error code:`),
			},
		},
	})
}

func testDeployApplication_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id     = huaweicloud_codearts_project.test.id
  name           = "%[2]s"
  description    = "test description"
  is_draft       = true
  create_type    = "template"
  trigger_source = "0"

  operation_list {
    name        = "Download Package"
    description = "download package description"
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
`, testProject_basic(name), name)
}

func testDeployApplication_resourcePoolId(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id             = huaweicloud_codearts_project.test.id
  name                   = "%[2]s"
  description            = "test description"
  is_draft               = true
  create_type            = "template"
  trigger_source         = "1"
  artifact_source_system = "CloudArtifact"
  artifact_type          = "generic"
  resource_pool_id       = "%[3]s"

  operation_list {
    name        = "Download Package"
    description = "download package description"
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
`, testProject_basic(name), name, acceptance.HW_CODEARTS_RESOURCE_POOL_ID)
}

func testDeployApplication_errorCheckInvalidArgument(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id     = huaweicloud_codearts_project.test.id
  name           = "%[2]s"
  description    = "test description"
  is_draft       = false
  create_type    = "template"
  trigger_source = "0"
}
`, testProject_basic(name), name)
}

func testDeployApplication_errorCheckErrorCode(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id     = huaweicloud_codearts_project.test.id
  name           = "%[2]s"
  description    = "test description"
  is_draft       = true
  create_type    = "error_type"
  trigger_source = "0"
}
`, testProject_basic(name), name)
}

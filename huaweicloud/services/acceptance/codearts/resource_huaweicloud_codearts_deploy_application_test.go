package codearts

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

func getDeployApplicationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/tasks/{task_id}"
		product = "codearts_deploy"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{task_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
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
	if errorCode == "Deploy.00011020" {
		// 'Deploy.00011020' means the application is not exist
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
			acceptance.TestAccPreCheckCodeArtsDeployTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployApplication_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "project_name", "huaweicloud_codearts_project.test", "name"),
					resource.TestCheckResourceAttr(rName, "template_id", acceptance.HW_CODEARTS_DEPLOY_TEMPLATE_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "configs.#", "4"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestCheckResourceAttrSet(rName, "description"),
					resource.TestCheckResourceAttrSet(rName, "owner_name"),
					resource.TestCheckResourceAttrSet(rName, "owner_id"),
					resource.TestCheckResourceAttrSet(rName, "can_modify"),
					resource.TestCheckResourceAttrSet(rName, "can_delete"),
					resource.TestCheckResourceAttrSet(rName, "can_view"),
					resource.TestCheckResourceAttrSet(rName, "can_execute"),
					resource.TestCheckResourceAttrSet(rName, "can_copy"),
					resource.TestCheckResourceAttrSet(rName, "can_manage"),
					resource.TestCheckResourceAttrSet(rName, "role_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"configs",
					"template_id",
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
			acceptance.TestAccPreCheckCodeArtsDeployTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployApplication_resourcePoolId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "project_name", "huaweicloud_codearts_project.test", "name"),
					resource.TestCheckResourceAttr(rName, "template_id", acceptance.HW_CODEARTS_DEPLOY_TEMPLATE_ID),
					resource.TestCheckResourceAttr(rName, "resource_pool_id", acceptance.HW_CODEARTS_RESOURCE_POOL_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestCheckResourceAttrSet(rName, "description"),
					resource.TestCheckResourceAttrSet(rName, "owner_name"),
					resource.TestCheckResourceAttrSet(rName, "owner_id"),
					resource.TestCheckResourceAttrSet(rName, "can_modify"),
					resource.TestCheckResourceAttrSet(rName, "can_delete"),
					resource.TestCheckResourceAttrSet(rName, "can_view"),
					resource.TestCheckResourceAttrSet(rName, "can_execute"),
					resource.TestCheckResourceAttrSet(rName, "can_copy"),
					resource.TestCheckResourceAttrSet(rName, "can_manage"),
					resource.TestCheckResourceAttrSet(rName, "role_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"configs",
					"template_id",
				},
			},
		},
	})
}

func testDeployApplication_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id   = huaweicloud_codearts_project.test.id
  project_name = huaweicloud_codearts_project.test.name
  template_id  = "%s"
  name         = "%s"

  configs {
    name          = "text_param"
    type          = "text"
    description   = "text param description"
    value         = "text value"
    static_status = 0
  }

  configs {
    name        = "host_group_param"
    type        = "host_group"
    description = "host group param description"
    value       = "host group value"
  }

  configs {
    name        = "enum_param"
    type        = "enum"
    description = "enum param description"
    value       = "Monday"
    limits      = ["Monday", "Tuesday", "Wednesday", "Thursday"]
  }

  configs {
    name        = "encrypt_param"
    type        = "encrypt"
    description = "encrypt param description"
    value       = "encrypt value"
  }
}
`, testProject_basic(name), acceptance.HW_CODEARTS_DEPLOY_TEMPLATE_ID, name)
}

func testDeployApplication_resourcePoolId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_application" "test" {
  project_id       = huaweicloud_codearts_project.test.id
  project_name     = huaweicloud_codearts_project.test.name
  template_id      = "%s"
  resource_pool_id = "%s"
  name             = "%s"
}
`, testProject_basic(name), acceptance.HW_CODEARTS_DEPLOY_TEMPLATE_ID, acceptance.HW_CODEARTS_RESOURCE_POOL_ID, name)
}

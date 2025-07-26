package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartspipeline"
)

func getPipelinePermissionsResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := state.Primary.Attributes["project_id"]
	pipelineId := state.Primary.Attributes["pipeline_id"]
	userId := state.Primary.Attributes["user_id"]
	roleId := state.Primary.Attributes["role_id"]

	var rst interface{}
	if userId != "" {
		rst, err = codeartspipeline.GetPipelineUesrPermissions(client, projectId, pipelineId, userId)
		if err != nil {
			return nil, err
		}
	} else {
		rst, err = codeartspipeline.GetPipelineRolePermissions(client, projectId, pipelineId, roleId)
		if err != nil {
			return nil, err
		}
	}

	return rst, nil
}

func TestAccPipelinePermission_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_permission.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelinePermissionsResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelinePermission_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "role_id", "3"),
					resource.TestCheckResourceAttr(rName, "operation_delete", "true"),
					resource.TestCheckResourceAttr(rName, "operation_execute", "true"),
					resource.TestCheckResourceAttr(rName, "operation_query", "true"),
					resource.TestCheckResourceAttr(rName, "operation_update", "true"),
				),
			},
			{
				Config: testPipelinePermission_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "role_id", "3"),
					resource.TestCheckResourceAttr(rName, "operation_delete", "false"),
					resource.TestCheckResourceAttr(rName, "operation_execute", "false"),
					resource.TestCheckResourceAttr(rName, "operation_query", "false"),
					resource.TestCheckResourceAttr(rName, "operation_update", "false"),
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

func testPipelinePermission_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_permission" "test" {
  project_id        = huaweicloud_codearts_project.test.id
  pipeline_id       = huaweicloud_codearts_pipeline.test.id
  role_id           = 3
  operation_delete  = true
  operation_execute = true
  operation_query   = true
  operation_update  = true
}
`, testPipeline_basic(name))
}

func testPipelinePermission_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_permission" "test" {
  project_id        = huaweicloud_codearts_project.test.id
  pipeline_id       = huaweicloud_codearts_pipeline.test.id
  role_id           = 3
  operation_delete  = false
  operation_execute = false
  operation_query   = false
  operation_update  = false
}
`, testPipeline_basic(name))
}

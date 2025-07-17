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

func getPipelineMicroServiceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	return codeartspipeline.GetPipelineMicroService(client, state.Primary.Attributes["project_id"], state.Primary.ID)
}

func TestAccPipelineMicroService_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_micro_service.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineMicroServiceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPipelineMicroService_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "is_followed", "true"),
					resource.TestCheckResourceAttr(rName, "description", "demo"),
					resource.TestCheckResourceAttr(rName, "type", "microservice"),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
				),
			},
			{
				Config: testPipelineMicroService_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "is_followed", "false"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "type", "microservice"),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPipelineImportState(rName),
			},
		},
	})
}

func testPipelineMicroService_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_micro_service" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[2]s"
  type        = "microservice"
  description = "demo"
  is_followed = true
}
`, testProject_basic(name), name)
}

func testPipelineMicroService_update(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_codearts_pipeline_micro_service" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[3]s"
  type        = "microservice"
  is_followed = false

  repos {
    type     = "codehub"
    repo_id  = huaweicloud_codearts_repository.test.repository_id
    http_url = huaweicloud_codearts_repository.test.https_url
    git_url  = huaweicloud_codearts_repository.test.ssh_url
    branch   = "master"
    language = "java"
  }
}
`, testProject_basic(name), testRepository_basic(name), name)
}

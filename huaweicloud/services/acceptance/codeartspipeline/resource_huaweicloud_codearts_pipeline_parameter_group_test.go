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

func getPipelineParameterGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	return codeartspipeline.GetPipelineParameterGroup(client, state.Primary.Attributes["project_id"], state.Primary.ID)
}

func TestAccPipelineParameterGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_parameter_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineParameterGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPipelineParameterGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "test"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPipelineImportState(rName),
			},
			{
				Config: testPipelineParameterGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
		},
	})
}

func testPipeline_parameterGroup(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_pipeline_parameter_group" "test" {
  project_id  = huaweicloud_codearts_project.test.id
  name        = "%[1]s"
  description = "test"

  variables {
    description = "test"
    is_secret   = false
    name        = "test"
    sequence    = 1
    type        = "string"
    value       = "1"
  }
}
`, name)
}

func testPipelineParameterGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s
`, testProject_basic(name), testPipeline_parameterGroup(name))
}

func testPipelineParameterGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_parameter_group" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s_update"
}
`, testProject_basic(name), name)
}

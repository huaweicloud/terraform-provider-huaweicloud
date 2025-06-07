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

func getPipelineServiceEndpointResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	return codeartspipeline.GetPipelineServiceEndpoint(client, acceptance.HW_REGION_NAME,
		state.Primary.Attributes["project_id"], state.Primary.ID)
}

func TestAccPipelineServiceEndpoint_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_service_endpoint.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineServiceEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPipelineServiceEndpoint_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "created_by.#"),
					resource.TestCheckResourceAttrSet(rName, "created_by.0.user_id"),
					resource.TestCheckResourceAttrSet(rName, "created_by.0.user_name"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testPipelineServiceEndpointImportState(rName),
				ImportStateVerifyIgnore: []string{"authorization"},
			},
		},
	})
}

func testProject_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_codearts_project" "test" {
  name = "%s"
  type = "scrum"
}
`, name)
}

func testPipelineServiceEndpoint_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_service_endpoint" "test" {
  project_id = huaweicloud_codearts_project.test.id
  module_id  = "devcloud2018.codesource-codehub-https-24.oauth07"
  url        = "https://github.test.com/test"
  name       = "%[2]s"

  authorization {
    scheme     = "endpoint-auth-scheme-basic"
    parameters = jsonencode({
      "username":"test",
      "password":"test"
    })
  }
}
`, testProject_basic(name), name)
}

func testPipelineServiceEndpointImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["project_id"] == "" {
			return "", fmt.Errorf("attribute (project_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["project_id"] + "/" + rs.Primary.ID, nil
	}
}

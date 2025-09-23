package codeartspipeline

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartspipeline"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPipelineTagResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := codeartspipeline.GetPipelineTag(client, state.Primary.Attributes["project_id"])
	if err != nil {
		return nil, fmt.Errorf("error getting pipeline tags: %s", err)
	}

	searchPath := fmt.Sprintf("[?tag_id=='%s']|[0]", state.Primary.ID)
	tag := utils.PathSearch(searchPath, getRespBody, nil)
	if tag == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return tag, nil
}

func TestAccPipelineTag_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_tag.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineTagResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineTag_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "color", "#0b81f6"),
				),
			},
			{
				Config: testPipelineTag_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrPair(rName, "project_id", "huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "color", "#4eb15e"),
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

func testPipelineTag_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_tag" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s"
  color      = "#0b81f6"
}
`, testProject_basic(name), name)
}

func testPipelineTag_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_tag" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-update"
  color      = "#4eb15e"
}
`, testProject_basic(name), name)
}

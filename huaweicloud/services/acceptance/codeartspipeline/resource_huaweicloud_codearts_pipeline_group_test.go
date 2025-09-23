package codeartspipeline

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/codeartspipeline"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPipelineGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_pipeline", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts pipeline client: %s", err)
	}

	groups, err := codeartspipeline.GetPipelineGroups(client, state.Primary.Attributes["project_id"])
	if err != nil {
		return nil, fmt.Errorf("error retrieving pipeline groups: %s", err)
	}

	// filter results by path
	paths := strings.Split(state.Primary.Attributes["path_id"], ".")
	jsonPaths := fmt.Sprintf("[?id=='%s']", paths[0])
	for i, path := range paths {
		if i == 0 {
			continue
		}
		jsonPaths += fmt.Sprintf(".children[]|[?id=='%s']", path)
	}
	jsonPaths = fmt.Sprintf("%s|[0]", jsonPaths)

	group := utils.PathSearch(jsonPaths, groups, nil)
	if group == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return group, nil
}

func TestAccPipelineGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "path_id"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
				),
			},
			{
				Config: testPipelineGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrSet(rName, "path_id"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
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

func testPipelineGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_group" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s"
}
`, testProject_basic(name), name)
}

func testPipelineGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_group" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-update"
}
`, testProject_basic(name), name)
}

func TestAccPipelineGroup_secondLevel(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_pipeline_group.level2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPipelineGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPipelineGroup_secondLevel(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_codearts_pipeline_group.level1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-2"),
					resource.TestCheckResourceAttrSet(rName, "path_id"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
				),
			},
			{
				Config: testPipelineGroup_secondLevel_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_codearts_pipeline_group.level1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-2-update"),
					resource.TestCheckResourceAttrSet(rName, "path_id"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
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

func testPipelineGroup_secondLevel(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_group" "level1" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-1"
}

resource "huaweicloud_codearts_pipeline_group" "level2" {
  project_id = huaweicloud_codearts_project.test.id
  parent_id  = huaweicloud_codearts_pipeline_group.level1.id
  name       = "%[2]s-2"
}
`, testProject_basic(name), name)
}

func testPipelineGroup_secondLevel_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_pipeline_group" "level1" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-1"
}

resource "huaweicloud_codearts_pipeline_group" "level2" {
  project_id = huaweicloud_codearts_project.test.id
  parent_id  = huaweicloud_codearts_pipeline_group.level1.id
  name       = "%[2]s-2-update"
}
`, testProject_basic(name), name)
}

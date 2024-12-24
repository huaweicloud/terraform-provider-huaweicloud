package codeartsdeploy

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

func getDeployApplicationGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("codearts_deploy", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CodeArts deploy client: %s", err)
	}

	listHttpUrl := "v1/projects/{project_id}/applications/groups"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", state.Primary.Attributes["project_id"])
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving application groups: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening application groups: %s", err)
	}

	// filter results by path
	paths := strings.Split(state.Primary.Attributes["path"], ".")
	jsonPaths := fmt.Sprintf("result[?id=='%s']", paths[0])
	for i, path := range paths {
		if i == 0 {
			continue
		}
		jsonPaths += fmt.Sprintf(".children[]|[?id=='%s']", path)
	}
	jsonPaths = fmt.Sprintf("%s|[0]", jsonPaths)

	group := utils.PathSearch(jsonPaths, listRespBody, nil)
	if group == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return group, nil
}

func TestAccDeployApplicationGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployApplicationGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "path"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "application_count"),
				),
			},
			{
				Config: testDeployApplicationGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrSet(rName, "path"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "application_count"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDeployApplicationGroupImportState(rName),
			},
		},
	})
}

func testDeployApplicationGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_group" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s"
}
`, testProject_basic(name), name)
}

func testDeployApplicationGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_group" "test" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-update"
}
`, testProject_basic(name), name)
}

func TestAccDeployApplicationGroup_secondLevel(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_codearts_deploy_application_group.level2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeployApplicationGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeployApplicationGroup_secondLevel(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_codearts_deploy_application_group.level1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-2"),
					resource.TestCheckResourceAttrSet(rName, "path"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "application_count"),
				),
			},
			{
				Config: testDeployApplicationGroup_secondLevel_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "project_id",
						"huaweicloud_codearts_project.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_codearts_deploy_application_group.level1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name+"-2-update"),
					resource.TestCheckResourceAttrSet(rName, "path"),
					resource.TestCheckResourceAttrSet(rName, "ordinal"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestCheckResourceAttrSet(rName, "updated_by"),
					resource.TestCheckResourceAttrSet(rName, "application_count"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDeployApplicationGroupImportState(rName),
			},
		},
	})
}

func testDeployApplicationGroup_secondLevel(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_group" "level1" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-1"
}

resource "huaweicloud_codearts_deploy_application_group" "level2" {
  project_id = huaweicloud_codearts_project.test.id
  parent_id  = huaweicloud_codearts_deploy_application_group.level1.id
  name       = "%[2]s-2"
}
`, testProject_basic(name), name)
}

func testDeployApplicationGroup_secondLevel_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_codearts_deploy_application_group" "level1" {
  project_id = huaweicloud_codearts_project.test.id
  name       = "%[2]s-1"
}

resource "huaweicloud_codearts_deploy_application_group" "level2" {
  project_id = huaweicloud_codearts_project.test.id
  parent_id  = huaweicloud_codearts_deploy_application_group.level1.id
  name       = "%[2]s-2-update"
}
`, testProject_basic(name), name)
}

// testDeployApplicationGroupImportState use to return an ID with format <project_id>/<id>
func testDeployApplicationGroupImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		projectId := rs.Primary.Attributes["project_id"]
		if projectId == "" {
			return "", fmt.Errorf("attribute (project_id) of resource (%s) not found: %s", name, rs)
		}
		return projectId + "/" + rs.Primary.ID, nil
	}
}

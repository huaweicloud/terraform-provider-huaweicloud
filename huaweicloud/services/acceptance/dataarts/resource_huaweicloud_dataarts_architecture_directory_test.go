package dataarts

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

func getArchitectureDirectoryResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getDirectoryClient, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio V2 client: %s", err)
	}

	//nolint:misspell
	getDirectoryHttpUrl := "v2/{project_id}/design/directorys?type={type}"

	getDirectoryPath := getDirectoryClient.Endpoint + getDirectoryHttpUrl
	getDirectoryPath = strings.ReplaceAll(getDirectoryPath, "{project_id}", getDirectoryClient.ProjectID)
	getDirectoryPath = strings.ReplaceAll(getDirectoryPath, "{type}", state.Primary.Attributes["type"])

	getDirectoryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}
	getDirectoryResp, err := getDirectoryClient.Request("GET", getDirectoryPath, &getDirectoryOpt)
	if err != nil {
		return nil, err
	}

	getDirectoryRespBody, err := utils.FlattenResponse(getDirectoryResp)
	if err != nil {
		return nil, err
	}

	paths := strings.Split(state.Primary.Attributes["qualified_name"], ".")
	jsonPaths := fmt.Sprintf("data.value[?name=='%s']", paths[0])
	for i, path := range paths {
		if i == 0 {
			continue
		}
		jsonPaths += fmt.Sprintf("[children][][?name=='%s'][]", path)
	}

	directories := utils.PathSearch(jsonPaths, getDirectoryRespBody, make([]interface{}, 0)).([]interface{})
	if len(directories) > 0 {
		return directories, nil
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccResourceArchitectureDirectory_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_architecture_directory.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getArchitectureDirectoryResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureDirectory_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "STANDARD_ELEMENT"),
					resource.TestCheckResourceAttrSet(resourceName, "qualified_name"),
					resource.TestCheckResourceAttrSet(resourceName, "root_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccArchitectureDirectory_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "CODE"),
					resource.TestCheckResourceAttrSet(resourceName, "qualified_name"),
					resource.TestCheckResourceAttrSet(resourceName, "root_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceArchitectureDirectoryImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccResourceArchitectureDirectoryImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		dirType := rs.Primary.Attributes["type"]
		qualifiedName := rs.Primary.Attributes["qualified_name"]
		if workspaceID == "" || dirType == "" || qualifiedName == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<workspace_id>/<type>/<qualified_name>', but got '%s/%s/%s'",
				workspaceID, dirType, qualifiedName)
		}
		return fmt.Sprintf("%s/%s/%s", workspaceID, dirType, qualifiedName), nil
	}
}

func testAccArchitectureDirectory_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_directory" "test" {
  workspace_id = "%s"
  name         = "%s"
  type         = "STANDARD_ELEMENT"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccArchitectureDirectory_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_directory" "test" {
  workspace_id = "%s"
  name         = "%s"
  type         = "CODE"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

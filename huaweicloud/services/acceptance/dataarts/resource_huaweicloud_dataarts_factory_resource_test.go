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

func getFactoryResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getResourceClient, err := conf.NewServiceClient("dataarts-dlf", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getResourceHttpUrl := "v1/{project_id}/resources/{resource_id} "

	getResourcePath := getResourceClient.Endpoint + getResourceHttpUrl
	getResourcePath = strings.ReplaceAll(getResourcePath, "{project_id}", getResourceClient.ProjectID)
	getResourcePath = strings.ReplaceAll(getResourcePath, "{resource_id}", state.Primary.ID)

	getResourceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}
	getResourceResp, err := getResourceClient.Request("GET", getResourcePath, &getResourceOpt)
	if err != nil {
		return nil, err
	}

	getResourceRespBody, err := utils.FlattenResponse(getResourceResp)
	if err != nil {
		return nil, err
	}

	return getResourceRespBody, nil
}

func TestAccFactoryResource_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_factory_resource.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getFactoryResourceFunc,
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
				Config: testAccResource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "jar"),
					resource.TestCheckResourceAttr(resourceName, "directory", "/"),
					resource.TestCheckResourceAttr(resourceName, "location", "obs://test/main.jar"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "depend_packages.0.type", "jar"),
					resource.TestCheckResourceAttr(resourceName, "depend_packages.0.location", "obs://test/depend.jar"),
				),
			},
			{
				Config: testAccResource_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "type", "jar"),
					resource.TestCheckResourceAttr(resourceName, "directory", "/"),
					resource.TestCheckResourceAttr(resourceName, "location", "obs://test/main_update.jar"),
					resource.TestCheckResourceAttr(resourceName, "description", "test update"),
					resource.TestCheckResourceAttr(resourceName, "depend_packages.0.type", "file"),
					resource.TestCheckResourceAttr(resourceName, "depend_packages.0.location", "obs://test/depend_update.properties"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		id := rs.Primary.ID
		if workspaceID == "" || id == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<workspace_id>/<resource_id>', but got '%s/%s'",
				workspaceID, id)
		}
		return fmt.Sprintf("%s/%s", workspaceID, id), nil
	}
}

func testAccResource_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_resource" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  type         = "jar"
  directory    = "/"
  location     = "obs://test/main.jar"
  description  = "test"

  depend_packages {
    type     = "jar"
    location = "obs://test/depend.jar"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccResource_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_factory_resource" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_update"
  type         = "jar"
  directory    = "/"
  location     = "obs://test/main_update.jar"
  description  = "test update"

  depend_packages {
    type     = "file"
    location = "obs://test/depend_update.properties"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

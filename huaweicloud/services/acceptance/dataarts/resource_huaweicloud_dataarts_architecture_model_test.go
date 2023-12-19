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

func getModelResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	modelClient, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}
	readModelHttpUrl := "v2/{project_id}/design/workspaces/{model_id}"
	readModelPath := modelClient.Endpoint + readModelHttpUrl
	readModelPath = strings.ReplaceAll(readModelPath, "{project_id}", modelClient.ProjectID)
	readModelPath = strings.ReplaceAll(readModelPath, "{model_id}", state.Primary.ID)
	workspaceID := state.Primary.Attributes["workspace_id"]
	readModelOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	readModelResp, err := modelClient.Request("GET", readModelPath, &readModelOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(readModelResp)
}
func TestAccResourceModel_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_architecture_model.test"
	rName := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getModelResourceFunc,
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
				Config: testAccModel_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "THIRD_NF"),
					resource.TestCheckResourceAttr(resourceName, "dw_type", "DWS"),
					resource.TestCheckResourceAttr(resourceName, "physical", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccModel_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "type", "THIRD_NF"),
					resource.TestCheckResourceAttr(resourceName, "dw_type", "DWS"),
					resource.TestCheckResourceAttr(resourceName, "physical", "true"),
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
				ImportStateIdFunc: testAccResourceModelImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccResourceModelImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		name := rs.Primary.Attributes["name"]
		if workspaceID == "" || name == "" {
			return "", fmt.Errorf("invalid format specified for import ID, "+
				"want '<workspace_id>/<name>', but got '%s/%s'",
				workspaceID, name)
		}
		return fmt.Sprintf("%s/%s", workspaceID, name), nil
	}
}
func testAccModel_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_model" "test"{
  workspace_id = "%s"
  name         = "%s"
  type         = "THIRD_NF"
  physical     = true
  dw_type      = "DWS"
}`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccModel_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_model" "test"{
  workspace_id = "%s"
  name         = "%s"
  type         = "THIRD_NF"
  physical     = true
  dw_type      = "DWS"
  description  = "test description"
}`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

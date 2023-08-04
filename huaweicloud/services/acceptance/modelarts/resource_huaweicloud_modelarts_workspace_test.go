package modelarts

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

func getModelartsWorkspaceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getWorkspace: Query the Modelarts workspace.
	var (
		getWorkspaceHttpUrl = "v1/{project_id}/workspaces/{id}"
		getWorkspaceProduct = "modelarts"
	)
	getWorkspaceClient, err := cfg.NewServiceClient(getWorkspaceProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts Client: %s", err)
	}

	getWorkspacePath := getWorkspaceClient.Endpoint + getWorkspaceHttpUrl
	getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{project_id}", getWorkspaceClient.ProjectID)
	getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{id}", state.Primary.ID)

	getWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getWorkspaceResp, err := getWorkspaceClient.Request("GET", getWorkspacePath, &getWorkspaceOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts workspace: %s", err)
	}

	getWorkspaceRespBody, err := utils.FlattenResponse(getWorkspaceResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Modelarts workspace: %s", err)
	}

	return getWorkspaceRespBody, nil
}

func TestAccModelartsWorkspace_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_modelarts_workspace.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsWorkspaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsWorkspace_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "public"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accptance test"),
					resource.TestCheckResourceAttr(rName, "status", "NORMAL"),
				),
			},
			{
				Config: testModelartsWorkspace_basic_update(name, "private", "This is an accptance test update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "private"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accptance test update"),
					resource.TestCheckResourceAttr(rName, "status", "NORMAL"),
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

func testModelartsWorkspace_basic_update(name, authType, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_workspace" "test" {
  name        = "%s"
  auth_type   = "%s"
  description = "%s"
}
`, name, authType, description)
}

func testModelartsWorkspace_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_workspace" "test" {
  name        = "%s"
  description = "This is an accptance test"
}
`, name)
}

func TestAccModelartsWorkspace_internal(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_modelarts_workspace.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelartsWorkspaceResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelartsWorkspace_internal(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "internal"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accptance test"),
					resource.TestCheckResourceAttr(rName, "status", "NORMAL"),
					resource.TestCheckResourceAttrPair(rName, "grants.0.user_id", "huaweicloud_identity_user.user_1", "id"),
				),
			},
			{
				Config: testModelartsWorkspace_internal_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "auth_type", "internal"),
					resource.TestCheckResourceAttr(rName, "description", "This is an accptance test update"),
					resource.TestCheckResourceAttr(rName, "status", "NORMAL"),
					resource.TestCheckResourceAttr(rName, "grants.#", "0"),
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

func testModelartsWorkspace_internal(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "user_1" {
  name        = "%[1]s"
  password    = "password123@!"
  enabled     = true
  email       = "%[1]s@abc.com"
  description = "tested by terraform"
}

resource "huaweicloud_modelarts_workspace" "test" {
  name        = "%[1]s"
  description = "This is an accptance test"
  auth_type   = "internal"
  grants {
    user_id = huaweicloud_identity_user.user_1.id
  }
}
`, name)
}

func testModelartsWorkspace_internal_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_workspace" "test" {
  name        = "%[1]s"
  description = "This is an accptance test update"
  auth_type   = "internal"
}
`, name)
}

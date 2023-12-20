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

func getSecurityPermissionSetResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	getPermissionSetClient, err := conf.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	getPermissionSetHttpUrl := "v1/{project_id}/security/permission-sets/{permission_set_id}"
	getPermissionSetPath := getPermissionSetClient.Endpoint + getPermissionSetHttpUrl
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{project_id}", getPermissionSetClient.ProjectID)
	getPermissionSetPath = strings.ReplaceAll(getPermissionSetPath, "{permission_set_id}", state.Primary.ID)

	getPermissionSetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": state.Primary.Attributes["workspace_id"]},
	}
	getPermissionSetResp, err := getPermissionSetClient.Request("GET", getPermissionSetPath, &getPermissionSetOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Security permission set: %s", err)
	}

	getPermissionSetRespBody, err := utils.FlattenResponse(getPermissionSetResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DataArts Security permission set: %s", err)
	}

	return getPermissionSetRespBody, nil
}

func TestAccResourceSecurityPermissionSet_basic(t *testing.T) {
	var obj interface{}

	resourceName := "huaweicloud_dataarts_security_permission_set.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSecurityPermissionSetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsManagerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPermissionSet_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "parent_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "manager_id", acceptance.HW_DATAARTS_MANAGER_ID),
					resource.TestCheckResourceAttr(resourceName, "description", "test_create"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_by"),
				),
			},
			{
				Config: testAccSecurityPermissionSet_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "parent_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "manager_id", acceptance.HW_DATAARTS_MANAGER_ID),
					resource.TestCheckResourceAttr(resourceName, "description", "test_update"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
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
				ImportStateIdFunc: testAccResourceSecurityPermissionSetImportStateIDFunc(resourceName),
			},
		},
	})
}

func testAccResourceSecurityPermissionSetImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		permissionSetID := rs.Primary.ID
		if workspaceID == "" || permissionSetID == "" {
			return "", fmt.Errorf("invalid format specified for import ID, " +
				"want '<workspace_id>/<id>'")
		}
		return fmt.Sprintf("%s/%s", workspaceID, permissionSetID), nil
	}
}

func testAccSecurityPermissionSet_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_permission_set" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  parent_id    = "0"
  manager_id   = "%[3]s"
  description  = "test_create"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_MANAGER_ID)
}

func testAccSecurityPermissionSet_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_security_permission_set" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s_update"
  parent_id    = "0"
  manager_id   = "%[3]s"
  description  = "test_update"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name, acceptance.HW_DATAARTS_MANAGER_ID)
}

package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getSecurityPermissionSetMemberResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	workspaceID := state.Primary.Attributes["workspace_id"]
	permissionSetId := state.Primary.Attributes["permission_set_id"]
	objectId := state.Primary.Attributes["object_id"]
	return dataarts.GetMemberByObjectId(cfg, acceptance.HW_REGION_NAME, workspaceID, permissionSetId, objectId)
}

func TestAccResourceSecurityPermissionSetMember_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dataarts_security_permission_set_member.test"
	rName := acceptance.RandomAccResourceName()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSecurityPermissionSetMemberResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsManagerID(t)
			acceptance.TestAccPreCheckDataArtsSecurityPermissionSetMember(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityPermissionSetMember_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "object_id", acceptance.HW_DATAARTS_SECURITY_PERMISSSIONSET_MEMBER_OBJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "name", acceptance.HW_DATAARTS_SECURITY_PERMISSSIONSET_MEMBER_OBJECT_NAME),
					resource.TestCheckResourceAttr(resourceName, "type", "USER"),
					resource.TestCheckResourceAttrSet(resourceName, "member_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceSecurityPermissionSetMenmerImportStateFunc(resourceName),
			},
		},
	})
}

func testAccResourceSecurityPermissionSetMenmerImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		workspaceID := rs.Primary.Attributes["workspace_id"]
		permissionSetId := rs.Primary.Attributes["permission_set_id"]
		objectId := rs.Primary.Attributes["object_id"]
		if workspaceID == "" || permissionSetId == "" || objectId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<permission_set_id>/<object_id>', but got '%s/%s/%s'",
				workspaceID, permissionSetId, objectId)
		}
		return fmt.Sprintf("%s/%s/%s", workspaceID, permissionSetId, objectId), nil
	}
}

func testAccSecurityPermissionSetMember_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dataarts_security_permission_set_member" "test" {
  workspace_id      = "%s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  object_id         = "%s"
  name              = "%s"
  type              = "USER"
}
`, testAccSecurityPermissionSet_basic(name), acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_SECURITY_PERMISSSIONSET_MEMBER_OBJECT_ID,
		acceptance.HW_DATAARTS_SECURITY_PERMISSSIONSET_MEMBER_OBJECT_NAME)
}

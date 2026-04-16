package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataStudioWorkspaceUserRoles_basic(t *testing.T) {
	var (
		allRolesUnderWorkspace   = "data.huaweicloud_dataarts_studio_workspace_user_roles.all_roles_under_workspace"
		dcAllRolesUnderWorkspace = acceptance.InitDataSourceCheck(allRolesUnderWorkspace)
		allRolesUnderInstance    = "data.huaweicloud_dataarts_studio_workspace_user_roles.all_roles_under_instance"
		dcAllRolesUnderInstance  = acceptance.InitDataSourceCheck(allRolesUnderInstance)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataStudioWorkspaceUserRoles_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Studio workspace user roles"),
			},
			{
				Config: testAccDataStudioWorkspaceUserRoles_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcAllRolesUnderWorkspace.CheckResourceExists(),
					resource.TestMatchResourceAttr(allRolesUnderWorkspace, "roles.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(allRolesUnderWorkspace, "roles.0.id"),
					resource.TestCheckResourceAttrSet(allRolesUnderWorkspace, "roles.0.name"),
					resource.TestCheckResourceAttrSet(allRolesUnderWorkspace, "roles.0.code"),
					dcAllRolesUnderInstance.CheckResourceExists(),
					resource.TestMatchResourceAttr(allRolesUnderInstance, "roles.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(allRolesUnderInstance, "roles.0.id"),
					resource.TestCheckResourceAttrSet(allRolesUnderInstance, "roles.0.name"),
					resource.TestCheckResourceAttrSet(allRolesUnderInstance, "roles.0.code"),
				),
			},
		},
	})
}

func testAccDataStudioWorkspaceUserRoles_nonExistentWorkspace() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = "%[1]s"
}
`, randUUID)
}

func testAccDataStudioWorkspaceUserRoles_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_user_roles" "all_roles_under_workspace" {
  workspace_id = "%[1]s"
}

data "huaweicloud_dataarts_studio_workspace_user_roles" "all_roles_under_instance" {
  instance_id = "%[2]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_INSTANCE_ID)
}

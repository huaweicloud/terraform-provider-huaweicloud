package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeployGroupPermissionModify_basic(t *testing.T) {
	resourceName := "huaweicloud_codearts_deploy_group_permission.test"
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployGroupPermissionModify_basic(rName, false),
			},
			{
				Config: testAccDeployGroupPermissionModify_basic(rName, true),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDeployGroupPermissionImportState(resourceName),
			},
		},
	})
}

func testAccDeployGroupPermissionModify_basic(rName string, value bool) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_group_permission" "test" {
  project_id       = huaweicloud_codearts_deploy_group.test.project_id
  group_id         = huaweicloud_codearts_deploy_group.test.id
  role_id          = try(huaweicloud_codearts_deploy_group.test.permission_matrix[2].role_id, "")
  permission_name  = "can_add_host"
  permission_value = %t
}`, testDeployGroup_basic(rName), value)
}

// testDeployGroupPermissionImportState use to return an ID with format <project_id>/<id>
func testDeployGroupPermissionImportState(name string) resource.ImportStateIdFunc {
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

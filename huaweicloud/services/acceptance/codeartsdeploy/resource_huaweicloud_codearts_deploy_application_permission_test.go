package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeployApplicationPermissionModify_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployApplicationPermissionModify_basic(rName, false),
			},
			{
				Config: testAccDeployApplicationPermissionModify_basic(rName, true),
			},
		},
	})
}

func testAccDeployApplicationPermissionModify_basic(rName string, value bool) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_application_permission" "test" {
  project_id      = huaweicloud_codearts_deploy_application.test.project_id
  application_ids = [huaweicloud_codearts_deploy_application.test.id]

  roles {
    role_id        = try(huaweicloud_codearts_deploy_application.test.permission_matrix[2].role_id, "")
    can_modify     = %t
    can_disable    = true
    can_delete     = true
    can_view       = true
    can_execute    = true
    can_copy       = true
    can_manage     = true
    can_create_env = true
  }
}`, testDeployApplication_basic(rName), value)
}

func TestAccDeployApplicationPermissionModify_conflict(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployApplicationPermissionModify_conflict(rName),
			},
		},
	})
}

func testAccDeployApplicationPermissionModify_conflict(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_application_permission" "test" {
  count = 2 

  project_id      = huaweicloud_codearts_deploy_application.test.project_id
  application_ids = [huaweicloud_codearts_deploy_application.test.id]

  roles {
    role_id        = try(huaweicloud_codearts_deploy_application.test.permission_matrix[2].role_id, "")
    can_modify     = false
    can_disable    = true
    can_delete     = true
    can_view       = true
    can_execute    = true
    can_copy       = true
    can_manage     = true
    can_create_env = true
  }
}`, testDeployApplication_basic(rName))
}

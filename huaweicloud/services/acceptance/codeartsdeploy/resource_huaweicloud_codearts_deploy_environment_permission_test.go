package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeployEnvironmentPermissionModify_basic(t *testing.T) {
	resourceName := "huaweicloud_codearts_deploy_environment_permission.test"
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployEnvironmentPermissionModify_basic(rName, false),
			},
			{
				Config: testAccDeployEnvironmentPermissionModify_basic(rName, true),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDeployEnvironmentPermissionModify_basic(rName string, value bool) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_environment_permission" "test" {
  application_id   = huaweicloud_codearts_deploy_application.test.id
  environment_id   = huaweicloud_codearts_deploy_environment.test.id
  role_id          = try(huaweicloud_codearts_deploy_environment.test.permission_matrix[2].role_id, "")
  permission_name  = "can_delete"
  permission_value = %t
}`, testDeployEnvironment_basic(rName), value)
}

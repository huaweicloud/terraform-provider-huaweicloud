package codeartsdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeployGroupPermissionModify_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDeployGroupPermissionModify_basic(rName),
			},
		},
	})
}

func testAccDeployGroupPermissionModify_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_codearts_deploy_group_permission" "test" {
  project_id       = huaweicloud_codearts_deploy_group.test.project_id
  group_id         = huaweicloud_codearts_deploy_group.test.id
  role_id          = try(huaweicloud_codearts_deploy_group.test.permission_matrix[2].role_id, "")
  permission_name  = "can_add_host"
  permission_value = false
}`, testDeployGroup_basic(rName))
}

package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAppRepoBucketAssign_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppRepoBucketAssign,
			},
			{
				Config: testAccResourceAppRepoBucketAssign_withName(name),
			},
		},
	})
}

const testAccResourceAppRepoBucketAssign = `resource "huaweicloud_workspace_app_repo_bucket_assign" "test" {}`

func testAccResourceAppRepoBucketAssign_withName(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_repo_bucket_assign" "test_with_name" {
  bucket_name = "%[1]s"
}`, name)
}

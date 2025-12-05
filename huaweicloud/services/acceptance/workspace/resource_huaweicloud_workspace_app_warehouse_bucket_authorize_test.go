package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAppWarehouseBucketAuthorize_basic(t *testing.T) {
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
				Config: testAccResourceAppWarehouseBucketAuthorize,
			},
			{
				Config: testAccResourceAppWarehouseBucketAuthorize_withName(name),
			},
		},
	})
}

const testAccResourceAppWarehouseBucketAuthorize = `resource "huaweicloud_workspace_app_warehouse_bucket_authorize" "test" {}`

func testAccResourceAppWarehouseBucketAuthorize_withName(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_warehouse_bucket_authorize" "test_with_name" {
  bucket_name = "%[1]s"
}`, name)
}

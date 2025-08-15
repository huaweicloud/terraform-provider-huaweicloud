package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAppBucketAuthorize_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppBucketAuthorize_basic(),
				Check:  resource.TestCheckOutput("is_success", "true"),
			},
		},
	})
}

func testAccResourceAppBucketAuthorize_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_bucket_authorize" "test" {}

data "huaweicloud_identity_projects" "test" {
  name = "%[1]s"
}

locals {
  bucket_name = format("wks-app-%%s", data.huaweicloud_identity_projects.test.projects[0].id)
}

# Exactly match the bucket name.
data "huaweicloud_obs_buckets" "test" {
  depends_on = [huaweicloud_workspace_app_bucket_authorize.test]
  bucket     = local.bucket_name
}

output "is_success" {
  value = try(data.huaweicloud_obs_buckets.test.buckets[0].bucket == local.bucket_name, false)
}
`, acceptance.HW_REGION_NAME)
}

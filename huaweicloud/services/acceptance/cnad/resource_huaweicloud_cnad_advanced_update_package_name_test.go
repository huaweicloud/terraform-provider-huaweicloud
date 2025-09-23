package cnad

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case only verifies the expected error scenario.
func TestAccResourceUpdatePackageName_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testUpdatePackageName_basic,
				ExpectError: regexp.MustCompile(`Resource Not Found.|资源不存在。`),
			},
		},
	})
}

// The value of field `package_id` is mock data.
const testUpdatePackageName_basic = `
resource "huaweicloud_cnad_advanced_update_package_name" "test" {
  package_id = "1d8c03c4-a1d0-49cf-808a-ab50bfd7f0d8"
  name       = "test-package-name"
}`

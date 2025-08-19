package antiddos

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to the lack of a test environment, this test case only verifies the expected error message.
func TestAccResourceUnblockIp_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testUnblockIp_basic,
				ExpectError: regexp.MustCompile(`不允许解封,因为IP \{0\} 未被封堵`),
			},
		},
	})
}

const testUnblockIp_basic = `
resource "huaweicloud_aad_unblock_ip" "test" {
  ip            = "*.*.*.*"
  blocking_time = 1688105685078
}`

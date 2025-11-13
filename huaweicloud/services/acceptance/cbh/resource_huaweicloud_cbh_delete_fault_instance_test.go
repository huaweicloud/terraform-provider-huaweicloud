package cbh

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDeleteFaultInstance_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCbhServerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testDeleteFaultInstance_basic,
				ExpectError: regexp.MustCompile(`租户无权限`),
			},
		},
	})
}

// The value of instance_id is a fixed value, not a variable.
const testDeleteFaultInstance_basic = `
resource "huaweicloud_cbh_delete_fault_instance" "test" {
  instance_id = "183743"
}
`

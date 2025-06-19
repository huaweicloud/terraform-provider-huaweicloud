package cbr

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Because there is a lack of scenarios for testing the API, the test case only tests one failure error.
func TestAccResourceCheckpointSync_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testResourceCheckpointSync_basic,
				ExpectError: regexp.MustCompile(`error creating CBR checkpoint sync`),
			},
		},
	})
}

const testResourceCheckpointSync_basic = `
resource "huaweicloud_cbr_checkpoint_sync" "test" {
  vault_id     = "not-exist-vault-id"
  auto_trigger = true
}
`

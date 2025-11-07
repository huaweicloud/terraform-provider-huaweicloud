package dew

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsDedicatedKeystoreAction_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testKmsDedicatedKeystoreAction_basic,
				ExpectError: regexp.MustCompile(`error disable DEW dedicated keystore in creation operation`),
			},
		},
	})
}

// The keystore_id used is dummy data for testing.
const testKmsDedicatedKeystoreAction_basic = `
resource "huaweicloud_kms_dedicated_keystore_action" "test" {
  keystore_id = "b8d0593a-69bc-40f3-b14d-a8b5839fc426"
  action      = "disable"
}`

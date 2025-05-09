package lts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccLogCollectionSwitch_basic(t *testing.T) {
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
				Config: resourceLogCollectionSwitch_basic("disable"),
			},
			{
				Config: resourceLogCollectionSwitch_basic("enable"),
			},
		},
	})
}

func resourceLogCollectionSwitch_basic(action string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_log_collection_switch" "test" {
  action           = "%[1]s"
  enable_force_new = "true"
}
`, action)
}

package eg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEventSubscriptionBatchAction_basic(t *testing.T) {
	var (
		rcName = "huaweicloud_eg_event_subscription_batch_action.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgEventSubscriptionIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testEventSubscriptionBatchAction_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rcName, "operation", "ENABLE"),
					resource.TestMatchResourceAttr(rcName, "subscription_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testEventSubscriptionBatchAction_basic() string {
	return fmt.Sprintf(`
locals {
  event_subscription_ids_str = "%[1]s"
  event_subscription_ids = split("," ,local.event_subscription_ids_str)
}

resource "huaweicloud_eg_event_subscription_batch_action" "test" {
  subscription_ids = local.event_subscription_ids
  operation        = "ENABLE"
}
`, acceptance.HW_EG_EVENT_SUBSCRIPTION_IDS)
}

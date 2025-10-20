package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeScheduledEventAccept_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSScheduledEventId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeScheduledEventAccept_basic(),
			},
		},
	})
}

func testAccComputeScheduledEventAccept_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_scheduled_event_accept" "test" {
  event_id = "%s"
}
`, acceptance.HW_ECS_SCHEDULED_EVENT_ID)
}

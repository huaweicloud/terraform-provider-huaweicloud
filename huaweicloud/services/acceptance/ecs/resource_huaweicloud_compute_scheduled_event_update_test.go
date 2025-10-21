package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeScheduledEventUpdate_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSScheduledEventId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeScheduledEventUpdate_basic(),
			},
		},
	})
}

func testAccComputeScheduledEventUpdate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_scheduled_event_update" "test" {
  event_id   = "%s"
  not_before = "2025-07-09T10:40:00Z"
}
`, acceptance.HW_ECS_SCHEDULED_EVENT_ID)
}

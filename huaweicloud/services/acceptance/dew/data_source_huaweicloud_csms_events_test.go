package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCsmsEvents_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_csms_events.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCsmsEvents_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "events.0.name"),
					resource.TestCheckResourceAttrSet(rName, "events.0.event_id"),
					resource.TestCheckResourceAttrSet(rName, "events.0.event_types.#"),
					resource.TestCheckResourceAttrSet(rName, "events.0.status"),
					resource.TestCheckResourceAttrSet(rName, "events.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "events.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "events.0.notification.#"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("eventId_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCsmsEvents_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_csms_events" "test" {
  depends_on = [huaweicloud_csms_event.test]
}

data "huaweicloud_csms_events" "name_filter" {
  name = huaweicloud_csms_event.test.name
}

locals {
  name = huaweicloud_csms_event.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_csms_events.name_filter.events) > 0 && alltrue(
    [for v in data.huaweicloud_csms_events.name_filter.events[*].name : v == local.name]
  )
}

data "huaweicloud_csms_events" "eventId_filter" {
  event_id = huaweicloud_csms_event.test.event_id
}

locals {
  event_id = huaweicloud_csms_event.test.event_id
}

output "eventId_filter_is_useful" {
  value = length(data.huaweicloud_csms_events.eventId_filter.events) > 0 && alltrue(
    [for v in data.huaweicloud_csms_events.eventId_filter.events[*].event_id : v == 
  local.event_id]
  )  
}

data "huaweicloud_csms_events" "status_filter" {
  status = huaweicloud_csms_event.test.status
}

locals {
  status = huaweicloud_csms_event.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_csms_events.status_filter.events) > 0 && alltrue(
    [for v in data.huaweicloud_csms_events.status_filter.events[*].status : v == local.status]
  )
}
`, testCsmsEvent_basic(name))
}

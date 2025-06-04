package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFunctionTriggerStatusAction_basic(t *testing.T) {
	var (
		obj interface{}

		rcTrigger = "huaweicloud_fgs_function_trigger.test"
		rc        = acceptance.InitResourceCheck(rcTrigger, &obj, getFunctionTriggerFunc)

		name = acceptance.RandomAccResourceName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTriggerStatusAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rcTrigger, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccFunctionTriggerStatusAction_basic_step2(name),
			},
			{
				Config: testAccFunctionTriggerStatusAction_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rcTrigger, "status", "DISABLED"),
				),
			},
		},
	})
}

func testAccFunctionTriggerStatusAction_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  description = "Created by terraform test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python3.6"
  code_type   = "inline"
  func_code   = "def handler (event, context):\n    return {'statusCode': 200, 'body': 'Hello World'}"
}

// Timing trigger (with rate schedule type)
resource "huaweicloud_fgs_function_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "ACTIVE"
  event_data   = jsonencode({
    "name": "%[1]s_rate",
    "schedule_type": "Rate",
    "sync_execution": false,
    "user_event": "Created by terraform script",
    "schedule": "3m"
  })

  lifecycle {
    ignore_changes = [
      status,
    ]
  }
}
`, name)
}

func testAccFunctionTriggerStatusAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function_trigger_status_action" "test" {
  function_urn      = huaweicloud_fgs_function.test.urn
  trigger_type_code = "TIMER"
  trigger_id        = huaweicloud_fgs_function_trigger.test.id
  trigger_status    = "DISABLED"
  event_data        = huaweicloud_fgs_function_trigger.test.event_data
}
`, testAccFunctionTriggerStatusAction_basic_step1(name))
}

// Do refresh to check the status is disabled
func testAccFunctionTriggerStatusAction_basic_step3(name string) string {
	return testAccFunctionTriggerStatusAction_basic_step2(name)
}

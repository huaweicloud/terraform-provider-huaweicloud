package fgs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTriggers_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_fgs_triggers.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byType   = "data.huaweicloud_fgs_triggers.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTriggers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "triggers.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "triggers.0.trigger_id"),
					resource.TestCheckResourceAttrSet(all, "triggers.0.trigger_type_code"),
					resource.TestCheckResourceAttrSet(all, "triggers.0.trigger_status"),
					resource.TestCheckResourceAttrSet(all, "triggers.0.func_urn"),
					resource.TestCheckResourceAttrSet(all, "triggers.0.event_data"),
					resource.TestCheckResourceAttrSet(all, "triggers.0.created_time"),
					resource.TestCheckResourceAttrSet(all, "triggers.0.last_updated_time"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataTriggers_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python3.6"
  code_type   = "inline"
  func_code   = "def handler (event, context):\n    return {'statusCode': 200, 'body': 'Hello World'}"
}

resource "huaweicloud_fgs_function_trigger" "test" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  event_data   = jsonencode({
    name           = "%[1]s_rate",
    schedule_type  = "Rate",
    sync_execution = false,
    user_event     = "Created by terraform script",
    schedule       = "3m"
  })
}
`, acceptance.RandomAccResourceName())
}

func testAccDataTriggers_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_fgs_triggers" "test" {
  depends_on = [huaweicloud_fgs_function_trigger.test]
}

locals {
  trigger_type = huaweicloud_fgs_function_trigger.test.type
}

data "huaweicloud_fgs_triggers" "filter_by_type" {
  trigger_type = local.trigger_type

  depends_on = [huaweicloud_fgs_function_trigger.test]
}

locals {
  type_filter_results = [for v in data.huaweicloud_fgs_triggers.filter_by_type.triggers[*].trigger_type_code : v == local.trigger_type]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_results) > 0 && alltrue(local.type_filter_results)
}
`, testAccDataTriggers_base())
}

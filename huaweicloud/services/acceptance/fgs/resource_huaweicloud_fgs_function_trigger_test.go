package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/trigger"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getFunctionTriggerFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("fgs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	return fgs.GetTriggerById(client, state.Primary.Attributes["function_urn"], state.Primary.Attributes["type"],
		state.Primary.ID)
}

func TestAccFunctionTrigger_basic(t *testing.T) {
	var (
		relatedFunc      function.Function
		timeTrigger      trigger.Trigger
		randName         = acceptance.RandomAccResourceName()
		resNameFunc      = "huaweicloud_fgs_function.test"
		resNameTimerRate = "huaweicloud_fgs_function_trigger.timer_rate"
		resNameTimerCron = "huaweicloud_fgs_function_trigger.timer_cron"

		rcFunc      = acceptance.InitResourceCheck(resNameFunc, &relatedFunc, getFunction)
		rcTimerRate = acceptance.InitResourceCheck(resNameTimerRate, &timeTrigger, getFunctionTriggerFunc)
		rcTimerCron = acceptance.InitResourceCheck(resNameTimerCron, &timeTrigger, getFunctionTriggerFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcFunc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTimingTrigger_basic_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					// Timing trigger (with rate schedule type)
					rcTimerRate.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resNameTimerRate, "function_urn",
						"${huaweicloud_fgs_function.test.urn}"),
					resource.TestCheckResourceAttr(resNameTimerRate, "type", "TIMER"),
					resource.TestCheckResourceAttr(resNameTimerRate, "status", "ACTIVE"),
					// Timing trigger (with cron schedule type)
					rcTimerCron.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerCron, "type", "TIMER"),
					resource.TestCheckResourceAttr(resNameTimerCron, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccFunctionTimingTrigger_basic_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					// Timing trigger (with rate schedule type)
					rcTimerRate.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerRate, "status", "DISABLED"),
					// Timing trigger (with cron schedule type)
					rcTimerCron.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerCron, "status", "DISABLED"),
				),
			},
			{
				ResourceName:      resNameTimerRate,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionTriggerImportStateFunc(resNameTimerRate),
			},
			{
				ResourceName:      resNameTimerCron,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionTriggerImportStateFunc(resNameTimerCron),
			},
		},
	})
}

func testAccFunctionTriggerImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var functionUrn, triggerType, triggerId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of function trigger is not found in the tfstate", rsName)
		}
		functionUrn = rs.Primary.Attributes["function_urn"]
		triggerType = rs.Primary.Attributes["type"]
		triggerId = rs.Primary.ID
		if functionUrn == "" || triggerType == "" || triggerId == "" {
			return "", fmt.Errorf("the function trigger is not exist or related function URN is missing")
		}
		return fmt.Sprintf("%s/%s/%s", functionUrn, triggerType, triggerId), nil
	}
}

func testAccFunctionTrigger_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgW1wcyhldybiBvdXRwdXQ="
}`, name)
}

// Test triggers with a limited number (except for Kafka triggers, when released, the elastic network card will be
// locked for one hour and the subnet cannot be deleted).
// The current quantity constraint rules for function triggers are as follows:
//   - The total number of DDS, GAUSSMONGO, DIS, LTS, Kafka and TIMER triggers that can be created under one function
//     version is up to 10.
//   - The maximum number of CTS triggers that can be created under one project is 10.
//   - There is no limit to the number of other triggers.
func testAccFunctionTimingTrigger_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

// Timing trigger (with rate schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_rate" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "ACTIVE"
  event_data   = jsonencode({
    "name": "%[2]s_rate",
    "schedule_type": "Rate",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "3m"
  })
}

// Timing trigger (with cron schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_cron" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "ACTIVE"
  event_data   = jsonencode({
    "name": "%[2]s_cron",
    "schedule_type": "Cron",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "@every 1h30m"
  })
}
`, testAccFunctionTrigger_base(name), name)
}

func testAccFunctionTimingTrigger_basic_step2(name string) string {
	return fmt.Sprintf(`
%s

// Timing trigger (with rate schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_rate" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"
  event_data   = jsonencode({
    "name": "%[2]s_rate",
    "schedule_type": "Rate",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "3m"
  })
}

// Timing trigger (with cron schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_cron" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"
  event_data   = jsonencode({
    "name": "%[2]s_cron",
    "schedule_type": "Cron",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "@every 1h30m"
  })
}
`, testAccFunctionTrigger_base(name), name)
}

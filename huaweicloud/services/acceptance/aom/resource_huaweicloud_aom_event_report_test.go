package aom

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEventReport_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testEventReport_expectError,
				ExpectError: regexp.MustCompile(`metadata is invalid`),
			},
			{
				Config: testEventReport_basic_step1(),
			},
			{
				Config: testEventReport_basic_step2(),
			},
			{
				Config: testEventReport_basic_step3(),
			},
			{
				Config: testEventReport_basic_step4(),
			},
			{
				Config: testEventReport_basic_step5(),
			},
		},
	})
}

const testEventReport_expectError string = `
resource "huaweicloud_aom_event_report" "test_with_error" {
  events {
    timeout = 432000000

    metadata = {
      event_name     = "normal_event"
      event_severity = "Major"
      event_type     = "alarm"
    }
  }
}
`

func testEventReport_basic_step1() string {
	timestamp := time.Now().UnixMilli()
	randomUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_aom_event_report" "test" {
  events {
    starts_at = %[1]d
    timeout   = 432000000

    metadata = {
      event_name        = "normal_event"
      event_severity    = "Major"
      event_type        = "alarm"
      resource_provider = "ecs"
      resource_type     = "vm"
      resource_id       = "%[2]s"
    }
  }
}
`, timestamp, randomUUID)
}

func testEventReport_basic_step2() string {
	timestamp := time.Now().UnixMilli()
	randomUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_aom_event_report" "test_with_annotations" {
  events {
    starts_at = %[1]d
    timeout   = 432000000

    metadata = {
      event_name        = "event_with_annotations"
      event_severity    = "Critical"
      event_type        = "alarm"
      resource_provider = "ecs"
      resource_type     = "vm"
      resource_id       = "%[2]s"
    }

    annotations = jsonencode({
      message = "This is a test alarm message"
    })
  }
}
`, timestamp, randomUUID)
}

func testEventReport_basic_step3() string {
	// 1 minute later
	timestamp := time.Now().UnixMilli()
	clearTimestamp := timestamp + 60000
	randomUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_aom_event_report" "test_with_clear_action" {
  action = "clear"

  events {
    ends_at = %[1]d

    metadata = {
      event_name        = "terraform_test_event_clear"
      event_severity    = "Major"
      event_type        = "alarm"
      resource_provider = "ecs"
      resource_type     = "vm"
      resource_id       = "%[2]s"
    }
  }
}
`, clearTimestamp, randomUUID)
}

func testEventReport_basic_step4() string {
	timestamp := time.Now().UnixMilli()
	randomUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_aom_event_report" "test_with_default_ep_id" {
  enterprise_project_id = "0"

  events {
    starts_at = %[1]d

    metadata = {
      event_name        = "event_with_ep_id"
      event_severity    = "Minor"
      event_type        = "event"
      resource_provider = "ecs"
      resource_type     = "vm"
      resource_id       = "%[2]s"
    }
  }
}
`, timestamp, randomUUID)
}

func testEventReport_basic_step5() string {
	timestamp := time.Now().UnixMilli()
	randomUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
resource "huaweicloud_aom_event_report" "test_with_ep_id" {
  enterprise_project_id = "%[1]s"

  events {
    starts_at = %[2]d

    metadata = {
      event_name        = "event_with_ep_id"
      event_severity    = "Minor"
      event_type        = "event"
      resource_provider = "ecs"
      resource_type     = "vm"
      resource_id       = "%[3]s"
    }
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID, timestamp, randomUUID)
}

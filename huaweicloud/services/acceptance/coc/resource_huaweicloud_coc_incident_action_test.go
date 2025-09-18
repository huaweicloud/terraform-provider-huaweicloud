package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceIncidentAction_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testIncidentAction_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testIncidentAction_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "time_sleep" "wait_10_seconds" {
  depends_on      = [huaweicloud_coc_incident.test]
  create_duration = "10s"
}

data "huaweicloud_coc_incident_tasks" "test" {
  depends_on  = [time_sleep.wait_10_seconds]
  incident_id = huaweicloud_coc_incident.test.id
}

locals {
  task_id = [for v in data.huaweicloud_coc_incident_tasks.test.data[0].operations[*].task_id : v if v != ""][0]
}

output "task_id" {
  value = local.task_id
}

resource "huaweicloud_coc_incident_action" "test" {
  incident_id = huaweicloud_coc_incident.test.id
  task_id     = local.task_id
  action      = "rejected"

  params = {
    virtual_confirm_comment = "test comment",
  }

  lifecycle {
    ignore_changes = [
      task_id
    ]
  }
}
`, testIncident_basic(rName))
}

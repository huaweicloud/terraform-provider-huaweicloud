package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAlarmAction_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPreCheckCocAlarmID(t)
			acceptance.TestAccPreCheckCocAgentSn(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAlarmAction_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAlarmAction_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  resource_id = "%[2]s"
  region      = "%[3]s"
}

data "huaweicloud_compute_instances" "test" {
  instance_id = local.resource_id
}

locals {
  name              = data.huaweicloud_compute_instances.test.instances[0].name
  access_ip_v4      = data.huaweicloud_compute_instances.test.instances[0].network[0].fixed_ip_v4
  availability_zone = data.huaweicloud_compute_instances.test.instances[0].availability_zone
  port              = data.huaweicloud_compute_instances.test.instances[0].network[0].port
  mac               = data.huaweicloud_compute_instances.test.instances[0].network[0].mac
}

resource "huaweicloud_coc_alarm_action" "test" {
  alarm_id                              = "%[4]s"
  task_type                             = "SCRIPT"
  associated_task_id                    = huaweicloud_coc_script.test.id
  associated_task_type                  = "CUSTOMIZATION"
  associated_task_name                  = huaweicloud_coc_script.test.name
  associated_task_enterprise_project_id = "0"
  runbook_instance_mode                 = "SAME"
  input_param = {
    timeout       = 300
    execute_user  = "root"
    success_rate  = 100
    project_id    = "%[5]s"
    script_params = jsonencode([{
      "paramName": "name",
      "paramValue": "world",
      "paramOrder": 1
    }])
  }
  target_instances {
    target_selection = "MANUAL"
    order_no         = 0
    batch_strategy   = "NONE"
    target_instances = jsonencode({
      "batches": [
        {
          "batchIndex": 1,
          "rotationStrategy": "CONTINUE",
          "targetInstances": [
            {
              "resourceId": local.resource_id,
              "regionId": local.region,
              "provider": "ECS",
              "type": "CLOUDSERVERS",
              "agentSn": "%[6]s",
              "agentStatus": "ONLINE",
              "nodeId": "",
              "enterpriseProjectId": "0",
              "properties": {
                "fixedIp": local.access_ip_v4,
                "regionId": local.region,
                "zoneId": local.availability_zone,
                "projectId": "%[5]s"
              }
            }
          ],
          "cmdbInstances": [
            {
              "resourceId": local.resource_id,
              "name": local.name,
              "projectId": "%[5]s",
              "regionId": local.region,
              "agentId": "%[6]s",
              "agentState": "ONLINE",
              "provider": "ecs",
              "enterpriseProjectId": "0",
              "type": "cloudservers",
              "properties": {
                "addresses": [
                  {
                    "OsExtIpsType": "fixed",
                    "OsExtIpsPortId": local.port,
                    "addr": local.access_ip_v4,
                    "version": 4,
                    "OsExtIpsMacAddr": local.mac,
                    "primary": true
                  }
                ],
                "metadata": {
                  "osType": "Linux"
                },
                "OsExtAz": local.availability_zone,
                "status": "ACTIVE"
              }
            }
          ]
        }
      ],
      "policy": "none",
      "all_rotation": "FIRST_PAUSE"
    })
  }
}
`, tesScript_basic(rName), acceptance.HW_COC_INSTANCE_ID, acceptance.HW_REGION_NAME, acceptance.HW_COC_ALARM_ID,
		acceptance.HW_PROJECT_ID, acceptance.HW_COC_AGENT_SN)
}

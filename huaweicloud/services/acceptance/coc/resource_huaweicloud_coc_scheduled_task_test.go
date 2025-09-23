package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getScheduledTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetScheduledTask(client, state.Primary.ID)
}

func TestAccResourceScheduledTask_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	rNameUpdate := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_coc_scheduled_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getScheduledTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testScheduledTask_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "LOW"),
					resource.TestCheckResourceAttrSet(resourceName, "approve_status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_user_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"ticket_infos", "associated_task_enterprise_project_id",
					"target_instances.0.batch_strategy"},
			},
			{
				Config: testScheduledTask_basic_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "risk_level", "HIGH"),
					resource.TestCheckResourceAttrSet(resourceName, "approve_status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_user_name"),
				),
			},
		},
	})
}

func testScheduledTask_basic(name string) string {
	currentTime := time.Now()
	tenMinutesLater := currentTime.Add(10*time.Minute).Unix() * 1e3
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "time_sleep" "wait_10_seconds" {
  depends_on = [huaweicloud_compute_instance.test]

  create_duration = "30s"
}

resource "huaweicloud_coc_scheduled_task" "test" {
  depends_on = [time_sleep.wait_10_seconds]
  name       = "%[3]s"
  version_no = "1.0.0"
  trigger_time {
    time_zone             = "Asia/Shanghai"
    policy                = "ONCE"
    single_scheduled_time = %[4]v
  }
  task_type            = "SCRIPT"
  associated_task_id   = huaweicloud_coc_script.test.id
  associated_task_type = "CUSTOMIZATION"
  associated_task_name = huaweicloud_coc_script.test.name
  risk_level           = "LOW"

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
    target_instances = jsonencode({
      "batches": [
        {
          "batchIndex": 1,
          "rotationStrategy": "CONTINUE",
          "targetInstances": [
            {
              "resourceId": huaweicloud_compute_instance.test.id,
              "regionId": huaweicloud_compute_instance.test.region,
              "provider": "ECS",
              "type": "CLOUDSERVERS",
              "agentStatus": "ONLINE",
              "nodeId": "",
              "enterpriseProjectId": "0",
              "properties": {
                "hostName": huaweicloud_compute_instance.test.hostname,
                "fixedIp": huaweicloud_compute_instance.test.access_ip_v4,
                "regionId": huaweicloud_compute_instance.test.region,
                "zoneId": huaweicloud_compute_instance.test.availability_zone,
                "projectId": "%[5]s"
              }
            }
          ]
        }
      ],
      "policy": "none",
      "all_rotation": "ALL_CONTINUE"
    })
  }
  enable_approve              = false
  enable_message_notification = false
  enterprise_project_id       = "0"
  agency_name                 = "ServiceAgencyForCOC"
  runbook_instance_mode       = "SAME"
  enabled                     = "true"
}
`, testAccComputeInstance_basic(name), tesScript_basic(name), name, tenMinutesLater, acceptance.HW_PROJECT_ID)
}

func testScheduledTask_basic_update(originalName string, name string) string {
	currentTime := time.Now()
	endTime := currentTime.Add(168*time.Hour).Unix() * 1e3
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "time_sleep" "wait_10_seconds" {
  depends_on = [huaweicloud_compute_instance.test]

  create_duration = "30s"
}

resource "huaweicloud_coc_scheduled_task" "test" {
  depends_on = [time_sleep.wait_10_seconds]
  name       = "%[3]s"
  version_no = "1.0.0"
  trigger_time {
    time_zone               = "Asia/Shanghai"
    policy                  = "PERIODIC"
    periodic_scheduled_time = "09:10:00"
    period                  = "2"
    scheduled_close_time    = %[4]v
  }
  task_type            = "SCRIPT"
  associated_task_id   = huaweicloud_coc_script.test.id
  associated_task_type = "CUSTOMIZATION"
  associated_task_name = huaweicloud_coc_script.test.name
  risk_level           = "HIGH"

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
    target_instances = jsonencode({
      "batches": [
        {
          "batchIndex": 1,
          "rotationStrategy": "CONTINUE",
          "targetInstances": [
            {
              "resourceId": huaweicloud_compute_instance.test.id,
              "regionId": huaweicloud_compute_instance.test.region,
              "provider": "ECS",
              "type": "CLOUDSERVERS",
              "agentStatus": "ONLINE",
              "nodeId": "",
              "enterpriseProjectId": "0",
              "properties": {
                "hostName": huaweicloud_compute_instance.test.hostname,
                "fixedIp": huaweicloud_compute_instance.test.access_ip_v4,
                "regionId": huaweicloud_compute_instance.test.region,
                "zoneId": huaweicloud_compute_instance.test.availability_zone,
                "projectId": "%[5]s"
              }
            }
          ]
        }
      ],
      "policy": "automatic",
      "all_rotation": "ALL_CONTINUE"
    })
    batch_strategy = "AUTO_BATCH"
  }
  enable_approve                         = false
  enable_message_notification            = false
  enterprise_project_id                  = "0"
  agency_name                            = "ServiceAgencyForCOC"
  associated_task_enterprise_project_id  = "0"
  runbook_instance_mode                  = "SAME"

  enabled = "false"
  lifecycle {
    ignore_changes = [
      target_instances.0.batch_strategy
    ]
  }
}
`, testAccComputeInstance_basic(originalName), tesScript_basic(originalName), name, endTime, acceptance.HW_PROJECT_ID)
}

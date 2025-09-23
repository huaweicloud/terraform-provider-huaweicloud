package aom

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAlarmRuleV4ResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getHttpUrl := "v4/{project_id}/alarm-rules?name={name}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{name}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(state),
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the rule: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening the rule: %s", err)
	}

	rule := utils.PathSearch("alarm_rules|[0]", getRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func TestAccAlarmRuleV4_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aomv4_alarm_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getAlarmRuleV4ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlarmRuleV4_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "event"),
					resource.TestCheckResourceAttr(resourceName, "enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName, "alarm_notifications.0.notification_type", "direct"),
					resource.TestCheckResourceAttr(resourceName, "alarm_notifications.0.notification_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "alarm_notifications.0.bind_notification_rule_id",
						"huaweicloud_aom_alarm_action_rule.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "event_alarm_spec.0.event_source", "CCE"),
					resource.TestCheckResourceAttr(resourceName, "event_alarm_spec.0.alarm_source", "systemEvent"),
				),
			},
			{
				Config: testAlarmRuleV4_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "event"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", "testUpdate"),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_notifications.0.notification_type", "direct"),
					resource.TestCheckResourceAttr(resourceName, "alarm_notifications.0.notification_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "alarm_notifications.0.bind_notification_rule_id",
						"huaweicloud_aom_alarm_action_rule.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "event_alarm_spec.0.event_source", "ModelArts"),
					resource.TestCheckResourceAttr(resourceName, "event_alarm_spec.0.alarm_source", "systemEvent"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAlarmRuleV4_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aomv4_alarm_rule" "test" {
  name        = "%[2]s"
  type        = "event"
  description = "test"

  alarm_notifications {
    notification_type         = "direct"
    notification_enable       = true
    bind_notification_rule_id = huaweicloud_aom_alarm_action_rule.test.name
  }

  event_alarm_spec {
    event_source = "CCE"
    alarm_source = "systemEvent"

    monitor_objects = [
      {
        event_name = "扩容节点超时##ScaleUpTimedOut;数据卷扩容失败##VolumeResizeFailed"
      },
    ]

    trigger_conditions {
      trigger_type = "immediately"
      event_name   = "扩容节点超时##ScaleUpTimedOut"

      thresholds = {
        "Critical" = 2
      }
    }

    trigger_conditions {
      trigger_type       = "accumulative"
      event_name         = "数据卷扩容失败##VolumeResizeFailed"
      aggregation_window = 300
      frequency          = "600"
      operator           = ">="

      thresholds = {
        "Info" = 5
      }
    }
  }
}
`, testAlarmActionRule_basic(name), name)
}

func testAlarmRuleV4_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aomv4_alarm_rule" "test" {
  name                  = "%[2]s"
  type                  = "event"
  enterprise_project_id = "0"
  description           = "testUpdate"
  enable                = true

  alarm_notifications {
    notification_type         = "direct"
    notification_enable       = true
    bind_notification_rule_id = huaweicloud_aom_alarm_action_rule.test.name
  }

  event_alarm_spec {
    event_source = "ModelArts"
    alarm_source = "systemEvent"

    monitor_objects = [
      {
        event_name = "Scheduled;PullingImage"
      },
    ]

    trigger_conditions {
      trigger_type = "immediately"
      event_name   = "Scheduled"

      thresholds = {
        "Major" = 1
      }
    }

    trigger_conditions {
      trigger_type       = "accumulative"
      event_name         = "PullingImage"
      aggregation_window = 1200
      frequency          = "900"
      operator           = ">"

      thresholds = {
        "Critical" = 7
      }
    }
  }
}
`, testAlarmActionRule_basic(name), name)
}

func TestAccAlarmRuleV4_metric(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aomv4_alarm_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getAlarmRuleV4ResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlarmRuleV4_metric(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "metric"),
					resource.TestCheckResourceAttr(resourceName, "enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "alarm_notifications.0.notification_type", "direct"),
					resource.TestCheckResourceAttr(resourceName, "metric_alarm_spec.0.monitor_type", "all_metric"),
					resource.TestCheckResourceAttr(resourceName, "metric_alarm_spec.0.recovery_conditions.0.recovery_timeframe", "1"),
				),
			},
			{
				Config: testAlarmRuleV4_metric_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "metric"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_notifications.0.notification_type", "direct"),
					resource.TestCheckResourceAttr(resourceName, "metric_alarm_spec.0.monitor_type", "all_metric"),
					resource.TestCheckResourceAttr(resourceName, "metric_alarm_spec.0.recovery_conditions.0.recovery_timeframe", "2"),
					resource.TestCheckResourceAttr(resourceName, "metric_alarm_spec.0.no_data_conditions.0.notify_no_data", "true"),
					resource.TestCheckResourceAttr(resourceName, "metric_alarm_spec.0.no_data_conditions.0.no_data_timeframe", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric_alarm_spec.0.no_data_conditions.0.no_data_alert_state", "no_data"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAlarmRuleV4_metric(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aomv4_alarm_rule" "test" {
  name             = "%[1]s"
  type             = "metric"
  prom_instance_id = "0"

  alarm_notifications {
    notification_type = "direct"
  }

  metric_alarm_spec {
    monitor_type = "all_metric"

    recovery_conditions {
      recovery_timeframe = 1
    }

    trigger_conditions {
      metric_query_mode       = "PROM"
      metric_name             = "duration"
      promql                  = "label_replace(avg_over_time(duration{}[59999ms]),\"__name__\",\"duration\",\"\",\"\")"
      trigger_times           = "3"
      trigger_type            = "FIXED_RATE"
      aggregation_window      = "1m"
      trigger_interval        = "30s"
      aggregation_type        = "average"
      operator                = ">"
      metric_statistic_method = "single"

      thresholds = {
        "Critical" = "1"
      }
    }
  }
}
`, name)
}

func testAlarmRuleV4_metric_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aomv4_alarm_rule" "test" {
  name                  = "%[2]s"
  type                  = "metric"
  enable                = true
  prom_instance_id      = "0"
  enterprise_project_id = "0"

  alarm_notifications {
    notification_type         = "direct"
    notify_resolved           = true
    notify_triggered          = true
    notification_enable       = true
    bind_notification_rule_id = huaweicloud_aom_alarm_action_rule.test.name
  }

  metric_alarm_spec {
    monitor_type = "all_metric"

    recovery_conditions {
      recovery_timeframe = 2
    }

    trigger_conditions {
      metric_query_mode       = "PROM"
      metric_name             = "ALERTS"
      promql                  = "label_replace(avg_over_time(ALERTS{}[59999ms]),\"__name__\",\"ALERTS\",\"\",\"\")"
      trigger_times           = "4"
      trigger_type            = "FIXED_RATE"
      aggregation_window      = "5m"
      trigger_interval        = "15s"
      aggregation_type        = "sum"
      operator                = "<="
      metric_statistic_method = "single"

      thresholds = {
        "Info" = "2"
      }
    }

    alarm_tags {
      custom_tags = ["key=value"]
    }

    no_data_conditions {
      notify_no_data      = true
      no_data_timeframe   = 1
      no_data_alert_state = "no_data"
    }
  }
}
`, testAlarmActionRule_basic(name), name)
}

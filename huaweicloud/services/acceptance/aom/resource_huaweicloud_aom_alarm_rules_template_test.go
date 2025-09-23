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

func getAlarmRulesTemplateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	getHttpUrl := "v4/{project_id}/alarm-rules-template?id={id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeaders(state),
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the alarm rules template: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening the alarm rules template: %s", err)
	}

	template := utils.PathSearch("[]|[0]", getRespBody, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return template, nil
}

func TestAccAlarmRulesTemplate_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_aom_alarm_rules_template.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getAlarmRulesTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAlarmRulesTemplate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "statics"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttr(resourceName,
						"alarm_template_spec_list.0.related_cloud_service", "CCEFromProm"),
					resource.TestCheckResourceAttr(resourceName,
						"alarm_template_spec_list.0.alarm_notification.0.notification_type", "direct"),
					resource.TestCheckResourceAttr(resourceName,
						"alarm_template_spec_list.0.alarm_notification.0.notification_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName,
						"alarm_template_spec_list.0.alarm_notification.0.bind_notification_rule_id",
						"huaweicloud_aom_alarm_action_rule.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "alarm_template_spec_list.0.related_cloud_service", "CCEFromProm"),
				),
			},
			{
				Config: testAlarmRulesTemplate_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "statics"),
					resource.TestCheckResourceAttr(resourceName, "description", "testUpdate"),
					resource.TestCheckResourceAttr(resourceName, "alarm_template_spec_list.0.related_cloud_service", "DRS"),
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

func testAlarmRulesTemplate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_alarm_rules_template" "test" {
  name        = "%[2]s"
  type        = "statics"
  description = "test"

  alarm_template_spec_list {
    related_cloud_service = "CCEFromProm"

    alarm_notification {
      notification_type         = "direct"
      notification_enable       = true
      bind_notification_rule_id = huaweicloud_aom_alarm_action_rule.test.name
    }

    alarm_template_spec_items {
      alarm_rule_name = "cce_event"
      alarm_rule_type = "event"

      event_alarm_spec {
        event_source = "CCE"

        monitor_objects = [
          {
            event_name = "扩容节点超时##ScaleUpTimedOut;数据卷扩容失败##VolumeResizeFailed"
          },
        ]

        monitor_object_templates = ["clusterId"]

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

    alarm_template_spec_items {
      alarm_rule_name = "cce_metric"
      alarm_rule_type = "metric"

      metric_alarm_spec {
        monitor_type = "promql"

        recovery_conditions {
          recovery_timeframe = 1
        }

        trigger_conditions {
          metric_query_mode  = "NATIVE_PROM"
          promql             = "increase(kube_pod_container_status_restarts_total[5m]) > 3"
          trigger_times      = "3"
          trigger_type       = "FIXED_RATE"
          aggregation_window = "1m"
          trigger_interval   = "30s"
          aggregation_type   = "average"
          operator           = ">"
          promql_for         = "1m"

          thresholds = {
            "Critical" = "1"
          }
        }
      }
    }
  }

  templating {
    list {
      name  = "key"
      type  = "constant"
      query = "value"
    }
  }
}
`, testAlarmActionRule_basic(name), name)
}

func testAlarmRulesTemplate_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_aom_alarm_rules_template" "test" {
  name        = "%[2]s"
  type        = "statics"
  description = "testUpdate"

  alarm_template_spec_list {
    related_cloud_service = "DRS"

    alarm_notification {
      notification_type = "direct"
    }

    alarm_template_spec_items {
      alarm_rule_name = "drs"
      alarm_rule_type = "metric"

      metric_alarm_spec {
        monitor_type = "resource"

        recovery_conditions {
          recovery_timeframe = 1
        }

        trigger_conditions {
          metric_query_mode       = "PROM"
          metric_name             = "huaweicloud_sys_drs_cpu_util"
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
  }
}
`, testAlarmActionRule_basic(name), name)
}

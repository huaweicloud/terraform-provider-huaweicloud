package aom

import (
	"encoding/json"
	"fmt"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAlarmPolicyResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, _ := httpclient_go.NewHttpClientGo(config)
	client.WithMethod(httpclient_go.MethodGet).WithUrlWithoutEndpoint(config, "aom", config.Region,
		"v4/"+state.Primary.Attributes["project_id"]+"/alarm-rules")

	response, err := client.Do()
	body, _ := client.CheckDeletedDiag(nil, err, response, "")
	if body == nil || string(body) == "[]" {
		return nil, fmt.Errorf("error getting HuaweiCloud Resource")
	}

	rlt := &[]entity.AddAlarmRuleParams{}
	err = json.Unmarshal(body, rlt)
	for _, params := range *rlt {
		if params.AlarmRuleName == state.Primary.Attributes["alarm_rule_name"] {
			return rlt, nil
		}
	}

	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("unable to find the persistent volume claim (%s)", state.Primary.ID)
	}
	return nil, fmt.Errorf("error getting HuaweiCloud Resource")
}

func TestAccAOMAlarmPolicy_basic(t *testing.T) {
	var ar []entity.AddAlarmRuleParams
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_aom_alarm_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckInternal(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAOMAlarmPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_rule_name", rName),
					resource.TestCheckResourceAttr(resourceName, "action_id", "add-alarm-action"),
					resource.TestCheckResourceAttr(resourceName, "alarm_rule_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_rule_type", "metric"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "project_id", "2a473356cca5487f8373be891bffc1cf"),
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

func testAOMAlarmPolicy_basic(rName string) string {
	return fmt.Sprintf(`
provider "huaweicloud" {
  region     = "cn-north-7"
  access_key = "KKILFTQC8DJYDINQQZXI"
  secret_key = "yM9DG7GaISH9Ob2n6zrq89IvdiZC68keqOad9oVu"
  auth_url = "https://iam.cn-north-7.myhuaweicloud.com"
  endpoints = {
    aom : "aom.cn-north-7.myhuaweicloud.com"
  }

  insecure = true
  domain_id = "40de487942a74a70b4666fa32d11ffa8"
  project_id  = "2a473356cca5487f8373be891bffc1cf"
}
  resource "huaweicloud_aom_alarm_policy" "test" {
         action_id	= "add-alarm-action"
         alarm_rule_description = "d"
         alarm_rule_enable        = true
         alarm_rule_name          = "%s"
         alarm_rule_type          = "metric"
         enterprise_project_id    = "0"
         project_id                 = "2a473356cca5487f8373be891bffc1cf"

         alarm_notifications {
	inhibit_enable          = false
	notification_enable   = false
	notification_type      = "direct"
	notify_resolved         = true
	route_group_enable = false
         }

         metric_alarm_spec {
              monitor_objects    = [
	 {
	     "lbaas_instance_id" = "0c9535c3-35b0-4216-8af0-19a748c423c8"
	 },
         ]
         monitor_type        = "promql"
         recovery_conditions = {
	"recovery_timeframe" = 1
         }

         no_data_conditions {
	notify_no_data = false
         }

         trigger_conditions {
	aggregation_type  = "average"
	aggregation_window = "1m"
	metric_labels          = [
	    "__name__",
	    "aom_monitor_level",
	    "lbaas_instance_id",
	    "lbaas_listener_id",
	    "listener_name",
	    "name",
	    "namespace",
	    "port",
	    "workspace_id",
	]
	metric_name      = "huaweicloud_sys_elb_m1_cps"
	metric_query_mode  =  "PROM"
	operator	           = ">="
	promql                = "huaweicloud_sys_elb_m1_cps{lbaas_instance_id=\"0c9535c3-35b0-4216-8af0-19a748c423c8\"}"
	thresholds           = {
	     "Critical" = "0"
	}
	trigger_interval   =  "1m"
	trigger_times    =  1
	trigger_type          = "FIXED_RATE"
             }
         }
     }
`, rName)
}

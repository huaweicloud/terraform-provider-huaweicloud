package lts

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

func getNotificationTemplateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getNotificationTemplate: Query the LTS notification template.
	var (
		getNotificationTemplateHttpUrl = "v2/{project_id}/{domain_id}/lts/events/notification/template/{id}"
		getNotificationTemplateProduct = "lts"
	)
	getNotificationTemplateClient, err := cfg.NewServiceClient(getNotificationTemplateProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	getNotificationTemplatePath := getNotificationTemplateClient.Endpoint + getNotificationTemplateHttpUrl
	getNotificationTemplatePath = strings.ReplaceAll(getNotificationTemplatePath, "{project_id}", getNotificationTemplateClient.ProjectID)
	getNotificationTemplatePath = strings.ReplaceAll(getNotificationTemplatePath, "{domain_id}", cfg.DomainID)
	getNotificationTemplatePath = strings.ReplaceAll(getNotificationTemplatePath, "{id}", state.Primary.ID)

	getNotificationTemplateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getNotificationTemplateResp, err := getNotificationTemplateClient.Request("GET", getNotificationTemplatePath, &getNotificationTemplateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving notification template: %s", err)
	}

	getNotificationTemplateRespBody, err := utils.FlattenResponse(getNotificationTemplateResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving notification template: %s", err)
	}

	return getNotificationTemplateRespBody, nil
}

func TestAccNotificationTemplate_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_notification_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getNotificationTemplateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testNotificationTemplate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source", "LTS"),
					resource.TestCheckResourceAttr(rName, "locale", "zh-cn"),
					resource.TestCheckResourceAttr(rName, "description", "This is acceptance test"),
					resource.TestCheckResourceAttr(rName, "templates.#", "1"),
					resource.TestCheckResourceAttr(rName, "templates.0.sub_type", "sms"),
				),
			},
			{
				Config: testNotificationTemplate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "source", "LTS"),
					resource.TestCheckResourceAttr(rName, "locale", "zh-cn"),
					resource.TestCheckResourceAttr(rName, "description", "This is acceptance test update"),
					resource.TestCheckResourceAttr(rName, "templates.#", "2"),
					resource.TestCheckResourceAttr(rName, "templates.0.sub_type", "sms"),
					resource.TestCheckResourceAttr(rName, "templates.1.sub_type", "email"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testNotificationTemplate_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_notification_template" "test" {
  name        = "%s"
  source      = "LTS"
  locale      = "zh-cn"
  description = "This is acceptance test"

  templates {
    sub_type = "sms"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }
}
`, name)
}

func testNotificationTemplate_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_notification_template" "test" {
  name        = "%s"
  description = "This is acceptance test update"
  source      = "LTS"
  locale      = "zh-cn"

  templates {
    sub_type = "sms"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }

  templates {
    sub_type = "email"
    content  = <<EOF
Account:$${domain_name};
Alarm Rules:<a href="$event.annotations.alarm_rule_url">$${event_name}</a>;
Alarm Status:$event.annotations.alarm_status;
Severity:<span style="color: red">$${event_severity}</span>;
Occurred:$${starts_at};
Type:Keywords;
Condition Expression:$event.annotations.condition_expression;
Current Value:$event.annotations.current_value;
Frequency:$event.annotations.frequency;
Log Group/Stream Name:$event.annotations.results[0].resource_id;
Query Time:$event.annotations.results[0].time;
Query URL:<a href="$event.annotations.results[0].url">details</a>;
EOF
  }
}
`, name)
}

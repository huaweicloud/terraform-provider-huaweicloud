package waf

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

func getAlarmNotificationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "waf"
		epsID   = state.Primary.Attributes["enterprise_project_id"]
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	return waf.GetAlarmNotificationDetail(client, state.Primary.ID, epsID)
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccAlarmNotification_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_waf_alarm_notification.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAlarmNotificationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAlarmNotification_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttr(rName, "notice_class", "threat_alert_notice"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "sendfreq", "120"),
					resource.TestCheckResourceAttr(rName, "locale", "en-us"),
					resource.TestCheckResourceAttr(rName, "times", "5"),
					resource.TestCheckResourceAttr(rName, "is_all_enterprise_project", "true"),
					resource.TestCheckResourceAttr(rName, "description", "test-alarm-notification"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "threat.#"),
				),
			},
			{
				Config: testAlarmNotification_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttrPair(rName, "topic_urn", "huaweicloud_smn_topic.test", "topic_urn"),
					resource.TestCheckResourceAttr(rName, "notice_class", "threat_alert_notice"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "sendfreq", "60"),
					resource.TestCheckResourceAttr(rName, "locale", "en-us"),
					resource.TestCheckResourceAttr(rName, "times", "6"),
					resource.TestCheckResourceAttr(rName, "is_all_enterprise_project", "false"),
					resource.TestCheckResourceAttr(rName, "description", "test-alarm-notification update"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrSet(rName, "threat.#"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAlarmNotificationImportState(rName),
			},
		},
	})
}

func testAlarmNotification_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_waf_alarm_notification" "test" {
  name                      = "%[1]s"
  topic_urn                 = huaweicloud_smn_topic.test.topic_urn
  notice_class              = "threat_alert_notice"
  enterprise_project_id     = "%[2]s"
  enabled                   = true
  sendfreq                  = 120
  locale                    = "en-us"
  times                     = 5
  threat                    = ["anticrawler", "cc"]
  is_all_enterprise_project = true
  description               = "test-alarm-notification"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAlarmNotification_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_waf_alarm_notification" "test" {
  name                      = "%[1]s_update"
  topic_urn                 = huaweicloud_smn_topic.test.topic_urn
  notice_class              = "threat_alert_notice"
  enterprise_project_id     = "%[2]s"
  enabled                   = false
  sendfreq                  = 60
  locale                    = "en-us"
  times                     = 6
  threat                    = ["anticrawler", "cc"]
  is_all_enterprise_project = false
  description               = "test-alarm-notification update"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAlarmNotificationImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}
		epsID := rs.Primary.Attributes["enterprise_project_id"]
		if epsID == "" {
			return "", errors.New("enterprise_project_id is empty")
		}
		return fmt.Sprintf("%s/%s", rs.Primary.ID, epsID), nil
	}
}

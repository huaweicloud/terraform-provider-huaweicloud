package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getNotificationRuleFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	return workspace.GetNotificationRuleById(client, state.Primary.ID)
}

func TestAccNotificationRule_basic(t *testing.T) {
	var (
		resourceName     = "huaweicloud_workspace_notification_rule.test"
		notificationRule interface{}
		rc               = acceptance.InitResourceCheck(resourceName, &notificationRule, getNotificationRuleFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationRule_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "metric_name", "desktop_idle_duration"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "notify_object"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "30"),
					resource.TestCheckResourceAttr(resourceName, "interval", "1"),
				),
			},
			{
				Config: testAccNotificationRule_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "metric_name", "desktop_idle_duration"),
					resource.TestCheckResourceAttr(resourceName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(resourceName, "enable", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "notify_object"),
					resource.TestCheckResourceAttr(resourceName, "threshold", "15"),
					resource.TestCheckResourceAttr(resourceName, "interval", "2"),
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

func testAccNotificationRule_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_workspace_notification_rule" "test" {
  depends_on = [huaweicloud_smn_topic.test]

  metric_name         = "desktop_idle_duration"
  comparison_operator = ">="
  enable             = true
  notify_object      = huaweicloud_smn_topic.test.topic_urn
  threshold          = 30
  interval           = 1
}
`, acceptance.RandomAccResourceNameWithDash())
}

func testAccNotificationRule_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_workspace_notification_rule" "test" {
  depends_on = [huaweicloud_smn_topic.test]

  metric_name         = "desktop_idle_duration"
  comparison_operator = ">="
  enable             = false
  notify_object      = huaweicloud_smn_topic.test.topic_urn
  threshold          = 15
  interval           = 2
}
`, acceptance.RandomAccResourceNameWithDash())
}

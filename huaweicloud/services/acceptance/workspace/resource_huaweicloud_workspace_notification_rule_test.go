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
		rName            = "huaweicloud_workspace_notification_rule.test"
		notificationRule interface{}
		rc               = acceptance.InitResourceCheck(rName, &notificationRule, getNotificationRuleFunc)
		baseConfig       = testAccNotificationRule_base(acceptance.RandomAccResourceNameWithDash())
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
				Config: testAccNotificationRule_basic_step1(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metric_name", "desktop_idle_duration"),
					resource.TestCheckResourceAttr(rName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(rName, "enable", "true"),
					resource.TestCheckResourceAttrSet(rName, "notify_object"),
					resource.TestCheckResourceAttr(rName, "threshold", "30"),
					resource.TestCheckResourceAttr(rName, "interval", "1"),
				),
			},
			{
				Config: testAccNotificationRule_basic_step2(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metric_name", "desktop_idle_duration"),
					resource.TestCheckResourceAttr(rName, "comparison_operator", ">="),
					resource.TestCheckResourceAttr(rName, "enable", "false"),
					resource.TestCheckResourceAttrSet(rName, "notify_object"),
					resource.TestCheckResourceAttr(rName, "threshold", "15"),
					resource.TestCheckResourceAttr(rName, "interval", "2"),
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

func testAccNotificationRule_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}
`, name)
}

func testAccNotificationRule_basic_step1(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_notification_rule" "test" {
  depends_on = [huaweicloud_smn_topic.test]

  metric_name         = "desktop_idle_duration"
  comparison_operator = ">="
  enable              = true
  notify_object       = huaweicloud_smn_topic.test.topic_urn
  threshold           = 30
  interval            = 1
}
`, baseConfig)
}

func testAccNotificationRule_basic_step2(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_notification_rule" "test" {
  depends_on = [huaweicloud_smn_topic.test]

  metric_name         = "desktop_idle_duration"
  comparison_operator = ">="
  enable              = false
  notify_object       = huaweicloud_smn_topic.test.topic_urn
  threshold           = 15
  interval            = 2
}
`, baseConfig)
}

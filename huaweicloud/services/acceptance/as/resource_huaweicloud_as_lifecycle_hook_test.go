package as

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/lifecyclehooks"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccASLifecycleHook_basic(t *testing.T) {
	var hook lifecyclehooks.Hook
	rName := acceptance.RandomAccResourceName()
	resourceGroupName := "huaweicloud_as_group.acc_as_group"
	resourceHookName := "huaweicloud_as_lifecycle_hook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASLifecycleHookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASLifecycleHook_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASLifecycleHookExists(resourceGroupName, resourceHookName, &hook),
					resource.TestCheckResourceAttr(resourceHookName, "name", rName),
					resource.TestCheckResourceAttr(resourceHookName, "type", "ADD"),
					resource.TestCheckResourceAttr(resourceHookName, "default_result", "ABANDON"),
					resource.TestCheckResourceAttr(resourceHookName, "timeout", "3600"),
					resource.TestCheckResourceAttr(resourceHookName, "notification_message", "This is a test message"),
					resource.TestMatchResourceAttr(resourceHookName, "notification_topic_urn",
						regexp.MustCompile(fmt.Sprintf(`^(urn:smn:%s:[0-9a-z]{32}:%s)$`, acceptance.HW_REGION_NAME, rName))),
				),
			},
			{
				Config: testASLifecycleHook_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASLifecycleHookExists(resourceGroupName, resourceHookName, &hook),
					resource.TestCheckResourceAttr(resourceHookName, "name", rName),
					resource.TestCheckResourceAttr(resourceHookName, "type", "REMOVE"),
					resource.TestCheckResourceAttr(resourceHookName, "default_result", "CONTINUE"),
					resource.TestCheckResourceAttr(resourceHookName, "timeout", "600"),
					resource.TestCheckResourceAttr(resourceHookName, "notification_message", "This is a update message"),
					resource.TestMatchResourceAttr(resourceHookName, "notification_topic_urn",
						regexp.MustCompile(fmt.Sprintf(`^(urn:smn:%s:[0-9a-z]{32}:%s-update)$`, acceptance.HW_REGION_NAME, rName))),
				),
			},
			{
				ResourceName:      resourceHookName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccASLifecycleHookImportStateIdFunc(resourceGroupName, resourceHookName),
			},
		},
	})
}

func testAccCheckASLifecycleHookDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating autoscaling client: %s", err)
	}

	var groupID string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_as_group" {
			groupID = rs.Primary.ID
			break
		}
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_as_lifecycle_hook" {
			continue
		}

		_, err := lifecyclehooks.Get(asClient, groupID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS lifecycle hook still exists")
		}
	}

	return nil
}

func testAccCheckASLifecycleHookExists(resGroup, resHook string, hook *lifecyclehooks.Hook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resGroup]
		if !ok {
			return fmt.Errorf("Not found: %s", resGroup)
		}
		groupID := rs.Primary.ID

		rs, ok = s.RootModule().Resources[resHook]
		if !ok {
			return fmt.Errorf("Not found: %s", resHook)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		asClient, err := config.AutoscalingV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating autoscaling client: %s", err)
		}
		found, err := lifecyclehooks.Get(asClient, groupID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		hook = found

		return nil
	}
}

func testAccASLifecycleHookImportStateIdFunc(groupRes, hookRes string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		group, ok := s.RootModule().Resources[groupRes]
		if !ok {
			return "", fmt.Errorf("Auto Scaling group not found: %s", group)
		}
		hook, ok := s.RootModule().Resources[hookRes]
		if !ok {
			return "", fmt.Errorf("Auto Scaling lifecycle hook not found: %s", hook)
		}
		if group.Primary.ID == "" || hook.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", group.Primary.ID, hook.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", group.Primary.ID, hook.Primary.ID), nil
	}
}

func testASLifecycleHook_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_topic" "test" {
  name = "%s"
}

resource "huaweicloud_smn_topic" "update" {
  name = "%s-update"
}
`, testASGroup_basic(rName), rName, rName)
}

func testASLifecycleHook_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_lifecycle_hook" "test" {
  name                   = "%s"
  type                   = "ADD"
  scaling_group_id       = huaweicloud_as_group.acc_as_group.id
  notification_topic_urn = huaweicloud_smn_topic.test.topic_urn
  notification_message   = "This is a test message"
}
`, testASLifecycleHook_base(rName), rName)
}

func testASLifecycleHook_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_lifecycle_hook" "test" {
  name                   = "%s"
  type                   = "REMOVE"
  scaling_group_id       = huaweicloud_as_group.acc_as_group.id
  default_result         = "CONTINUE"
  notification_topic_urn = huaweicloud_smn_topic.update.topic_urn
  notification_message   = "This is a update message"
  timeout                = 600
}
`, testASLifecycleHook_base(rName), rName)
}

package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/lifecyclehooks"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccASLifecycleHook_basic(t *testing.T) {
	var hook lifecyclehooks.Hook
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	// If the group name of the testASV1Group_basic method is updated, the resource name must also be updated.
	resourceGroupName := "huaweicloud_as_group.hth_as_group"
	resourceHookName := "huaweicloud_as_lifecycle_hook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASLifecycleHookDestroy,
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
						regexp.MustCompile(fmt.Sprintf("^(urn:smn:%s:%s:%s)$", HW_REGION_NAME, HW_PROJECT_ID, rName))),
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
					resource.TestCheckResourceAttr(resourceHookName, "notification_message",
						"This is a update message"),
					resource.TestMatchResourceAttr(resourceHookName, "notification_topic_urn",
						regexp.MustCompile(fmt.Sprintf("^(urn:smn:%s:%s:%s-update)$",
							HW_REGION_NAME, HW_PROJECT_ID, rName))),
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
	config := testAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating huaweicloud autoscaling client: %s", err)
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
			return fmtp.Errorf("AS lifecycle hook still exists")
		}
	}

	return nil
}

func testAccCheckASLifecycleHookExists(resGroup, resHook string, hook *lifecyclehooks.Hook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resGroup]
		if !ok {
			return fmtp.Errorf("Not found: %s", resGroup)
		}
		groupID := rs.Primary.ID

		rs, ok = s.RootModule().Resources[resHook]
		if !ok {
			return fmtp.Errorf("Not found: %s", resHook)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		asClient, err := config.AutoscalingV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating huaweicloud autoscaling client: %s", err)
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
			return "", fmtp.Errorf("Auto Scaling group not found: %s", group)
		}
		hook, ok := s.RootModule().Resources[hookRes]
		if !ok {
			return "", fmtp.Errorf("Auto Scaling lifecycle hook not found: %s", hook)
		}
		if group.Primary.ID == "" || hook.Primary.ID == "" {
			return "", fmtp.Errorf("resource not found: %s/%s", group.Primary.ID, hook.Primary.ID)
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
`, testASV1Group_basic(rName), rName, rName)
}

func testASLifecycleHook_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_lifecycle_hook" "test" {
  name                   = "%s"
  type                   = "ADD"
  scaling_group_id       = huaweicloud_as_group.hth_as_group.id
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
  scaling_group_id       = huaweicloud_as_group.hth_as_group.id
  default_result         = "CONTINUE"
  notification_topic_urn = huaweicloud_smn_topic.update.topic_urn
  notification_message   = "This is a update message"
  timeout                = 600
}
`, testASLifecycleHook_base(rName), rName)
}

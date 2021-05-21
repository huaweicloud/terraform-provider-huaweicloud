package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/lifecyclehooks"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccASLifecycleHook_basic(t *testing.T) {
	var hook lifecyclehooks.Hook
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceGroupName := "huaweicloud_as_group.test"
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
					resource.TestCheckResourceAttr(resourceHookName, "notification_topic_urn",
						fmt.Sprintf("urn:smn:%s:%s:default", HW_REGION_NAME, HW_PROJECT_ID)),
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
					resource.TestCheckResourceAttr(resourceHookName, "notification_topic_urn",
						fmt.Sprintf("urn:smn:%s:%s:update", HW_REGION_NAME, HW_PROJECT_ID)),
				),
			},
			{
				ResourceName:      resourceHookName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccASLifecycleHookImportStateIdFunc(),
			},
		},
	})
}

func testAccCheckASLifecycleHookDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud autoscaling client: %s", err)
	}

	var groupID string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_as_group" {
			groupID = rs.Primary.ID
			continue
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

		config := testAccProvider.Meta().(*config.Config)
		asClient, err := config.AutoscalingV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud autoscaling client: %s", err)
		}
		found, err := lifecyclehooks.Get(asClient, groupID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		hook = found

		return nil
	}
}

func testAccASLifecycleHookImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		group, ok := s.RootModule().Resources["huaweicloud_as_group.test"]
		if !ok {
			return "", fmt.Errorf("Auto Scaling group not found: %s", group)
		}
		hook, ok := s.RootModule().Resources["huaweicloud_as_lifecycle_hook.test"]
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
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
}

resource "huaweicloud_compute_keypair" "test" {
  name = "%s"

  lifecycle {
    ignore_changes = [
      public_key,
    ]
  }
}

resource "huaweicloud_lb_loadbalancer" "test" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.subnet_id
}

resource "huaweicloud_lb_listener" "test" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.test.id
}

resource "huaweicloud_lb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.test.id
}

resource "huaweicloud_as_configuration" "test"{
  scaling_configuration_name = "%s"

  instance_config {
    image    = data.huaweicloud_images_image.test.id
    flavor   = data.huaweicloud_compute_flavors.test.ids[0]
    key_name = huaweicloud_compute_keypair.test.id

    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
  }
}

resource "huaweicloud_as_group" "test"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.test.id
  vpc_id                   = data.huaweicloud_vpc.test.id

  networks {
    id = data.huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.test.id
    protocol_port = huaweicloud_lb_listener.test.protocol_port
  }
}
`, rName, rName, rName, rName, rName, rName, rName)
}

// Please make sure the smn topic is exist in specifies region and project.
func testASLifecycleHook_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_lifecycle_hook" "test" {
  name                   = "%s"
  type                   = "ADD"
  scaling_group_id       = huaweicloud_as_group.test.id
  default_result         = "ABANDON"
  notification_topic_urn = "%s"
  notification_message   = "This is a test message"
}	  
`, testASLifecycleHook_base(rName), rName, fmt.Sprintf("urn:smn:%s:%s:default", HW_REGION_NAME, HW_PROJECT_ID))
}

// Please make sure the smn topic is exist in specifies region and project.
func testASLifecycleHook_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_lifecycle_hook" "test" {
  name                   = "%s"
  type                   = "REMOVE"
  scaling_group_id       = huaweicloud_as_group.test.id
  default_result         = "CONTINUE"
  notification_topic_urn = "%s"
  notification_message   = "This is a update message"
  timeout                = 600
}	  
`, testASLifecycleHook_base(rName), rName, fmt.Sprintf("urn:smn:%s:%s:update", HW_REGION_NAME, HW_PROJECT_ID))
}

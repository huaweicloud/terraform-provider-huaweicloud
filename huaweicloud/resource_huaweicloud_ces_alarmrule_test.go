package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/cloudeyeservice/alarmrule"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccCESAlarmRule_basic(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))
	resourceName := "huaweicloud_ces_alarmrule.alarmrule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCESAlarmRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testCESAlarmRuleExists(resourceName, &ar),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "6"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testCESAlarmRule_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "3"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "60"),
				),
			},
		},
	})
}

func testCESAlarmRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.CesV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud ces client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_ces_alarmrule" {
			continue
		}

		id := rs.Primary.ID
		_, err := alarmrule.Get(networkingClient, id).Extract()
		if err == nil {
			return fmt.Errorf("Alarm rule still exists")
		}
	}

	return nil
}

func testCESAlarmRuleExists(n string, ar *alarmrule.AlarmRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.CesV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud ces client: %s", err)
		}

		id := rs.Primary.ID
		found, err := alarmrule.Get(networkingClient, id).Extract()
		if err != nil {
			return err
		}

		*ar = *found

		return nil
	}
}

func testCESAlarmRule_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_compute_instance" "vm_1" {
  name              = "ecs-%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
resource "huaweicloud_smn_topic" "topic_1" {
  name         = "smn-%s"
  display_name = "The display name of smn topic"
}
`, rName, rName)
}

func testCESAlarmRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.vm_1.id
    }
  }
  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6
    unit                = "B/s"
    count               = 1
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.topic_1.topic_urn
    ]
  }
}
`, testCESAlarmRule_base(rName), rName)
}

func testCESAlarmRule_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "rule-%s-update"
  alarm_action_enabled = true
  alarm_enabled        = false
  alarm_level          = 3

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.vm_1.id
    }
  }
  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 60
    unit                = "B/s"
    count               = 1
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.topic_1.topic_urn
    ]
  }
}
`, testCESAlarmRule_base(rName), rName)
}

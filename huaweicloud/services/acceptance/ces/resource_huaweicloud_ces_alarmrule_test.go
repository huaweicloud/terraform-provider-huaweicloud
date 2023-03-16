package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cloudeyeservice/v1/alarmrule"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getAlarmRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CesV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud CES v1 client: %s", err)
	}
	return alarmrule.Get(c, state.Primary.ID).Extract()
}

func TestAccCESAlarmRule_basic(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.alarmrule_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_type", "MULTI_INSTANCE"),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "6.5"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.suppress_duration", "300"),
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

func TestAccCESAlarmRule_withEpsId(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.alarmrule_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "6"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccCESAlarmRule_sysEvent(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.alarmrule_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_sysEvent(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "1"),
				),
			},
		},
	})
}

func TestAccCESAlarmRule_multiConditions(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.alarmrule_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_multiConditions(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.alarm_level", "3"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "6.5"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.period", "300"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.metric_name", "network_incoming_bytes_rate_inband"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.alarm_level", "3"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.value", "6.5"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.period", "300"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.metric_name", "network_outgoing_bytes_rate_inband"),
				),
			},
			{
				Config: testCESAlarmRule_multiConditions_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s-update", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.alarm_level", "4"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "6.5"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.period", "1200"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.metric_name", "network_outgoing_bytes_rate_inband"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.alarm_level", "4"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.value", "20"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.period", "3600"),
					resource.TestCheckResourceAttr(resourceName, "condition.1.metric_name", "network_outgoing_bytes_rate_inband"),
				),
			},
		},
	})
}

func testCESAlarmRule_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "vm_1" {
  name               = "ecs-%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_smn_topic" "topic_1" {
  name         = "smn-%[2]s"
  display_name = "The display name of smn topic"
}
`, common.TestBaseComputeResources(rName), rName)
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
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
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
    suppress_duration   = 300
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

func testCESAlarmRule_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "alarmrule_1" {
  alarm_name            = "rule-%s"
  alarm_action_enabled  = true
  enterprise_project_id = "%s"

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
    suppress_duration   = 300
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.topic_1.topic_urn
    ]
  }
}
`, testCESAlarmRule_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testCESAlarmRule_sysEvent(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true
  alarm_type           = "EVENT.SYS"

  metric {
    namespace   = "SYS.ECS"
    metric_name = "stopServer"
  }
  
  condition  {
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    suppress_duration   = 0
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

func testCESAlarmRule_multiConditions(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true
  alarm_type           = "MULTI_INSTANCE"

  metric {
    namespace   = "SYS.ECS"

    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.vm_1.id
    }
  }

  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    metric_name         = "network_incoming_bytes_rate_inband"
    alarm_level         = 3
  }

  condition  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    metric_name         = "network_outgoing_bytes_rate_inband"
    alarm_level         = 3
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

func testCESAlarmRule_multiConditions_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "rule-%s-update"
  alarm_action_enabled = true
  alarm_enabled        = false
  alarm_type           = "MULTI_INSTANCE"

  metric {
    namespace   = "SYS.ECS"

    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.vm_1.id
    }
  }

  condition  {
    period              = 1200
    filter              = "average"
    comparison_operator = ">"
    value               = 6.5
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    metric_name         = "network_outgoing_bytes_rate_inband"
    alarm_level         = 4
  }

  condition  {
    period              = 3600
    filter              = "average"
    comparison_operator = ">="
    value               = 20
    unit                = "B/s"
    count               = 1
    suppress_duration   = 300
    metric_name         = "network_outgoing_bytes_rate_inband"
    alarm_level         = 4
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

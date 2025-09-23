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
	resourceName := "huaweicloud_ces_alarmrule.test"

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
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testCESAlarmRule_update(rName),
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

func TestAccCESAlarmRule_withEpsId(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&ar,
		getAlarmRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCESAlarmRule_basic(rName),
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
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testCESAlarmRule_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.alarm_level", "3"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "6.5"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.period", "300"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.metric_name", "network_incoming_bytes_rate_inband"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccCESAlarmRule_sysEvent(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.test"

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
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.alarm_level", "2"),
				),
			},
		},
	})
}

func TestAccCESAlarmRule_allInstance(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.test"

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
				Config: testCESAlarmRule_allInstance(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", fmt.Sprintf("rule-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_type", "ALL_INSTANCE"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.alarm_level", "2"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "80"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.period", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.metric_name", "disk_usedPercent"),
				),
			},
		},
	})
}

func TestAccCESAlarmRule_withResourceGroup(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.test"

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
				Config: testCESAlarmRule_withResourceGroup(rName),
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
					resource.TestCheckResourceAttrPair(resourceName, "resource_group_id",
						"huaweicloud_ces_resource_group.test", "id"),
				),
			},
		},
	})
}

func TestAccCESAlarmRule_old(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.test"

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
				Config: testCESAlarmRule_old(rName),
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
				Config: testCESAlarmRule_old_update(rName),
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

func TestAccCESAlarmRule_withAlarmTemplate(t *testing.T) {
	var ar alarmrule.AlarmRule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ces_alarmrule.test"

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
				Config: testCESAlarmRule_withAlarmTemplate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "alarm_name", rName),
					resource.TestCheckResourceAttr(resourceName, "alarm_type", "MULTI_INSTANCE"),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_action_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "alarm_level", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "alarm_template_id",
						"huaweicloud_ces_alarm_template.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.metric_name", "disk_util_inband"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.period", "1"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.filter", "average"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.value", "90"),
					resource.TestCheckResourceAttr(resourceName, "condition.0.suppress_duration", "86400"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.dimensions.0.value", "huaweicloud_compute_instance.test.0", "id"),
				),
			},
			{
				Config: testCESAlarmRule_withAlarmTemplate_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "alarm_name", rName),
					resource.TestCheckResourceAttr(resourceName, "alarm_enabled", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.dimensions.0.value", "huaweicloud_compute_instance.test.1", "id"),
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

func testCESAlarmRule_instanceBase(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = "ecs-%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}
`, common.TestBaseComputeResources(rName), rName)
}

func testCESAlarmRule_topicBase(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name         = "smn-%s"
  display_name = "The display name of smn topic"
}
`, rName)
}

func testCESAlarmRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true
  alarm_type           = "MULTI_INSTANCE"

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[0].id
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
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_update(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%s-update"
  alarm_action_enabled = true
  alarm_enabled        = false
  alarm_type           = "MULTI_INSTANCE"

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[0].id
    }
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[1].id
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
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name            = "rule-%s"
  alarm_action_enabled  = true
  alarm_type            = "MULTI_INSTANCE"
  enterprise_project_id = "%s"

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[0].id
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

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testCESAlarmRule_sysEvent(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true
  alarm_type           = "EVENT.SYS"

  metric {
    namespace = "SYS.ECS"
  }
  
  condition  {
    metric_name         = "stopServer"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = 1
    unit                = "count"
    count               = 1
    suppress_duration   = 0
    alarm_level         = 2
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_allInstance(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true
  alarm_type           = "ALL_INSTANCE"

  metric {
    namespace = "AGT.ECS"
  }

  resources {
    dimensions {
      name = "instance_id"
    }

    dimensions {
      name = "mount_point"
    }
  }

  condition  {
    alarm_level         = 2
    suppress_duration   = 0
    period              = 1
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    count               = 1
    metric_name         = "disk_usedPercent"
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_withResourceGroup(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_ces_resource_group" "test" {
  name = "test"

  resources {
    namespace = "SYS.ECS"
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[0].id
    }
  }
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true
  alarm_type           = "RESOURCE_GROUP"
  resource_group_id    = huaweicloud_ces_resource_group.test.id

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
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
      huaweicloud_smn_topic.test.topic_urn
    ]
  }

  lifecycle {
    ignore_changes = [
      resources
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_old(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%s"
  alarm_action_enabled = true

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"

    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[0].id
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
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_old_update(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%s-update"
  alarm_action_enabled = true
  alarm_enabled        = false
  alarm_level          = 3

  metric {
    namespace   = "SYS.ECS"
    metric_name = "network_outgoing_bytes_rate_inband"

    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[0].id
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
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_withAlarmTemplate(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_ces_alarm_template" "test" {
  name        = "%[3]s"
  description = "It is an template"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "disk_util_inband"
    period              = 1
    filter              = "average"
    comparison_operator = ">="
    value               = "90"
    unit                = "%%"
    count               = 3
    alarm_level         = 1
    suppress_duration   = 86400
  }

  depends_on = [
    huaweicloud_compute_instance.test,
    huaweicloud_smn_topic.test
  ]
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "%[3]s"
  alarm_enabled        = true
  alarm_action_enabled = true
  alarm_type           = "MULTI_INSTANCE"
  alarm_template_id    = huaweicloud_ces_alarm_template.test.id

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[0].id
    }
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName)
}

func testCESAlarmRule_withAlarmTemplate_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_ces_alarm_template" "test" {
  name        = "%[3]s"
  description = "It is an template"

  policies {
    namespace           = "SYS.ECS"
    dimension_name      = "instance_id"
    metric_name         = "disk_util_inband"
    period              = 1
    filter              = "average"
    comparison_operator = ">="
    value               = "90"
    unit                = "%%"
    count               = 3
    alarm_level         = 1
    suppress_duration   = 86400
  }

  depends_on = [
    huaweicloud_compute_instance.test,
    huaweicloud_smn_topic.test
  ]
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "%[3]s"
  alarm_action_enabled = true
  alarm_enabled        = false
  alarm_type           = "MULTI_INSTANCE"
  alarm_template_id    = huaweicloud_ces_alarm_template.test.id

  metric {
    namespace = "SYS.ECS"
  }

  resources {
    dimensions {
      name  = "instance_id"
      value = huaweicloud_compute_instance.test[1].id
    }
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }
}
`, testCESAlarmRule_instanceBase(rName), testCESAlarmRule_topicBase(rName), rName)
}

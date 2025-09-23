package iotda

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

func getDeviceLinkageRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "iotda"
		httpUrl = "v5/iot/{project_id}/rules/{rule_id}"
	)

	isDerived := WithDerivedAuth()
	client, err := conf.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccDeviceLinkageRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_device_linkage_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceLinkageRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeviceLinkageRule_deviceData(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "DEVICE_DATA"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_data_condition.0.path", "temp_1/demo_1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_data_condition.0.operator", "="),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_data_condition.0.value", "4"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_data_condition.0.trigger_strategy", "pulse"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_data_condition.0.data_validatiy_period",
						"300"),
					resource.TestCheckResourceAttrPair(rName, "triggers.0.device_data_condition.0.product_id",
						"huaweicloud_iotda_product.test", "id"),
					resource.TestCheckResourceAttr(rName, "actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "actions.0.type", "SMN_FORWARDING"),
					resource.TestCheckResourceAttr(rName, "actions.0.smn_forwarding.0.message_title", "title"),
					resource.TestCheckResourceAttr(rName, "actions.0.smn_forwarding.0.message_content", "content"),
					resource.TestCheckResourceAttrPair(rName, "actions.0.smn_forwarding.0.topic_name",
						"huaweicloud_smn_topic.topic", "name"),
					resource.TestCheckResourceAttrPair(rName, "actions.0.smn_forwarding.0.topic_urn",
						"huaweicloud_smn_topic.topic", "topic_urn"),
					resource.TestCheckResourceAttrPair(rName, "actions.0.smn_forwarding.0.message_template_name",
						"huaweicloud_smn_message_template.test", "name"),
				),
			},
			{
				Config: testDeviceLinkageRule_deviceStatus(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "DEVICE_LINKAGE_STATUS"),
					resource.TestCheckResourceAttrPair(rName, "triggers.0.device_linkage_status_condition.0.device_id",
						"huaweicloud_iotda_device.test", "id"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_linkage_status_condition.0.status_list.#", "2"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_linkage_status_condition.0.duration", "10"),
					resource.TestCheckResourceAttr(rName, "actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "actions.0.type", "DEVICE_ALARM"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_alarm.0.name", name),
					resource.TestCheckResourceAttr(rName, "actions.0.device_alarm.0.type", "fault"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_alarm.0.severity", "warning"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_alarm.0.dimension", "device"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_alarm.0.description", "description test"),
				),
			},
			{
				Config: testDeviceLinkageRule_simpleTimer(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "SIMPLE_TIMER"),
					resource.TestCheckResourceAttr(rName, "triggers.0.simple_timer_condition.0.start_time", "20220622T160000Z"),
					resource.TestCheckResourceAttr(rName, "triggers.0.simple_timer_condition.0.repeat_interval", "2"),
					resource.TestCheckResourceAttr(rName, "triggers.0.simple_timer_condition.0.repeat_count", "2"),
					resource.TestCheckResourceAttr(rName, "actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "actions.0.type", "SMN_FORWARDING"),
					resource.TestCheckResourceAttr(rName, "actions.0.smn_forwarding.0.message_title", "title"),
					resource.TestCheckResourceAttr(rName, "actions.0.smn_forwarding.0.message_content", "content"),
					resource.TestCheckResourceAttrPair(rName, "actions.0.smn_forwarding.0.topic_name",
						"huaweicloud_smn_topic.topic", "name"),
					resource.TestCheckResourceAttrPair(rName, "actions.0.smn_forwarding.0.topic_urn",
						"huaweicloud_smn_topic.topic", "topic_urn"),
				),
			},
			{
				Config: testDeviceLinkageRule_dailyTimer(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "DAILY_TIMER"),
					resource.TestCheckResourceAttr(rName, "triggers.0.daily_timer_condition.0.start_time", "19:02"),
					resource.TestCheckResourceAttr(rName, "triggers.0.daily_timer_condition.0.days_of_week",
						"1,2,3,4,5,6,7"),
					resource.TestCheckResourceAttr(rName, "actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "actions.0.type", "DEVICE_CMD"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.service_id", "temp_2"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.command_name", "cmd_1"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.command_body",
						"{\"cmd_p_1\":\"3\"}"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.buffer_timeout", "180"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.response_timeout", "60"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.mode", "PASSIVE"),
					resource.TestCheckResourceAttrPair(rName, "actions.0.device_command.0.device_id",
						"huaweicloud_iotda_device.test", "id"),
					resource.TestCheckResourceAttr(rName, "effective_period.0.start_time", "00:00"),
					resource.TestCheckResourceAttr(rName, "effective_period.0.end_time", "23:59"),
					resource.TestCheckResourceAttr(rName, "effective_period.0.days_of_week", "1,2,3,4,5,6"),
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

func testDeviceLinkageRuleBase(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic" {
  name = "%[1]s"
}

resource "huaweicloud_smn_message_template" "test" {
  name     = "%[1]s"
  protocol = "default"
  content  = "content test"
}
`, name)
}

func testDeviceLinkageRule_deviceData(name string) string {
	deviceConfig := testDevice_basic(name, name)
	deviceLinkageRuleBase := testDeviceLinkageRuleBase(name)

	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  name     = "%[3]s"
  space_id = huaweicloud_iotda_space.test.id

  triggers {
    type = "DEVICE_DATA"

    device_data_condition {
      product_id            = huaweicloud_iotda_device.test.product_id
      path                  = "temp_1/demo_1"
      operator              = "="
      value                 = 4
      trigger_strategy      = "pulse"
      data_validatiy_period = 300
    }
  }

  actions {
    type = "SMN_FORWARDING"

    smn_forwarding {
      region                = huaweicloud_smn_topic.topic.region
      topic_name            = huaweicloud_smn_topic.topic.name
      topic_urn             = huaweicloud_smn_topic.topic.topic_urn
      message_title         = "title"
      message_content       = "content"
      message_template_name = huaweicloud_smn_message_template.test.name
    }
  }

  depends_on = [
    huaweicloud_iotda_space.test,
    huaweicloud_iotda_product.test,
    huaweicloud_iotda_device.test,
    huaweicloud_smn_topic.topic,
    huaweicloud_smn_message_template.test,
  ]
}
`, deviceLinkageRuleBase, deviceConfig, name)
}

func testDeviceLinkageRule_deviceStatus(name string) string {
	deviceConfig := testDevice_basic(name, name)
	deviceLinkageRuleBase := testDeviceLinkageRuleBase(name)

	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  name     = "%[3]s_update"
  space_id = huaweicloud_iotda_space.test.id
  enabled  = false

  triggers {
    type = "DEVICE_LINKAGE_STATUS"

    device_linkage_status_condition {
      device_id   = huaweicloud_iotda_device.test.id
      status_list = ["ONLINE", "OFFLINE"]
      duration    = 10
    }
  }

  actions {
    type = "DEVICE_ALARM"

    device_alarm {
      name        = "%[3]s"
      type        = "fault"
      severity    = "warning"
      dimension   = "device"
      description = "description test"
    }
  }

  depends_on = [
    huaweicloud_iotda_device.test,
    huaweicloud_iotda_product.test,
    huaweicloud_smn_topic.topic,
    huaweicloud_smn_message_template.test,
  ]
}
`, deviceLinkageRuleBase, deviceConfig, name)
}

func testDeviceLinkageRule_simpleTimer(name string) string {
	deviceConfig := testDevice_basic(name, name)
	deviceLinkageRuleBase := testDeviceLinkageRuleBase(name)

	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  name     = "%[3]s"
  space_id = huaweicloud_iotda_space.test.id
  enabled  = true

  triggers {
    type = "SIMPLE_TIMER"

    simple_timer_condition {
      start_time      = "20220622T160000Z"
      repeat_interval = 2
      repeat_count    = 2
    }
  }

  actions {
    type = "SMN_FORWARDING"
    smn_forwarding {
      region                = huaweicloud_smn_topic.topic.region
      topic_name            = huaweicloud_smn_topic.topic.name
      topic_urn             = huaweicloud_smn_topic.topic.topic_urn
      message_title         = "title"
      message_content       = "content"
      message_template_name = huaweicloud_smn_message_template.test.name
    }
  }

  depends_on = [
    huaweicloud_iotda_device.test,
    huaweicloud_iotda_product.test,
    huaweicloud_smn_topic.topic,
    huaweicloud_smn_message_template.test,
  ]
}
`, deviceLinkageRuleBase, deviceConfig, name)
}

func testDeviceLinkageRule_dailyTimer(name string) string {
	deviceConfig := testDevice_basic(name, name)
	deviceLinkageRuleBase := testDeviceLinkageRuleBase(name)

	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  name     = "%[3]s"
  space_id = huaweicloud_iotda_space.test.id
  enabled  = true

  triggers {
    type = "DAILY_TIMER"

    daily_timer_condition {
      start_time   = "19:02"
      days_of_week = "1,2,3,4,5,6,7"
    }
  }

  actions {
    type = "DEVICE_CMD"

    device_command {
      device_id        = huaweicloud_iotda_device.test.id
      service_id       = "temp_2"
      command_name     = "cmd_1"
      command_body     = "{\"cmd_p_1\":\"3\"}"
      buffer_timeout   = 180
      response_timeout = 60
      mode             = "PASSIVE"
    }
  }

  effective_period {
    start_time   = "00:00"
    end_time     = "23:59"
    days_of_week = "1,2,3,4,5,6"
  }

  depends_on = [
    huaweicloud_iotda_device.test,
    huaweicloud_iotda_product.test,
    huaweicloud_smn_topic.topic,
    huaweicloud_smn_message_template.test,
  ]
}
`, deviceLinkageRuleBase, deviceConfig, name)
}

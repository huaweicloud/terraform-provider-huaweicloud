package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDeviceLinkageRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	return client.ShowRule(&model.ShowRuleRequest{RuleId: state.Primary.ID})
}

func TestAccDeviceLinkageRule_basic(t *testing.T) {
	var obj model.ShowRuleResponse

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_device_linkage_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceLinkageRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeviceLinkageRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "DEVICE_DATA"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_data_condition.0.path", "service_1/p_1"),
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
				),
			},
			{
				Config: testDeviceLinkageRule_timer(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
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
				Config: testDeviceLinkageRule_daily(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "DAILY_TIMER"),
					resource.TestCheckResourceAttr(rName, "triggers.0.daily_timer_condition.0.start_time", "19:02"),
					resource.TestCheckResourceAttr(rName, "triggers.0.daily_timer_condition.0.days_of_week",
						"1,2,3,4,5,6,7"),
					resource.TestCheckResourceAttr(rName, "actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "actions.0.type", "DEVICE_CMD"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.service_id", "service_1"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.command_name", "cmd_1"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.command_body",
						"{\"cmd_p_1\":\"3\"}"),
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

func TestAccDeviceLinkageRule_derived(t *testing.T) {
	var obj model.ShowRuleResponse

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
				Config: testDeviceLinkageRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "DEVICE_DATA"),
					resource.TestCheckResourceAttr(rName, "triggers.0.device_data_condition.0.path", "service_1/p_1"),
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
				),
			},
			{
				Config: testDeviceLinkageRule_timer(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
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
				Config: testDeviceLinkageRule_daily(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "trigger_logic", "and"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "triggers.#", "1"),
					resource.TestCheckResourceAttr(rName, "triggers.0.type", "DAILY_TIMER"),
					resource.TestCheckResourceAttr(rName, "triggers.0.daily_timer_condition.0.start_time", "19:02"),
					resource.TestCheckResourceAttr(rName, "triggers.0.daily_timer_condition.0.days_of_week",
						"1,2,3,4,5,6,7"),
					resource.TestCheckResourceAttr(rName, "actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "actions.0.type", "DEVICE_CMD"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.service_id", "service_1"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.command_name", "cmd_1"),
					resource.TestCheckResourceAttr(rName, "actions.0.device_command.0.command_body",
						"{\"cmd_p_1\":\"3\"}"),
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

func testDeviceLinkageRule_basic(name string) string {
	deviceConfig := testDevice_basic(name, name)
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_topic" "topic" {
  name = "%s"
}

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  name     = "%s"
  space_id = huaweicloud_iotda_space.test.id

  triggers {
    type = "DEVICE_DATA"
    device_data_condition {
      product_id            = huaweicloud_iotda_device.test.product_id
      path                  = "service_1/p_1"
      operator              = "="
      value                 = 4
      trigger_strategy      = "pulse"
      data_validatiy_period = 300
    }
  }

  actions {
    type = "SMN_FORWARDING"
    smn_forwarding {
      region          = huaweicloud_smn_topic.topic.region
      topic_name      = huaweicloud_smn_topic.topic.name
      topic_urn       = huaweicloud_smn_topic.topic.topic_urn
      message_title   = "title"
      message_content = "content"
    }
  }

  depends_on = [
    huaweicloud_iotda_device.test,
    huaweicloud_iotda_product.test,
  ]
}
`, deviceConfig, name, name)
}

func testDeviceLinkageRule_timer(name string) string {
	deviceConfig := testDevice_basic(name, name)
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_topic" "topic" {
  name = "%s"
}

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  name     = "%s"
  space_id = huaweicloud_iotda_space.test.id
  enabled  = false

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
      region          = huaweicloud_smn_topic.topic.region
      topic_name      = huaweicloud_smn_topic.topic.name
      topic_urn       = huaweicloud_smn_topic.topic.topic_urn
      message_title   = "title"
      message_content = "content"
    }
  }

  depends_on = [
    huaweicloud_iotda_device.test,
    huaweicloud_iotda_product.test,
  ]
}
`, deviceConfig, name, name)
}

func testDeviceLinkageRule_daily(name string) string {
	deviceConfig := testDevice_basic(name, name)
	return fmt.Sprintf(`
%s

resource "huaweicloud_smn_topic" "topic" {
  name = "%s"
}

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  name     = "%s"
  space_id = huaweicloud_iotda_space.test.id
  enabled  = false

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
      device_id    = huaweicloud_iotda_device.test.id
      service_id   = "service_1"
      command_name = "cmd_1"
      command_body = "{\"cmd_p_1\":\"3\"}"
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
  ]
}
`, deviceConfig, name, name)
}

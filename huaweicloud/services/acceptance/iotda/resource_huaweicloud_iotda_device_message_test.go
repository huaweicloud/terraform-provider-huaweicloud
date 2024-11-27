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

func getDeviceMessageResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	queryOpts := &model.ShowDeviceMessageRequest{
		DeviceId:  state.Primary.Attributes["device_id"],
		MessageId: state.Primary.ID,
	}
	return client.ShowDeviceMessage(queryOpts)
}

func TestAccDeviceMessage_basic(t *testing.T) {
	var (
		obj   model.ShowDeviceMessageResponse
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_iotda_device_message.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceMessageResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeviceMessage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "device_id", "huaweicloud_iotda_device.test", "id"),
					resource.TestCheckResourceAttr(rName, "message", "message_content"),
					resource.TestCheckResourceAttr(rName, "message_id", name),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "encoding", "none"),
					resource.TestCheckResourceAttr(rName, "payload_format", "standard"),
					resource.TestCheckResourceAttr(rName, "topic_full_name", "topic_test"),
					resource.TestCheckResourceAttr(rName, "properties.#", "1"),
					resource.TestCheckResourceAttr(rName, "properties.0.correlation_data", "data_test"),
					resource.TestCheckResourceAttr(rName, "properties.0.response_topic", "resp_test"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.#", "2"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.0.prop_key", "test_key1"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.0.prop_value", "test_val1"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.1.prop_key", "test_key2"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.1.prop_value", "test_val2"),
					resource.TestCheckResourceAttrSet(rName, "created_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
		},
	})
}

func TestAccDeviceMessage_derived(t *testing.T) {
	var (
		obj   model.ShowDeviceMessageResponse
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_iotda_device_message.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceMessageResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeviceMessage_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "device_id", "huaweicloud_iotda_device.test", "id"),
					resource.TestCheckResourceAttr(rName, "message", "message_content"),
					resource.TestCheckResourceAttr(rName, "message_id", name),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "encoding", "none"),
					resource.TestCheckResourceAttr(rName, "payload_format", "standard"),
					resource.TestCheckResourceAttr(rName, "topic_full_name", "topic_test"),
					resource.TestCheckResourceAttr(rName, "properties.#", "1"),
					resource.TestCheckResourceAttr(rName, "properties.0.correlation_data", "data_test"),
					resource.TestCheckResourceAttr(rName, "properties.0.response_topic", "resp_test"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.#", "2"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.0.prop_key", "test_key1"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.0.prop_value", "test_val1"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.1.prop_key", "test_key2"),
					resource.TestCheckResourceAttr(rName, "properties.0.user_properties.1.prop_value", "test_val2"),
					resource.TestCheckResourceAttrSet(rName, "created_time"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
		},
	})
}

func testDeviceMessage_base(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_device" "test" {
  node_id    = "%[3]s"
  name       = "%[3]s"
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
  secret     = "1234567890"
}
`, buildIoTDAEndpoint(), testAccDevice_base(name), name)
}

func testDeviceMessage_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_message" "test" {
  device_id       = huaweicloud_iotda_device.test.id
  message         = "message_content"
  message_id      = "%[2]s"
  name            = "%[2]s"
  encoding        = "none"
  payload_format  = "standard"
  topic_full_name = "topic_test"
  ttl             = "100"

  properties {
    correlation_data = "data_test"
    response_topic   = "resp_test"

    user_properties {
      prop_key   = "test_key1"
      prop_value = "test_val1"
    }

    user_properties {
      prop_key   = "test_key2"
      prop_value = "test_val2"
    }
  }
}
`, testDeviceMessage_base(name), name)
}

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDeviceMessageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/messages/{message_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{device_id}", state.Primary.Attributes["device_id"])
	getPath = strings.ReplaceAll(getPath, "{message_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When the parent resource (device_id) does not exist, query API will return `404` error code.
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	// When the resource does not exist, query API will return `200` status code.
	// Therefore, it is necessary to check whether the message ID is returned in the response body.
	messageIdResp := utils.PathSearch("message_id", getRespBody, "").(string)
	if messageIdResp == "" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccDeviceMessage_basic(t *testing.T) {
	var (
		obj   interface{}
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

data "huaweicloud_iotda_spaces" "test" {
  is_default = true
}

resource "huaweicloud_iotda_product" "test" {
  name        = "%[2]s"
  device_type = "test"
  protocol    = "MQTT"
  space_id    = data.huaweicloud_iotda_spaces.test.spaces[0].id
  data_type   = "json"

  services {
    id   = "service_1"
    type = "serv_type"
  }
}

resource "huaweicloud_iotda_device" "test" {
  node_id    = "%[2]s"
  name       = "%[2]s"
  space_id   = data.huaweicloud_iotda_spaces.test.spaces[0].id
  product_id = huaweicloud_iotda_product.test.id
  secret     = "1234567890"
}
`, buildIoTDAEndpoint(), name)
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

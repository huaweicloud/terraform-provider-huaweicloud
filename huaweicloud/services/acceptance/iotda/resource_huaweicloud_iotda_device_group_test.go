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

func getDeviceGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + "v5/iot/{project_id}/device-group/{group_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{group_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IoTDA device group: %s", err)

	}

	return utils.FlattenResponse(getResp)
}

func TestAccDeviceGroup_basic(t *testing.T) {
	var (
		obj        interface{}
		deviceName = acceptance.RandomAccResourceName()
		name       = acceptance.RandomAccResourceName()
		updateName = name + "_update"
		rName      = "huaweicloud_iotda_device_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceGroupResourceFunc,
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
				Config: testDeviceGroup_basic(name, deviceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "data.huaweicloud_iotda_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttr(rName, "device_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "device_ids.0", "huaweicloud_iotda_device.test", "id"),
				),
			},
			{
				Config: testDeviceGroup_basic_update(updateName, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "description update"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "data.huaweicloud_iotda_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttr(rName, "device_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "device_ids.0", "huaweicloud_iotda_device.test2", "id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_id"},
			},
		},
	})
}

func TestAccDeviceGroup_withDynamicGroup(t *testing.T) {
	var (
		obj        interface{}
		deviceName = acceptance.RandomAccResourceName()
		name       = acceptance.RandomAccResourceName()
		updateName = name + "_update"
		rName      = "huaweicloud_iotda_device_group.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceGroupResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Only the standard and enterprise versions of IoTDA instances support dynamic device groups.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDeviceGroup_dynamicGroup_basic(name, deviceName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttr(rName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "data.huaweicloud_iotda_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttr(rName, "device_ids.#", "2"),
				),
			},
			{
				Config: testDeviceGroup_dynamicGroup_update(updateName, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "description update"),
					resource.TestCheckResourceAttr(rName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "data.huaweicloud_iotda_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttr(rName, "device_ids.#", "2"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_id"},
			},
		},
	})
}

func testDeviceGroup_base(deviceName string) string {
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
  node_id    = "%[2]s_1"
  name       = "%[2]s_1"
  space_id   = data.huaweicloud_iotda_spaces.test.spaces[0].id
  product_id = huaweicloud_iotda_product.test.id
}

resource "huaweicloud_iotda_device" "test2" {
  node_id    = "%[2]s_2"
  name       = "%[2]s_2"
  space_id   = data.huaweicloud_iotda_spaces.test.spaces[0].id
  product_id = huaweicloud_iotda_product.test.id
}
`, buildIoTDAEndpoint(), deviceName)
}

func testDeviceGroup_basic(name, deviceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_group" "test" {
  name        = "%[2]s"
  space_id    = data.huaweicloud_iotda_spaces.test.spaces[0].id
  description = "description test"
  device_ids  = [huaweicloud_iotda_device.test.id]
}
`, testDeviceGroup_base(deviceName), name)
}

func testDeviceGroup_basic_update(name, deviceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_group" "test" {
  name        = "%[2]s"
  space_id    = data.huaweicloud_iotda_spaces.test.spaces[0].id
  description = "description update"
  device_ids  = [huaweicloud_iotda_device.test2.id]
}
`, testDeviceGroup_base(deviceName), name)
}

func testDeviceGroup_dynamicGroup_basic(name, deviceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_group" "test" {
  name               = "%[2]s"
  space_id           = data.huaweicloud_iotda_spaces.test.spaces[0].id
  description        = "description test"
  type               = "DYNAMIC"
  dynamic_group_rule = "device_id = '${huaweicloud_iotda_device.test.id}' or device_id = '${huaweicloud_iotda_device.test2.id}'"
}
`, testDeviceGroup_base(deviceName), name)
}

func testDeviceGroup_dynamicGroup_update(name, deviceName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_group" "test" {
  name               = "%[2]s"
  space_id           = data.huaweicloud_iotda_spaces.test.spaces[0].id
  description        = "description update"
  type               = "DYNAMIC"
  dynamic_group_rule = "device_id = '${huaweicloud_iotda_device.test.id}' or device_id = '${huaweicloud_iotda_device.test2.id}'"
}
`, testDeviceGroup_base(deviceName), name)
}

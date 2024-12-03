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

func getDeviceGroupResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	return client.ShowDeviceGroup(&model.ShowDeviceGroupRequest{GroupId: state.Primary.ID})
}

func TestAccDeviceGroup_basic(t *testing.T) {
	var obj model.ShowDeviceGroupResponse

	deviceName := acceptance.RandomAccResourceName()
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_device_group.test"

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
					resource.TestCheckResourceAttr(rName, "description", name),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "device_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "device_ids.0", "huaweicloud_iotda_device.test", "id"),
				),
			},
			{
				Config: testDeviceGroup_basic_update(updateName, deviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", updateName),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
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
		obj        model.ShowDeviceGroupResponse
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
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
				Config: testDeviceGroup_dynamicGroup(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttr(rName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "device_ids.#", "2"),
				),
			},
			{
				Config: testDeviceGroup_dynamicGroup_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "description", "description update"),
					resource.TestCheckResourceAttr(rName, "type", "DYNAMIC"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
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

func testDeviceGroup_basic(name, deviceName string) string {
	baseConfig := testDevice_basic(deviceName, deviceName)
	return fmt.Sprintf(`
%s

resource "huaweicloud_iotda_device_group" "test" {
  name        = "%s"
  space_id    = huaweicloud_iotda_space.test.id
  description = "%s"
  device_ids  = [huaweicloud_iotda_device.test.id]
}
`, baseConfig, name, name)
}

func testDeviceGroup_basic_update(name, deviceName string) string {
	baseConfig := testDevice_basic(deviceName, deviceName)
	return fmt.Sprintf(`
%s

resource "huaweicloud_iotda_device_group" "test" {
  name        = "%s"
  space_id    = huaweicloud_iotda_space.test.id
  description = "%s"
  device_ids  = [huaweicloud_iotda_device.test2.id]
}
`, baseConfig, name, name)
}

func testDeviceGroup_dynamicGroup(name string) string {
	baseConfig := testDevice_basic(name, name)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_group" "test" {
  name               = "%[2]s"
  space_id           = huaweicloud_iotda_space.test.id
  description        = "description test"
  type               = "DYNAMIC"
  dynamic_group_rule = "device_id = '${huaweicloud_iotda_device.test.id}' or device_id = '${huaweicloud_iotda_device.test2.id}'"
}
`, baseConfig, name)
}

func testDeviceGroup_dynamicGroup_update(name string) string {
	baseConfig := testDevice_basic(name, name)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_group" "test" {
  name               = "%[2]s"
  space_id           = huaweicloud_iotda_space.test.id
  description        = "description update"
  type               = "DYNAMIC"
  dynamic_group_rule = "device_id = '${huaweicloud_iotda_device.test.id}' or device_id = '${huaweicloud_iotda_device.test2.id}'"
}
`, baseConfig, name)
}

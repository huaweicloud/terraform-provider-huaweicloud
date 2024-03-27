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

func getDeviceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	return client.ShowDevice(&model.ShowDeviceRequest{DeviceId: state.Primary.ID})
}

func TestAccDevice_basic(t *testing.T) {
	var obj model.ShowDeviceResponse

	nodeId := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_device.test"
	childDeviceName := "huaweicloud_iotda_device.test2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDevice_basic(nodeId, nodeId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nodeId),
					resource.TestCheckResourceAttr(rName, "node_id", nodeId),
					resource.TestCheckResourceAttr(rName, "secret", "1234567890"),
					resource.TestCheckResourceAttr(rName, "secondary_secret", "test123456"),
					resource.TestCheckResourceAttr(rName, "secure_access", "true"),
					resource.TestCheckResourceAttr(rName, "description", "demo"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(rName, "auth_type", "SECRET"),
					resource.TestCheckResourceAttr(rName, "node_type", "GATEWAY"),
					resource.TestCheckResourceAttr(rName, "frozen", "false"),
					resource.TestCheckResourceAttr(childDeviceName, "name", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "node_id", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(childDeviceName, "node_type", "ENDPOINT"),
					resource.TestCheckResourceAttrPair(rName, "id", childDeviceName, "gateway_id"),
				),
			},
			{
				Config: testDevice_basic_update(updateName, nodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_id", nodeId),
					resource.TestCheckResourceAttr(rName, "fingerprint", "1234567890123456789012345678901234567890"),
					resource.TestCheckResourceAttr(rName, "secondary_fingerprint", "dc0f1016f495157344ac5f1296335cff725ef22f"),
					resource.TestCheckResourceAttr(rName, "secure_access", "false"),
					resource.TestCheckResourceAttr(rName, "description", "demo_update"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "status", "FROZEN"),
					resource.TestCheckResourceAttr(rName, "frozen", "true"),
					resource.TestCheckResourceAttr(rName, "auth_type", "CERTIFICATES"),
					resource.TestCheckResourceAttr(rName, "node_type", "GATEWAY"),
					resource.TestCheckResourceAttr(childDeviceName, "name", updateName+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "node_id", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(childDeviceName, "node_type", "ENDPOINT"),
					resource.TestCheckResourceAttrPair(rName, "id", childDeviceName, "gateway_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_disconnect",
				},
			},
		},
	})
}

func TestAccDevice_derived(t *testing.T) {
	var obj model.ShowDeviceResponse

	nodeId := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_device.test"
	childDeviceName := "huaweicloud_iotda_device.test2"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceResourceFunc,
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
				Config: testDevice_basic(nodeId, nodeId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", nodeId),
					resource.TestCheckResourceAttr(rName, "node_id", nodeId),
					resource.TestCheckResourceAttr(rName, "secret", "1234567890"),
					resource.TestCheckResourceAttr(rName, "secondary_secret", "test123456"),
					resource.TestCheckResourceAttr(rName, "secure_access", "true"),
					resource.TestCheckResourceAttr(rName, "description", "demo"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(rName, "auth_type", "SECRET"),
					resource.TestCheckResourceAttr(rName, "node_type", "GATEWAY"),
					resource.TestCheckResourceAttr(rName, "frozen", "false"),
					resource.TestCheckResourceAttr(childDeviceName, "name", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "node_id", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(childDeviceName, "node_type", "ENDPOINT"),
					resource.TestCheckResourceAttrPair(rName, "id", childDeviceName, "gateway_id"),
				),
			},
			{
				Config: testDevice_basic_update(updateName, nodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "node_id", nodeId),
					resource.TestCheckResourceAttr(rName, "fingerprint", "1234567890123456789012345678901234567890"),
					resource.TestCheckResourceAttr(rName, "secondary_fingerprint", "dc0f1016f495157344ac5f1296335cff725ef22f"),
					resource.TestCheckResourceAttr(rName, "secure_access", "false"),
					resource.TestCheckResourceAttr(rName, "description", "demo_update"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(rName, "status", "FROZEN"),
					resource.TestCheckResourceAttr(rName, "frozen", "true"),
					resource.TestCheckResourceAttr(rName, "auth_type", "CERTIFICATES"),
					resource.TestCheckResourceAttr(rName, "node_type", "GATEWAY"),
					resource.TestCheckResourceAttr(childDeviceName, "name", updateName+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "node_id", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(childDeviceName, "node_type", "ENDPOINT"),
					resource.TestCheckResourceAttrPair(rName, "id", childDeviceName, "gateway_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_disconnect",
				},
			},
		},
	})
}

func testDevice_basic(name, nodeId string) string {
	productbasic := testProduct_basic(name)
	return fmt.Sprintf(`
%s

resource "huaweicloud_iotda_device" "test" {
  node_id          = "%[3]s"
  name             = "%[2]s"
  space_id         = huaweicloud_iotda_space.test.id
  product_id       = huaweicloud_iotda_product.test.id
  secret           = "1234567890"
  secondary_secret = "test123456"
  secure_access    = true
  force_disconnect = true
  description      = "demo"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_iotda_device" "test2" {
  node_id    = "%[3]s_2"
  name       = "%[2]s_2"
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
  gateway_id = huaweicloud_iotda_device.test.id
}
`, productbasic, name, nodeId)
}

func testDevice_basic_update(name, nodeId string) string {
	productbasic := testProduct_basic(name)
	return fmt.Sprintf(`
%s

resource "huaweicloud_iotda_device" "test" {
  node_id               = "%[3]s"
  name                  = "%[2]s"
  space_id              = huaweicloud_iotda_space.test.id
  product_id            = huaweicloud_iotda_product.test.id
  fingerprint           = "1234567890123456789012345678901234567890"
  secondary_fingerprint = "dc0f1016f495157344ac5f1296335cff725ef22f"
  secure_access         = false
  description           = "demo_update"
  frozen                = true

  tags = {
    foo = "bar_update"
    key = "value"
  }
}

resource "huaweicloud_iotda_device" "test2" {
  node_id    = "%[3]s_2"
  name       = "%[2]s_2"
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
  gateway_id = huaweicloud_iotda_device.test.id
}
`, productbasic, name, nodeId)
}

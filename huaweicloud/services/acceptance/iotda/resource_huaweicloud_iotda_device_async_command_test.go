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

func getDeviceAsyncCommandFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	return client.ShowAsyncDeviceCommand(&model.ShowAsyncDeviceCommandRequest{DeviceId: state.Primary.Attributes["device_id"],
		CommandId: state.Primary.ID})
}

func TestAccDeviceAsyncCommand_basic(t *testing.T) {
	var obj model.ShowDeviceResponse

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_device_async_command.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceAsyncCommandFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceAsyncCommand_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "device_id", "huaweicloud_iotda_device.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "service_id", "huaweicloud_iotda_product.test", "services.0.id"),
					resource.TestCheckResourceAttrPair(rName, "name", "huaweicloud_iotda_product.test", "services.0.commands.0.name"),
					resource.TestCheckResourceAttr(rName, "send_strategy", "delay"),
					resource.TestCheckResourceAttr(rName, "expire_time", "80000"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
		},
	})
}

func TestAccDeviceAsyncCommand_derived(t *testing.T) {
	var obj model.ShowDeviceResponse

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_device_async_command.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDeviceAsyncCommandFunc,
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
				Config: testAccDeviceAsyncCommand_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "device_id", "huaweicloud_iotda_device.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "service_id", "huaweicloud_iotda_product.test", "services.0.id"),
					resource.TestCheckResourceAttrPair(rName, "name", "huaweicloud_iotda_product.test", "services.0.commands.0.name"),
					resource.TestCheckResourceAttr(rName, "send_strategy", "delay"),
					resource.TestCheckResourceAttr(rName, "expire_time", "80000"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
				),
			},
		},
	})
}

func testAccDeviceAsyncCommand_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iotda_space" "test" {
  name = "%[1]s"
}

resource "huaweicloud_iotda_product" "test" {
  name        = "%[1]s"
  device_type = "AI"
  protocol    = "CoAP"
  space_id    = huaweicloud_iotda_space.test.id
  data_type   = "json"
  industry    = "smart-home"

  services {
    id   = "001211985996"
    type = "001002"

    commands {
      name = "cmd-test"

      paras {
        name       = "cmd-req"
        type       = "string"
        max_length = 20
      }

      responses {
        name       = "cmd-resp"
        type       = "string"
        max_length = 20
      }
    }
  }
}

resource "huaweicloud_iotda_device" "test" {
  node_id     = "101112026"
  name        = "%[1]s"
  space_id    = huaweicloud_iotda_space.test.id
  product_id  = huaweicloud_iotda_product.test.id
  description = "demo"
}
`, name)
}

func testAccDeviceAsyncCommand_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device_async_command" "test" {
  device_id     = huaweicloud_iotda_device.test.id
  service_id    = huaweicloud_iotda_product.test.services.0.id
  name          = huaweicloud_iotda_product.test.services.0.commands.0.name
  send_strategy = "delay"
  expire_time   = "80000"

  paras = {
    (huaweicloud_iotda_product.test.services.0.commands.0.paras.0.name) = "tf-acc"
  }
}
`, testAccDeviceAsyncCommand_base(name))
}

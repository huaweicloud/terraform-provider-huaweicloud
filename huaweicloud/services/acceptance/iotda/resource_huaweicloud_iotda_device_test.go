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

func getDeviceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "iotda"
		httpUrl = "v5/iot/{project_id}/devices/{device_id}"
	)

	isDerived := WithDerivedAuth()
	client, err := conf.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{device_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccDevice_basic(t *testing.T) {
	var obj interface{}

	nodeId := acceptance.RandomAccResourceName()
	name := acceptance.RandomAccResourceName()
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
				Config: testDevice_basic(name, nodeId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
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
					resource.TestCheckResourceAttr(rName, "extension_info.tf", "terraform"),
					resource.TestCheckResourceAttr(rName, "shadow.#", "2"),
					resource.TestCheckResourceAttr(childDeviceName, "name", name+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "node_id", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(childDeviceName, "node_type", "ENDPOINT"),
					resource.TestCheckResourceAttrPair(rName, "id", childDeviceName, "gateway_id"),
				),
			},
			{
				Config: testDevice_updateSecret(name, nodeId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "secret", "1234567890123"),
					resource.TestCheckResourceAttr(rName, "secondary_secret", "test123456123"),
				),
			},
			{
				Config: testDevice_update(name, nodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
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
					resource.TestCheckResourceAttr(rName, "extension_info.tf", "update"),
					resource.TestCheckResourceAttr(rName, "extension_info.test", "acc"),
					resource.TestCheckResourceAttr(childDeviceName, "name", name+"_2_update"),
					resource.TestCheckResourceAttr(childDeviceName, "node_id", nodeId+"_2"),
					resource.TestCheckResourceAttr(childDeviceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(childDeviceName, "node_type", "ENDPOINT"),
					resource.TestCheckResourceAttrPair(rName, "id", childDeviceName, "gateway_id"),
				),
			},
			{
				Config: testDevice_unfreeze(name, nodeId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "fingerprint", "1234567890123456789012345678901234567123"),
					resource.TestCheckResourceAttr(rName, "secondary_fingerprint", "dc0f1016f495157344ac5f1296335cff725ef123"),
					resource.TestCheckResourceAttr(rName, "frozen", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_disconnect", "extension_info", "shadow",
				},
			},
		},
	})
}

func testAccDevice_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_product" "test" {
  name        = "%[2]s-pro"
  device_type = "Thermometer"
  protocol    = "MQTT"
  space_id    = huaweicloud_iotda_space.test.id
  data_type   = "json"

  services {
    id   = "temp_1"
    type = "temperature_a"

    properties {
      name       = "demo_1"
      type       = "string"
      max_length = 256
      method     = "RW"
    }
  }

  services {
    id   = "temp_2"
    type = "temperature_b"

    properties {
      name   = "demo_2"
      type   = "int"
      method = "RW"
    }

    properties {
      name       = "demo_3"
      type       = "string"
      max_length = 256
      method     = "RW"
    }

    commands {
      name = "cmd_1"

      paras {
        name        = "cmd_p_1"
        type        = "int"
        required    = false
        description = "desc"
        min         = "3"
        max         = 33
      }

      responses {
        name     = "cmd_r_1"
        type     = "int"
        required = false
        min      = "1"
        max      = "22"
      }
    }
  }
}
`, testSpace_basic(name), name)
}

func testDevice_basic(name, nodeId string) string {
	return fmt.Sprintf(`
%[1]s

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

  extension_info = {
    tf = "terraform"
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[0].id

    desired = {
      (huaweicloud_iotda_product.test.services[0].properties[0].name) = "test"
    }
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[1].id

    desired = {
      (huaweicloud_iotda_product.test.services[1].properties[0].name) = "acc"
    }
  }

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
`, testAccDevice_base(name), name, nodeId)
}

func testDevice_updateSecret(name, nodeId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device" "test" {
  node_id          = "%[3]s"
  name             = "%[2]s"
  space_id         = huaweicloud_iotda_space.test.id
  product_id       = huaweicloud_iotda_product.test.id
  secret           = "1234567890123"
  secondary_secret = "test123456123"
  secure_access    = true
  force_disconnect = true
  description      = "demo"

  extension_info = {
    tf = "terraform"
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[0].id

    desired = {
      (huaweicloud_iotda_product.test.services[0].properties[0].name) = "test"
    }
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[1].id

    desired = {
      (huaweicloud_iotda_product.test.services[1].properties[0].name) = "acc"
    }
  }

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
`, testAccDevice_base(name), name, nodeId)
}

func testDevice_update(name, nodeId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device" "test" {
  node_id               = "%[3]s"
  name                  = "%[2]s_update"
  space_id              = huaweicloud_iotda_space.test.id
  product_id            = huaweicloud_iotda_product.test.id
  fingerprint           = "1234567890123456789012345678901234567890"
  secondary_fingerprint = "dc0f1016f495157344ac5f1296335cff725ef22f"
  secure_access         = false
  description           = "demo_update"
  frozen                = true

  extension_info = {
    tf   = "update"
    test = "acc"
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[0].id

    desired = {}
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[1].id

    desired = {
      (huaweicloud_iotda_product.test.services[1].properties[0].name) = "update"
      (huaweicloud_iotda_product.test.services[1].properties[1].name) = "retest"
    }
  }

  tags = {
    foo = "bar_update"
    key = "value"
  }
}

resource "huaweicloud_iotda_device" "test2" {
  node_id    = "%[3]s_2"
  name       = "%[2]s_2_update"
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
  gateway_id = huaweicloud_iotda_device.test.id
}
`, testAccDevice_base(name), name, nodeId)
}

func testDevice_unfreeze(name, nodeId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device" "test" {
  node_id               = "%[3]s"
  name                  = "%[2]s_update"
  space_id              = huaweicloud_iotda_space.test.id
  product_id            = huaweicloud_iotda_product.test.id
  fingerprint           = "1234567890123456789012345678901234567123"
  secondary_fingerprint = "dc0f1016f495157344ac5f1296335cff725ef123"
  secure_access         = false
  description           = "demo_update"
  frozen                = false

  extension_info = {
    tf   = "update"
    test = "acc"
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[0].id

    desired = {}
  }

  shadow {
    service_id = huaweicloud_iotda_product.test.services[1].id

    desired = {
      (huaweicloud_iotda_product.test.services[1].properties[0].name) = "update"
      (huaweicloud_iotda_product.test.services[1].properties[1].name) = "retest"
    }
  }

  tags = {
    foo = "bar_update"
    key = "value"
  }
}

resource "huaweicloud_iotda_device" "test2" {
  node_id    = "%[3]s_2"
  name       = "%[2]s_2_update"
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
  gateway_id = huaweicloud_iotda_device.test.id
}
`, testAccDevice_base(name), name, nodeId)
}

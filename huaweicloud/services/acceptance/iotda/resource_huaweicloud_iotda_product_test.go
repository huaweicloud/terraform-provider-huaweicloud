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

func getProductResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	return client.ShowProduct(&model.ShowProductRequest{ProductId: state.Primary.ID})
}

func TestAccProduct_basic(t *testing.T) {
	var obj model.ShowProductResponse

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_product.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProductResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testProduct_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "device_type", "test"),
					resource.TestCheckResourceAttr(rName, "protocol", "MQTT"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "data_type", "json"),
					resource.TestCheckResourceAttr(rName, "manufacturer_name", "demo_manufacturer_name"),
					resource.TestCheckResourceAttr(rName, "industry", "demo_industry"),
					resource.TestCheckResourceAttr(rName, "services.0.id", "service_1"),
					resource.TestCheckResourceAttr(rName, "services.0.type", "serv_type"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.#", "4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.name", "p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.min", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.max", "666"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.method", "RW"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.name", "p_2"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.type", "string"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.max_length", "20"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.method", "R"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.method", "W"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.3.type", "decimal"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.3.min", "3.1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.3.max", "666.99"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.name", "cmd_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.name", "cmd_p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.min", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.max", "33"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.name", "cmd_r_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.min", "1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.max", "22"),
				),
			},
			{
				Config: testProduct_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.#", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.name", "p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.description", "desc_update"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.min", "4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.max", "5"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.method", "RW"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.name", "p_3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.type", "string"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.max_length", "20"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.method", "W"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.name", "p_4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.type", "decimal"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.min", "3.2"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.max", "777.99"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.name", "cmd_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.name", "cmd_p_2"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.min", "5"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.max", "33"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.name", "cmd_r_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.min", "2"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.max", "33"),
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

func TestAccProduct_derived(t *testing.T) {
	var obj model.ShowProductResponse

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_iotda_product.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getProductResourceFunc,
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
				Config: testProduct_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "device_type", "test"),
					resource.TestCheckResourceAttr(rName, "protocol", "MQTT"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "huaweicloud_iotda_space.test", "id"),
					resource.TestCheckResourceAttr(rName, "data_type", "json"),
					resource.TestCheckResourceAttr(rName, "manufacturer_name", "demo_manufacturer_name"),
					resource.TestCheckResourceAttr(rName, "industry", "demo_industry"),
					resource.TestCheckResourceAttr(rName, "services.0.id", "service_1"),
					resource.TestCheckResourceAttr(rName, "services.0.type", "serv_type"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.#", "4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.name", "p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.min", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.max", "666"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.method", "RW"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.name", "p_2"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.type", "string"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.max_length", "20"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.method", "R"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.method", "W"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.3.type", "decimal"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.3.min", "3.1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.3.max", "666.99"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.name", "cmd_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.name", "cmd_p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.min", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.max", "33"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.name", "cmd_r_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.min", "1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.max", "22"),
				),
			},
			{
				Config: testProduct_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.#", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.name", "p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.description", "desc_update"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.min", "4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.max", "5"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.method", "RW"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.name", "p_3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.type", "string"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.max_length", "20"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.1.method", "W"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.name", "p_4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.type", "decimal"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.min", "3.2"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.2.max", "777.99"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.name", "cmd_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.name", "cmd_p_2"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.min", "5"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.max", "33"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.name", "cmd_r_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.min", "2"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.max", "33"),
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

func testProduct_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_space" "test" {
  name = "%[2]s"
}

resource "huaweicloud_iotda_product" "test" {
  name              = "%[3]s"
  device_type       = "test"
  protocol          = "MQTT"
  space_id          = huaweicloud_iotda_space.test.id
  data_type         = "json"
  manufacturer_name = "demo_manufacturer_name"
  industry          = "demo_industry"

  services {
    id   = "service_1"
    type = "serv_type"

    properties {
      name        = "p_1"
      type        = "int"
      min         = "3"
      max         = "666"
      description = "desc"
      method      = "RW"
    }

    properties {
      name       = "p_2"
      type       = "string"
      max_length = 20
      enum_list  = ["1", "E"]
      method     = "R"
    }

    properties {
      name       = "p_3"
      type       = "string"
      method     = "W"
      max_length = 200
    }

    properties {
      name   = "p_4"
      type   = "decimal"
      method = "W"
      min    = "3.1"
      max    = "666.99"
    }

    commands {
      name = "cmd_1"

      paras {
        name        = "cmd_p_1"
        type        = "int"
        description = "desc"
        min         = "3"
        max         = 33
      }

      responses {
        name = "cmd_r_1"
        type = "int"
        min  = "1"
        max  = "22"
      }
    }
  }
}
`, buildIoTDAEndpoint(), name, name)
}

func testProduct_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_space" "test" {
  name = "%[2]s"
}

resource "huaweicloud_iotda_product" "test" {
  name              = "%[3]s"
  device_type       = "test"
  protocol          = "MQTT"
  space_id          = huaweicloud_iotda_space.test.id
  data_type         = "json"
  manufacturer_name = "demo_manufacturer_name"
  industry          = "demo_industry"

  services {
    id   = "service_1"
    type = "serv_type"

    properties {
      name        = "p_1"
      type        = "int"
      min         = "4"
      max         = "5"
      description = "desc_update"
      method      = "RW"
    }

    properties {
      name       = "p_3"
      type       = "string"
      method     = "W"
      max_length = 20
    }

    properties {
      name   = "p_4"
      type   = "decimal"
      method = "W"
      min    = "3.2"
      max    = "777.99"
    }

    commands {
      name = "cmd_1"

      paras {
        name        = "cmd_p_2"
        type        = "int"
        description = "desc"
        min         = "5"
        max         = 33
      }

      responses {
        name = "cmd_r_1"
        type = "int"
        min  = "2"
        max  = "33"
      }
    }
  }
}
`, buildIoTDAEndpoint(), name, name)
}

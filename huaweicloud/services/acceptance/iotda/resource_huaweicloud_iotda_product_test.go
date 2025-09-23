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

func getProductResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "iotda"
		httpUrl = "v5/iot/{project_id}/products/{product_id}"
	)

	isDerived := WithDerivedAuth()
	client, err := conf.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{product_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccProduct_basic(t *testing.T) {
	var obj interface{}

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
					resource.TestCheckResourceAttr(rName, "services.0.option", "Master"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.#", "4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.name", "p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.required", "false"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.min", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.max", "666"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.method", "RW"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.default_value", "{\"foo\":\"bar\"}"),
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
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.required", "false"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.min", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.max", "33"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.name", "cmd_r_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.required", "false"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.min", "1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.max", "22"),
				),
			},
			{
				Config: testProduct_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "services.0.option", "Optional"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.#", "3"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.name", "p_1"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.required", "true"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.description", "desc_update"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.min", "4"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.max", "5"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.method", "RW"),
					resource.TestCheckResourceAttr(rName, "services.0.properties.0.default_value", "{\"tf\":\"test\"}"),
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
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.required", "true"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.description", "desc"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.min", "5"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.paras.0.max", "33"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.name", "cmd_r_1"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.type", "int"),
					resource.TestCheckResourceAttr(rName, "services.0.commands.0.responses.0.required", "true"),
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

resource "huaweicloud_iotda_product" "test" {
  name              = "%[2]s"
  device_type       = "test"
  protocol          = "MQTT"
  space_id          = huaweicloud_iotda_space.test.id
  data_type         = "json"
  manufacturer_name = "demo_manufacturer_name"
  industry          = "demo_industry"

  services {
    id     = "service_1"
    type   = "serv_type"
    option = "Master"

    properties {
      name          = "p_1"
      type          = "int"
      required      = false
      min           = "3"
      max           = "666"
      description   = "desc"
      method        = "RW"
      default_value = "{\"foo\":\"bar\"}"
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

func testProduct_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_product" "test" {
  name              = "%[2]s_update"
  device_type       = "test"
  protocol          = "MQTT"
  space_id          = huaweicloud_iotda_space.test.id
  data_type         = "json"
  manufacturer_name = "demo_manufacturer_name"
  industry          = "demo_industry"

  services {
    id     = "service_1"
    type   = "serv_type"
    option = "Optional"

    properties {
      name          = "p_1"
      type          = "int"
      required      = true
      min           = "4"
      max           = "5"
      description   = "desc_update"
      method        = "RW"
      default_value = "{\"tf\":\"test\"}"
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
        required    = true
        description = "desc"
        min         = "5"
        max         = 33
      }

      responses {
        name     = "cmd_r_1"
        type     = "int"
        required = true
        min      = "2"
        max      = "33"
      }
    }
  }
}
`, testSpace_basic(name), name)
}

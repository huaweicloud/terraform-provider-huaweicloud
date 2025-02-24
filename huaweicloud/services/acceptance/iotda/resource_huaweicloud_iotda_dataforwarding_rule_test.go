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

func getDataForwardingRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "iotda"
		httpUrl = "v5/iot/{project_id}/routing-rule/rules/{rule_id}"
	)

	isDerived := WithDerivedAuth()
	client, err := conf.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

// The forwarding target of **DMS_KAFKA_FORWARDING** type requires the use of a public IP address, which may pose a
// security port risk, so it will not be tested temporarily.
func TestAccDataForwardingRule_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_iotda_dataforwarding_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataForwardingRuleResourceFunc,
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
				Config: testDataForwardingRule_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "HTTP_FORWARDING"),
					resource.TestCheckResourceAttr(rName, "targets.0.http_forwarding.0.url", "http://www.exampletest.com"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
				),
			},
			{
				Config: testDataForwardingRule_basicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "HTTP_FORWARDING"),
					resource.TestCheckResourceAttr(rName, "targets.0.http_forwarding.0.url", "http://www.example.com"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
				),
			},
			{
				Config: testDataForwardingRule_DIS(name, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "DIS_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.dis_forwarding.0.stream_id",
						"huaweicloud_dis_stream.test.0", "id"),
				),
			},
			{
				Config: testDataForwardingRule_DISUpdate(name, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "DIS_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.dis_forwarding.0.stream_id",
						"huaweicloud_dis_stream.test.1", "id"),
				),
			},
			{
				Config: testDataForwardingRule_OBS(name, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "OBS_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.obs_forwarding.0.bucket",
						"huaweicloud_obs_bucket.test.0", "bucket"),
				),
			},
			{
				Config: testDataForwardingRule_OBSUpdate(name, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "OBS_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.obs_forwarding.0.bucket",
						"huaweicloud_obs_bucket.test.1", "bucket"),
				),
			},
			{
				Config: testDataForwardingRule_FGS(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "FUNCTIONGRAPH_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.fgs_forwarding.0.func_urn",
						"huaweicloud_fgs_function.test.0", "urn"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.fgs_forwarding.0.func_name",
						"huaweicloud_fgs_function.test.0", "name"),
				),
			},
			{
				Config: testDataForwardingRule_FGSUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "FUNCTIONGRAPH_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.fgs_forwarding.0.func_urn",
						"huaweicloud_fgs_function.test.1", "urn"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.fgs_forwarding.0.func_name",
						"huaweicloud_fgs_function.test.1", "name"),
				),
			},
			{
				Config: testDataForwardingRule_AMQP(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "AMQP_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.amqp_forwarding.0.queue_name",
						"huaweicloud_iotda_amqp.test.0", "name"),
				),
			},
			{
				Config: testDataForwardingRule_AMQPUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "AMQP_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.amqp_forwarding.0.queue_name",
						"huaweicloud_iotda_amqp.test.1", "name"),
				),
			},
			{
				Config: testDataForwardingRule_multipleTargets(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "true"),
					resource.TestCheckResourceAttr(rName, "targets.#", "2"),
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

func testDataForwardingRule_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[2]s"
  trigger = "product:delete"
  enabled = true
  
  targets {
    type = "HTTP_FORWARDING"

    http_forwarding {
      url = "http://www.exampletest.com"
    }
  }
}
`, buildIoTDAEndpoint(), name)
}

func testDataForwardingRule_basicUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[2]s_update"
  trigger = "product:delete"
  enabled = false
  
  targets {
    type = "HTTP_FORWARDING"

    http_forwarding {
      url = "http://www.example.com"
    }
  }
}
`, buildIoTDAEndpoint(), name)
}

func testDataForwardingRule_DIS(name, region string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dis_stream" "test" {
  count = 2

  stream_name     = format("%[2]s_%%d", count.index)
  partition_count = 1
}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[3]s"
  trigger = "product:delete"
  enabled = false

  targets {
    type = "DIS_FORWARDING"

    dis_forwarding {
      region    = "%[4]s"
      stream_id = huaweicloud_dis_stream.test[0].id
    }
  }
}
`, buildIoTDAEndpoint(), name, name, region)
}

func testDataForwardingRule_DISUpdate(name, region string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dis_stream" "test" {
  count = 2

  stream_name     = format("%[2]s_%%d", count.index)
  partition_count = 1
}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[3]s"
  trigger = "product:delete"
  enabled = true

  targets {
    type = "DIS_FORWARDING"

    dis_forwarding {
      region    = "%[4]s"
      stream_id = huaweicloud_dis_stream.test[1].id
    }
  }
}
`, buildIoTDAEndpoint(), name, name, region)
}

func testDataForwardingRule_OBS(name, region string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "test" {
  count = 2

  bucket        = format("tf-test-bucket-%%d", count.index)
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[2]s"
  trigger = "product:delete"
  enabled = false

  targets {
    type = "OBS_FORWARDING"

    obs_forwarding {
      region = "%[3]s"
      bucket = huaweicloud_obs_bucket.test[0].bucket
    }
  }
}
`, buildIoTDAEndpoint(), name, region)
}

func testDataForwardingRule_OBSUpdate(name, region string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket" "test" {
  count = 2

  bucket        = format("tf-test-bucket-%%d", count.index)
  storage_class = "STANDARD"
  acl           = "private"
}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[2]s"
  trigger = "product:delete"
  enabled = true

  targets {
    type = "OBS_FORWARDING"

    obs_forwarding {
      region = "%[3]s"
      bucket = huaweicloud_obs_bucket.test[1].bucket
    }
  }
}
`, buildIoTDAEndpoint(), name, region)
}

func testDataForwardingRule_FGS(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[3]s"
  trigger = "product:delete"
  enabled = false

  targets {
    type = "FUNCTIONGRAPH_FORWARDING"

    fgs_forwarding {
      func_urn  = huaweicloud_fgs_function.test[0].urn
      func_name = huaweicloud_fgs_function.test[0].name
    }
  }
}
`, buildIoTDAEndpoint(), testDataForwardingRule_FGS_base(), name)
}

func testDataForwardingRule_FGSUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s
%[2]s

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[3]s"
  trigger = "product:delete"
  enabled = true

  targets {
    type = "FUNCTIONGRAPH_FORWARDING"

    fgs_forwarding {
      func_urn  = huaweicloud_fgs_function.test[1].urn
      func_name = huaweicloud_fgs_function.test[1].name
    }
  }
}
`, buildIoTDAEndpoint(), testDataForwardingRule_FGS_base(), name)
}

func testDataForwardingRule_FGS_base() string {
	rName := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
variable "request_resp_print_script_content" {
  default = <<EOT
exports.handler = async (event, context) => {
    const result =
    {
        'repsonse_code': 200,
        'headers':
        {
            'Content-Type': 'application/json'
        },
        'isBase64Encoded': false,
        'body': JSON.stringify(event)
    }
    return result
}
EOT
}

resource "huaweicloud_fgs_function" "test" {
  count = 2

  name        = format("%s-%%d", count.index)
  app         = "default"
  description = "function test"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 3
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = base64encode(var.request_resp_print_script_content)

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testDataForwardingRule_AMQP(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_amqp" "test" {
  count = 2

  name = format("%[2]s_%%d", count.index)
}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[2]s"
  trigger = "product:delete"
  enabled = false

  targets {
    type = "AMQP_FORWARDING"

    amqp_forwarding {
      queue_name = huaweicloud_iotda_amqp.test[0].name
    }
  }
}
`, buildIoTDAEndpoint(), name)
}

func testDataForwardingRule_AMQPUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_amqp" "test" {
  count = 2

  name = format("%[2]s_%%d", count.index)
}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[2]s"
  trigger = "product:delete"
  enabled = true

  targets {
    type = "AMQP_FORWARDING"

    amqp_forwarding {
      queue_name = huaweicloud_iotda_amqp.test[1].name
    }
  }
}
`, buildIoTDAEndpoint(), name)
}

func testDataForwardingRule_multipleTargets(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_amqp" "test" {
  count = 2

  name = format("%[2]s_%%d", count.index)
}

resource "huaweicloud_iotda_dataforwarding_rule" "test" {
  name    = "%[2]s"
  trigger = "product:delete"
  enabled = true

  targets {
    type = "AMQP_FORWARDING"

    amqp_forwarding {
      queue_name = huaweicloud_iotda_amqp.test[1].name
    }
  }

  targets {
    type = "HTTP_FORWARDING"

    http_forwarding {
      url = "http://www.exampletest.com"
    }
  }
}
`, buildIoTDAEndpoint(), name)
}

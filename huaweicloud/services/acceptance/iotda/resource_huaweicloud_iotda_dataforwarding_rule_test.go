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

func getDataForwardingRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	return client.ShowRoutingRule(&model.ShowRoutingRuleRequest{RuleId: state.Primary.ID})
}

func TestAccDataForwardingRule_basic(t *testing.T) {
	var obj model.ShowRoutingRuleResponse

	name := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	rName := "huaweicloud_iotda_dataforwarding_rule.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataForwardingRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
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
				Config: testDataForwardingRule_dis(updateName, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "DIS_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.dis_forwarding.0.stream_id",
						"huaweicloud_dis_stream.test", "id"),
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

func TestAccDataForwardingRule_derived(t *testing.T) {
	var obj model.ShowRoutingRuleResponse

	name := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
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
				Config: testDataForwardingRule_dis(updateName, acceptance.HW_REGION_NAME),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "trigger", "product:delete"),
					resource.TestCheckResourceAttr(rName, "enabled", "false"),
					resource.TestCheckResourceAttr(rName, "targets.#", "1"),
					resource.TestCheckResourceAttr(rName, "targets.0.type", "DIS_FORWARDING"),
					resource.TestCheckResourceAttrSet(rName, "targets.0.id"),
					resource.TestCheckResourceAttrPair(rName, "targets.0.dis_forwarding.0.stream_id",
						"huaweicloud_dis_stream.test", "id"),
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

func testDataForwardingRule_dis(name, region string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dis_stream" "test" {
  stream_name     = "%[2]s"
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
      stream_id = huaweicloud_dis_stream.test.id
    }
  }
}
`, buildIoTDAEndpoint(), name, name, region)
}

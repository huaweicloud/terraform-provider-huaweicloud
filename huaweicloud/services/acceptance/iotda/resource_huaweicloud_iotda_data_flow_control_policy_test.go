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

func getDataFlowControlPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When the resource does not exist, query API will return `404` error code.
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccDataFlowControlPolicy_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_iotda_data_flow_control_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataFlowControlPolicyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource only supports standard and enterprise version IoTDA instances.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataFlowControlPolicy_userScope(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description_test"),
					resource.TestCheckResourceAttr(rName, "scope", "USER"),
					resource.TestCheckResourceAttr(rName, "limit", "1000"),
				),
			},
			{
				Config: testDataFlowControlPolicy_userScope_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "description_update"),
					resource.TestCheckResourceAttr(rName, "scope", "USER"),
					resource.TestCheckResourceAttr(rName, "limit", "200"),
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

func TestAccDataFlowControlPolicy_channelScope(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_iotda_data_flow_control_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataFlowControlPolicyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource only supports standard and enterprise version IoTDA instances.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataFlowControlPolicy_channelScope(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description_test"),
					resource.TestCheckResourceAttr(rName, "scope", "CHANNEL"),
					resource.TestCheckResourceAttr(rName, "scope_value", "DMS_KAFKA_FORWARDING"),
					resource.TestCheckResourceAttr(rName, "limit", "1000"),
				),
			},
			{
				Config: testDataFlowControlPolicy_channelScope_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "description_update"),
					resource.TestCheckResourceAttr(rName, "scope", "CHANNEL"),
					resource.TestCheckResourceAttr(rName, "scope_value", "DMS_KAFKA_FORWARDING"),
					resource.TestCheckResourceAttr(rName, "limit", "200"),
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

func TestAccDataFlowControlPolicy_ruleScope(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_iotda_data_flow_control_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataFlowControlPolicyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource only supports standard and enterprise version IoTDA instances.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataFlowControlPolicy_ruleScope(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description_test"),
					resource.TestCheckResourceAttr(rName, "scope", "RULE"),
					resource.TestCheckResourceAttr(rName, "limit", "1000"),
					resource.TestCheckResourceAttrPair(rName, "scope_value",
						"huaweicloud_iotda_dataforwarding_rule.test", "id"),
				),
			},
			{
				Config: testDataFlowControlPolicy_ruleScope_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "description_update"),
					resource.TestCheckResourceAttr(rName, "scope", "RULE"),
					resource.TestCheckResourceAttr(rName, "limit", "200"),
					resource.TestCheckResourceAttrPair(rName, "scope_value",
						"huaweicloud_iotda_dataforwarding_rule.test", "id"),
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

func TestAccDataFlowControlPolicy_actionScope(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_iotda_data_flow_control_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataFlowControlPolicyResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource only supports standard and enterprise version IoTDA instances.
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataFlowControlPolicy_actionScope(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description_test"),
					resource.TestCheckResourceAttr(rName, "scope", "ACTION"),
					resource.TestCheckResourceAttr(rName, "limit", "1000"),
					resource.TestCheckResourceAttrPair(rName, "scope_value",
						"huaweicloud_iotda_dataforwarding_rule.test", "targets.0.id"),
				),
			},
			{
				Config: testDataFlowControlPolicy_actionScope_update(name + "_update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "description_update"),
					resource.TestCheckResourceAttr(rName, "scope", "ACTION"),
					resource.TestCheckResourceAttr(rName, "limit", "200"),
					resource.TestCheckResourceAttrPair(rName, "scope_value",
						"huaweicloud_iotda_dataforwarding_rule.test", "targets.0.id"),
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

func testDataFlowControlPolicy_userScope(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name        = "%[2]s"
  description = "description_test"
  scope       = "USER"
}
`, buildIoTDAEndpoint(), name)
}

func testDataFlowControlPolicy_userScope_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name        = "%[2]s"
  description = "description_update"
  scope       = "USER"
  limit       = 200
}
`, buildIoTDAEndpoint(), name)
}

func testDataFlowControlPolicy_channelScope(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name        = "%[2]s"
  description = "description_test"
  scope       = "CHANNEL"
  scope_value = "DMS_KAFKA_FORWARDING"
}
`, buildIoTDAEndpoint(), name)
}

func testDataFlowControlPolicy_channelScope_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name        = "%[2]s"
  description = "description_update"
  scope       = "CHANNEL"
  scope_value = "DMS_KAFKA_FORWARDING"
  limit       = 200
}
`, buildIoTDAEndpoint(), name)
}

func testDataFlowControlPolicy_ruleScope(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name        = "%[2]s"
  description = "description_test"
  scope       = "RULE"
  scope_value = huaweicloud_iotda_dataforwarding_rule.test.id
}
`, testDataForwardingRule_basic(name), name)
}

func testDataFlowControlPolicy_ruleScope_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  name        = "%[2]s_update"
  description = "description_update"
  scope       = "RULE"
  scope_value = huaweicloud_iotda_dataforwarding_rule.test.id
  limit       = 200
}
`, testDataForwardingRule_basic(name), name)
}

func testDataFlowControlPolicy_actionScope(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  depends_on = [
    huaweicloud_iotda_dataforwarding_rule.test
  ]

  name        = "%[2]s"
  description = "description_test"
  scope       = "ACTION"
  scope_value = [for target in huaweicloud_iotda_dataforwarding_rule.test.targets : target.id][0]
}
`, testDataForwardingRule_basic(name), name)
}

func testDataFlowControlPolicy_actionScope_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_flow_control_policy" "test" {
  depends_on = [
    huaweicloud_iotda_dataforwarding_rule.test
  ]

  name        = "%[2]s"
  description = "description_update"
  scope       = "ACTION"
  scope_value = [for target in huaweicloud_iotda_dataforwarding_rule.test.targets : target.id][0]
  limit       = 200
}
`, testDataForwardingRule_basic(name), name)
}

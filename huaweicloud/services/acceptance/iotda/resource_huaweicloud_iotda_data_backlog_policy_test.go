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

func getDataBacklogPolicyResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "iotda"
		httpUrl = "v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}"
	)

	isDerived := WithDerivedAuth()
	client, err := conf.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", state.Primary.ID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func TestAccDataBacklogPolicy_basic(t *testing.T) {
	var (
		obj   interface{}
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_iotda_data_backlog_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDataBacklogPolicyResourceFunc,
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
				Config: testDataBacklogPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "description_test"),
					// The `backlog_size` field default value is `1073741823`.
					resource.TestCheckResourceAttr(rName, "backlog_size", "1073741823"),
					// The `backlog_time` field default value is `86399`.
					resource.TestCheckResourceAttr(rName, "backlog_time", "86399"),
				),
			},
			{
				Config: testDataBacklogPolicy_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "description_update"),
					resource.TestCheckResourceAttr(rName, "backlog_size", "100"),
					resource.TestCheckResourceAttr(rName, "backlog_time", "100"),
				),
			},
			{
				Config: testDataBacklogPolicy_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"_update"),
					resource.TestCheckResourceAttr(rName, "description", "description_update"),
					resource.TestCheckResourceAttr(rName, "backlog_size", "0"),
					resource.TestCheckResourceAttr(rName, "backlog_time", "0"),
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

func testDataBacklogPolicy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_backlog_policy" "test" {
  # Before creating a data backlog policy, it is necessary to ensure that there is a data forwarding rule in place.
  depends_on = [
    huaweicloud_iotda_dataforwarding_rule.test
  ]

  name        = "%[2]s"
  description = "description_test"
}
`, testDataForwardingRule_basic(name), name)
}

func testDataBacklogPolicy_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_backlog_policy" "test" {
  depends_on = [
    huaweicloud_iotda_dataforwarding_rule.test
  ]

  name         = "%[2]s_update"
  description  = "description_update"
  backlog_size = "100"
  backlog_time = "100"
}
`, testDataForwardingRule_basic(name), name)
}

func testDataBacklogPolicy_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_data_backlog_policy" "test" {
  depends_on = [
    huaweicloud_iotda_dataforwarding_rule.test
  ]

  name         = "%[2]s_update"
  description  = "description_update"
  backlog_size = "0"
  backlog_time = "0"
}
`, testDataForwardingRule_basic(name), name)
}

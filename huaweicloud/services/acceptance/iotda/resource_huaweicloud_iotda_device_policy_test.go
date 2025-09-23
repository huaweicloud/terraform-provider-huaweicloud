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

func getDevicePolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-policies/{policy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", state.Primary.ID)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IoTDA device policy: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccDevicePolicy_basic(t *testing.T) {
	var (
		devicePolicyObj interface{}
		name            = acceptance.RandomAccResourceName()
		updateName      = name + "_update"
		rName           = "huaweicloud_iotda_device_policy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&devicePolicyObj,
		getDevicePolicyResourceFunc,
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
				Config: testDevicePolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "policy_name", name),
					resource.TestCheckResourceAttr(rName, "statement.#", "1"),
					resource.TestCheckResourceAttr(rName, "statement.0.effect", "ALLOW"),
					resource.TestCheckResourceAttr(rName, "statement.0.actions.#", "1"),
					resource.TestCheckResourceAttr(rName, "statement.0.resources.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "data.huaweicloud_iotda_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
				),
			},
			{
				Config: testDevicePolicy_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "policy_name", updateName),
					resource.TestCheckResourceAttr(rName, "statement.#", "2"),
					resource.TestCheckResourceAttr(rName, "statement.0.effect", "DENY"),
					resource.TestCheckResourceAttr(rName, "statement.0.actions.#", "2"),
					resource.TestCheckResourceAttr(rName, "statement.0.resources.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "space_id", "data.huaweicloud_iotda_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
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

func testDevicePolicy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_spaces" "test" {
  is_default = true
}

resource "huaweicloud_iotda_device_policy" "test" {
  policy_name = "%[2]s"

  statement {
    effect    = "ALLOW"
    actions   = ["iotda:devices:publish"]
    resources = ["topic:/v1/test/hello"]
  }

  space_id = data.huaweicloud_iotda_spaces.test.spaces[0].id
}
`, buildIoTDAEndpoint(), name)
}

func testDevicePolicy_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_spaces" "test" {
  is_default = true
}

resource "huaweicloud_iotda_device_policy" "test" {
  policy_name = "%[2]s"

  statement {
    effect    = "DENY"
    actions   = ["iotda:devices:subscribe", "iotda:devices:publish"]
    resources = ["topic:/v1/test/hello1", "topic:/v1/test/hello2"]
  }

  statement {
    effect    = "ALLOW"
    actions   = ["iotda:devices:subscribe"]
    resources = ["topic:/v1/test/hello1", "topic:/v1/test/hello2"]
  }

  space_id = data.huaweicloud_iotda_spaces.test.spaces[0].id
}
`, buildIoTDAEndpoint(), name)
}

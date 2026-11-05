package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getTmsAssociatedResourceRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("tms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating TMS client: %s", err)
	}

	return getTmsAssociatedResourceRule(client, state.Primary.Attributes["setting_name"],
		state.Primary.Attributes["region_id"])
}

func getTmsAssociatedResourceRule(client *golangsdk.ServiceClient, settingName, regionId string) (interface{}, error) {
	listPath := client.Endpoint + "v2/tms/associated-resources/rules"
	listPath += fmt.Sprintf("?setting_name=%s&region_id=%s", settingName, regionId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving TMS associated resource rule: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	rule := utils.PathSearch("rules[0]", listRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method: "GET",
				URL:    "/v2/tms/associated-resources/rules",
			},
		}
	}

	return rule, nil
}

func TestAccTmsAssociatedResourceRule_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_tms_associated_resource_rule.test"
	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTmsAssociatedResourceRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// The CheckDestroy verifies the DELETE /v2/tms/associated-resources/rules/{setting_name} behavior.
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// The step covers the create behavior (POST /v2/tms/associated-resources/rules/batch-create)
				// and the read behavior (GET /v2/tms/associated-resources/rules).
				Config: testAccTmsAssociatedResourceRule_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "setting_name"),
					resource.TestCheckResourceAttr(rName, "region_id", "cn-north-4"),
					resource.TestCheckResourceAttr(rName, "tag_keys.#", "1"),
					resource.TestCheckResourceAttr(rName, "tag_keys.0", "key1"),
					resource.TestCheckResourceAttr(rName, "existing_resource_status", "enable"),
					resource.TestCheckResourceAttr(rName, "auto_delete_status", "disable"),
				),
			},
			{
				// The step covers the update behavior (PUT /v2/tms/associated-resources/rules/batch-update),
				// changing `tag_keys`, `existing_resource_status` and `auto_delete_status`.
				Config: testAccTmsAssociatedResourceRule_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "setting_name"),
					resource.TestCheckResourceAttr(rName, "region_id", "cn-north-4"),
					resource.TestCheckResourceAttr(rName, "tag_keys.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_keys.0", "key1"),
					resource.TestCheckResourceAttr(rName, "tag_keys.1", "key2"),
					resource.TestCheckResourceAttr(rName, "existing_resource_status", "disable"),
					resource.TestCheckResourceAttr(rName, "auto_delete_status", "enable"),
				),
			},
			{
				// The step covers the import behavior, the import ID must be `<setting_name>/<region_id>`.
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// The tag keys returned by the API may not be in the same order as the configuration.
					"tag_keys",
				},
			},
		},
	})
}

func testAccTmsAssociatedResourceRule_basic() string {
	return `
data "huaweicloud_tms_associated_resource_settings" "test" {}

resource "huaweicloud_tms_associated_resource_status" "test" {
  status = "enable"
}

resource "huaweicloud_tms_associated_resource_rule" "test" {
  setting_name             = data.huaweicloud_tms_associated_resource_settings.test.settings[0].setting_name
  region_id                = "cn-north-4"
  tag_keys                 = ["key1"]
  existing_resource_status = "enable"
  auto_delete_status       = "disable"

  depends_on = [huaweicloud_tms_associated_resource_status.test]
}
`
}

func testAccTmsAssociatedResourceRule_update() string {
	return `
data "huaweicloud_tms_associated_resource_settings" "test" {}

resource "huaweicloud_tms_associated_resource_status" "test" {
  status = "enable"
}

resource "huaweicloud_tms_associated_resource_rule" "test" {
  setting_name             = data.huaweicloud_tms_associated_resource_settings.test.settings[0].setting_name
  region_id                = "cn-north-4"
  tag_keys                 = ["key1", "key2"]
  existing_resource_status = "disable"
  auto_delete_status       = "enable"

  depends_on = [huaweicloud_tms_associated_resource_status.test]
}
`
}

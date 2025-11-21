package identitycenter

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

func getIdentityCenterMfaManagementSettingResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getHttpUrl = "v1/instances/{instance_id}/mfa-devices/management-settings"
		product    = "identitycenter"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center client: %s", err)
	}

	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center mfa management settings: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccIdentityCenterMfaManagementSetting_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_mfa_management_setting.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterMfaManagementSettingResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testMfaManagementSetting_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "data.huaweicloud_identitycenter_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "user_permission", "READ_ACTIONS"),
				),
			},
			{
				Config: testMfaManagementSetting_basic_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "data.huaweicloud_identitycenter_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "user_permission", "ALL_ACTIONS"),
				),
			},
		},
	})
}

const testMfaManagementSetting_basic = `
data "huaweicloud_identitycenter_instance" "test" {}

resource "huaweicloud_identitycenter_mfa_management_setting" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.test.id
  user_permission   = "READ_ACTIONS"
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`

const testMfaManagementSetting_basic_update = `
data "huaweicloud_identitycenter_instance" "test" {}

resource "huaweicloud_identitycenter_mfa_management_setting" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.test.id
  user_permission   = "ALL_ACTIONS"
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
`

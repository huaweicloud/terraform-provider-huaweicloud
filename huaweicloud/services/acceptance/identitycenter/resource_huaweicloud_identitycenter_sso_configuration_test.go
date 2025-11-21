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

func getIdentityCenterSsoConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		getHttpUrl = "v1/instances/{instance_id}/sso-configuration"
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
		return nil, fmt.Errorf("error retrieving sso configuration: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccIdentityCenterSsoConfiguration_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_sso_configuration.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterSsoConfigurationResourceFunc,
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
				Config: testSsoConfiguration_basic,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "mfa_mode", "CONTEXT_AWARE"),
					resource.TestCheckResourceAttrSet(rName, "allowed_mfa_types.#"),
					resource.TestCheckResourceAttr(rName, "no_mfa_signin_behavior", "ALLOWED_WITH_ENROLLMENT"),
					resource.TestCheckResourceAttr(rName, "no_password_signin_behavior", "BLOCKED"),
					resource.TestCheckResourceAttr(rName, "max_authentication_age", "PT12H"),
				),
			},
			{
				Config: testSsoConfiguration_basic_update,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "mfa_mode", "DISABLED"),
					resource.TestCheckResourceAttrSet(rName, "allowed_mfa_types.#"),
					resource.TestCheckResourceAttr(rName, "no_mfa_signin_behavior", "EMAIL_OTP"),
					resource.TestCheckResourceAttr(rName, "no_password_signin_behavior", "EMAIL_OTP"),
					resource.TestCheckResourceAttr(rName, "max_authentication_age", "PT4H"),
				),
			},
		},
	})
}

const testSsoConfiguration_basic = `
data "huaweicloud_identitycenter_instance" "test" {}

resource "huaweicloud_identitycenter_sso_configuration" "test" {
  instance_id                 = data.huaweicloud_identitycenter_instance.test.id
  mfa_mode                    = "CONTEXT_AWARE"
  allowed_mfa_types           = ["TOTP"]
  no_mfa_signin_behavior      = "ALLOWED_WITH_ENROLLMENT"
  no_password_signin_behavior = "BLOCKED"
  max_authentication_age      = "PT12H"
  configuration_type          = "APP_AUTHENTICATION_CONFIGURATION"
}
`

const testSsoConfiguration_basic_update = `
data "huaweicloud_identitycenter_instance" "test" {}

resource "huaweicloud_identitycenter_sso_configuration" "test" {
  instance_id                 = data.huaweicloud_identitycenter_instance.test.id
  mfa_mode                    = "DISABLED"
  allowed_mfa_types           = ["TOTP","WEBAUTHN_SECURITY_KEY"]
  no_mfa_signin_behavior      = "EMAIL_OTP"
  no_password_signin_behavior = "EMAIL_OTP"
  max_authentication_age      = "PT4H"
  configuration_type          = "APP_AUTHENTICATION_CONFIGURATION"
}
`

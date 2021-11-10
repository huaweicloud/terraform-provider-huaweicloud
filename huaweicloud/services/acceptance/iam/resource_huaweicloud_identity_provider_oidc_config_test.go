package iam

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/oidcconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProviderOidcConfigFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMNoVersionClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud IAM without version: %s", err)
	}
	return oidcconfig.Get(client, state.Primary.ID)
}

func TestAccIdentityProviderOidcConfig_basic(t *testing.T) {
	var oidcConfig oidcconfig.OpenIDConnectConfig
	var name = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_provider_oidc_config.config"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&oidcConfig,
		getProviderOidcConfigFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderOidcConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "client_id", "client_id_example"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "program"),
					//resource.TestCheckResourceAttr(resourceName, "scopes.#", "0"),
				),
			},
			{
				Config: testAccIdentityProviderOidcConfig_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "client_id", "client_id_demo"),
					resource.TestCheckResourceAttr(resourceName, "access_type", "program_console"),
					resource.TestCheckResourceAttr(resourceName, "scopes.1", "email"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityProviderOidcConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "oidc"
}

resource "huaweicloud_identity_provider_oidc_config" "config" {
  provider_id            = huaweicloud_identity_provider.provider_1.id
  access_type            = "program"
  provider_url           = "https://accounts.example.com"
  client_id              = "client_id_example"
  signing_key            = jsonencode(
  {
    keys = [
      {
        alg = "RS256"
        e   = "AQAB"
        kid = "d05ef20c4512645vv1..."
        kty = "RSA"
        n   = "cws_cnjiwsbvweolwn_-vnl..."
        use = "sig"
      },
    ]
  }
  )
}
`, name)
}

func testAccIdentityProviderOidcConfig_basic_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "oidc"
}

resource "huaweicloud_identity_provider_oidc_config" "config" {
  provider_id            = huaweicloud_identity_provider.provider_1.id
  access_type            = "program_console"
  provider_url           = "https://accounts.example.com"
  client_id              = "client_id_demo"
  authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
  scopes                 = ["openid", "email"]
  signing_key            = jsonencode(
  {
    keys = [
      {
        alg = "RS256"
        e   = "AQAB"
        kid = "d05ef20c4512645vv1..."
        kty = "RSA"
        n   = "cws_cnjiwsbvweolwn_-vnl..."
        use = "sig"
      },
    ]
  }
  )
}
`, name)
}

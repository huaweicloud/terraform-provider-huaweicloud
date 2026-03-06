package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getV3ProviderResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMNoVersionClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client without version: %s", err)
	}
	return providers.Get(client, state.Primary.ID)
}

func TestAccV3Provider_basic(t *testing.T) {
	var (
		obj interface{}

		protocolSaml   = "huaweicloud_identity_provider.protocol_saml"
		rcProtocolSaml = acceptance.InitResourceCheck(protocolSaml, &obj, getV3ProviderResourceFunc)

		protocolOidc   = "huaweicloud_identity_provider.protocol_oidc"
		rcProtocolOidc = acceptance.InitResourceCheck(protocolOidc, &obj, getV3ProviderResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcProtocolSaml.CheckResourceDestroy(),
			rcProtocolOidc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Provider_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcProtocolSaml.CheckResourceExists(),
					resource.TestCheckResourceAttr(protocolSaml, "name", name+"_saml"),
					resource.TestCheckResourceAttr(protocolSaml, "protocol", "saml"),
					rcProtocolOidc.CheckResourceExists(),
					resource.TestCheckResourceAttr(protocolOidc, "name", name+"_oidc"),
					resource.TestCheckResourceAttr(protocolOidc, "protocol", "oidc"),
					resource.TestCheckResourceAttr(protocolOidc, "access_config.0.access_type", "program_console"),
					resource.TestCheckResourceAttr(protocolOidc, "access_config.0.client_id", name+"_step1"),
				),
			},
			{
				Config: testAccV3Provider_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcProtocolSaml.CheckResourceExists(),
					resource.TestCheckResourceAttr(protocolSaml, "name", name+"_saml"),
					resource.TestCheckResourceAttr(protocolSaml, "protocol", "saml"),
					resource.TestCheckResourceAttr(protocolSaml, "status", "false"),
					rcProtocolOidc.CheckResourceExists(),
					resource.TestCheckResourceAttr(protocolOidc, "name", name+"_oidc"),
					resource.TestCheckResourceAttr(protocolOidc, "protocol", "oidc"),
					resource.TestCheckResourceAttr(protocolOidc, "access_config.0.access_type", "program_console"),
					resource.TestCheckResourceAttr(protocolOidc, "access_config.0.client_id", name+"_step2"),
				),
			},
			{
				Config: testAccV3Provider_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rcProtocolOidc.CheckResourceExists(),
					resource.TestCheckResourceAttr(protocolOidc, "name", name+"_oidc"),
					resource.TestCheckResourceAttr(protocolOidc, "protocol", "oidc"),
					resource.TestCheckResourceAttr(protocolOidc, "status", "false"),
					resource.TestCheckResourceAttr(protocolOidc, "access_config.0.access_type", "program"),
					resource.TestCheckResourceAttr(protocolOidc, "access_config.0.client_id", name+"_step3"),
				),
			},
			{
				ResourceName:      protocolSaml,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      protocolOidc,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV3Provider_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "protocol_saml" {
  name     = "%[1]s_saml"
  protocol = "saml"
}

resource "huaweicloud_identity_provider" "protocol_oidc" {
  name        = "%[1]s_oidc"
  protocol    = "oidc"
  description = "unit test"

  access_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "%[1]s_step1"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
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
}
`, name)
}

func testAccV3Provider_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "protocol_saml" {
  name     = "%[1]s_saml"
  protocol = "saml"
  status   = false
}

resource "huaweicloud_identity_provider" "protocol_oidc" {
  name        = "%[1]s_oidc"
  protocol    = "oidc"
  description = "unit test"

  access_config {
    access_type            = "program_console"
    provider_url           = "https://accounts.example.com"
    client_id              = "%[1]s_step2"
    authorization_endpoint = "https://accounts.example.com/o/oauth2/v2/auth"
    scopes                 = ["openid"]
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
}
`, name)
}

func testAccV3Provider_basic_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "protocol_saml" {
  name     = "%[1]s_saml"
  protocol = "saml"
  status   = false
}

resource "huaweicloud_identity_provider" "protocol_oidc" {
  name        = "%[1]s_oidc"
  protocol    = "oidc"
  status      = false
  description = "unit test"

  access_config {
    access_type  = "program"
    provider_url = "https://accounts.example.com"
    client_id    = "%[1]s_step3"
    signing_key  = jsonencode(
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
}
`, name)
}

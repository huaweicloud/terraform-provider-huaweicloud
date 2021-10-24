package iam

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProviderResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMNoVersionClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud IAM without version: %s", err)
	}
	return providers.Get(client, state.Primary.ID)
}

func TestAccIdentityProvider_basic(t *testing.T) {
	var provider providers.Provider
	var name = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_provider.provider_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&provider,
		getProviderResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProvider_saml(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "protocol", "saml"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.remote.#", "1"),
				),
			},
			{
				Config: testAccIdentityProvider_saml_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.1.remote.0.condition", "any_one_of"),
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

func TestAccIdentityProvider_oidc(t *testing.T) {
	var provider providers.Provider
	var name = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_provider.provider_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&provider,
		getProviderResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProvider_oidc(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccIdentityProvider_oidc_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "scopes.#", "3"),
				),
			},
		},
	})
}

func testAccIdentityProvider_saml(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "saml"
  status   = "enabled"

  conversion_rules {
    local {
      username = "federateduser"
    }
    remote {
      attribute = "federatedgroup"
    }
  }
}
`, name)
}

func testAccIdentityProvider_saml_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider_1" {
  name     = "%s"
  protocol = "saml"
  status   = "disabled"

  conversion_rules {
    local {
      username = "federateduser"
    }
    local {
      username = "Tom"
    }
    remote {
      attribute = "federatedgroup"
    }
    remote {
      attribute = "Tom"
    }
  }

  conversion_rules {
    local {
      username = "Jams"
    }
    remote {
      attribute = "username"
      condition = "any_one_of"
      value     = ["Tom", "Jerry"]
    }
  }
}
`, name)
}

func testAccIdentityProvider_oidc(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider_1" {
  name        = "%s"
  protocol    = "oidc"
  status      = "enabled"
  description = "unit test"

  access_type            = "program_console"
  provider_url           = "https://accounts.example.com"
  client_id              = "client_id_example"
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

  conversion_rules {
    local {
      username = "federateduser"
    }
    remote {
      attribute = "federatedgroup"
    }
  }
}
`, name)
}

func testAccIdentityProvider_oidc_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider_1" {
  name        = "%s"
  protocol    = "oidc"
  status      = "disabled"
  description = "unit test"

  access_type            = "program_console"
  provider_url           = "https://new.accounts.example.com"
  client_id              = "client_id_example_new"
  authorization_endpoint = "https://new.accounts.example.com/o/oauth2/v2/auth"
  scopes                 = ["openid", "email", "aba"]
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

  conversion_rules {
    local {
      username = "federateduser"
    }
    local {
      username = "Tom"
    }
    remote {
      attribute = "federatedgroup"
    }
    remote {
      attribute = "Tom"
    }
  }

  conversion_rules {
    local {
      username = "Jams"
    }
    remote {
      attribute = "username"
      condition = "any_one_of"
      value     = ["Tom", "Jerry"]
    }
  }
}
`, name)
}

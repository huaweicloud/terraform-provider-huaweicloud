package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/mappings"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProviderConversionFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMNoVersionClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client without version: %s", err)
	}
	providerID := state.Primary.Attributes["provider_id"]
	conversionID := "mapping_" + providerID
	return mappings.Get(client, conversionID)
}

func TestAccProviderConversion_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_provider_conversion.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getProviderConversionFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConversion_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.0.username", "Tom"),
				),
			},
			{
				Config: testAccProviderConversion_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.0.local.0.username", "Tom"),
					resource.TestCheckResourceAttr(resourceName, "conversion_rules.1.remote.0.value.#", "2"),
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

func testAccProviderConversion_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%[1]s"
  protocol = "oidc"
}

resource "huaweicloud_identity_provider_conversion" "test" {
  provider_id = huaweicloud_identity_provider.test.id

  conversion_rules {
    local {
      username = "Tom"
    }
    remote {
      attribute = "Tom"
    }
  }
}
`, name)
}

func testAccProviderConversion_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%[1]s"
  protocol = "oidc"
}

resource "huaweicloud_identity_provider_conversion" "test" {
  provider_id = huaweicloud_identity_provider.test.id

  conversion_rules {
    local {
      username = "Tom"
    }
    local {
      username = "federateduser"
    }

    remote {
      attribute = "Tom"
    }
    remote {
      attribute = "federatedgroup"
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

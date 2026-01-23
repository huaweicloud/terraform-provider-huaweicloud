package iam

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

func getProviderMappingFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.IAMNoVersionClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client without version: %s", err)
	}
	providerID := state.Primary.Attributes["provider_id"]
	mappingID := "mapping_" + providerID

	getMappingHttpUrl := "v3/OS-FEDERATION/mappings/{id}"
	getMappingPath := client.Endpoint + getMappingHttpUrl
	getMappingPath = strings.ReplaceAll(getMappingPath, "{id}", mappingID)
	getMappingOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getMappingResp, err := client.Request("GET", getMappingPath, &getMappingOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getMappingResp)
}

func TestAccProviderMapping_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_provider_mapping.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getProviderMappingFunc)

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
				Config: testAccProviderMapping_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// nothing to check
				),
			},
			{
				Config: testAccProviderMapping_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// nothing to check
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

func testAccProviderMapping_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%[1]s"
  protocol = "oidc"
}

resource "huaweicloud_identity_provider_mapping" "test" {
  provider_id = huaweicloud_identity_provider.test.id

  mapping_rules = <<RULES
    [
      {
        "local": [
          {
            "user": {
              "name": "{0}"
            }
          },
          {
            "group": {
              "name": "admin"
            }
          }
        ],
        "remote": [
          {
            "type": "UserName"
          },
          {
            "type": "Groups",
            "any_one_of": [
              ".*@mail.com$"
            ],
            "regex": true
          }
        ]
      }
    ]
  RULES
}
`, name)
}

func testAccProviderMapping_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name     = "%[1]s"
  protocol = "oidc"
}

resource "huaweicloud_identity_provider_mapping" "test" {
  provider_id = huaweicloud_identity_provider.test.id

  mapping_rules = <<RULES
    [
      {
        "local": [
          {
            "user": {
              "name": "{0} {1}"
            }
          },
          {
            "group": {
              "name": "{2}"
            }
          }
        ],
        "remote": [
          {
            "type": "FirstName"
          },
          {
            "type": "LastName"
          },
          {
            "type": "Group"
          }
        ]
      }
    ]
  RULES
}
`, name)
}

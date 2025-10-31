package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/federatedauth/providers"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iam/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getProviderProtocolFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	iamV3Client, err := c.HcIamV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err), err
	}
	request := &model.KeystoneShowProtocolRequest{
		IdpId:      state.Primary.Attributes["provider_id"],
		ProtocolId: state.Primary.Attributes["protocol_id"],
	}
	response, err := iamV3Client.KeystoneShowProtocol(request)
	if err != nil {
		return fmt.Errorf("KeystoneShowProtocol error : %s", err), err
	}
	return response.Protocol, nil
}

func generateTestMappingID(providerID string) string {
	return "test_" + providerID
}

func TestAccIdentityProviderProtocol_basic(t *testing.T) {
	var provider providers.Provider
	var randomName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_provider_protocol.protocol"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&provider,
		getProviderProtocolFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProviderProtocolCreate(randomName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "provider_id", randomName),
					resource.TestCheckResourceAttr(resourceName, "protocol_id", "saml"),
					resource.TestCheckResourceAttr(resourceName, "mapping_id", "mapping_"+generateTestMappingID(randomName)),
					resource.TestCheckResourceAttr(resourceName, "links.#", "1"),
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

// huaweicloud_identity_provider创建一个saml协议的idp，mapping_id为mapping_{provider_id}
// huaweicloud_identity_provider_conversion创建一个新的mapping，mapping_id为mapping_generateTestMappingID(name)
// huaweicloud_identity_provider_protocol更新协议，修改mapping为mapping_generateTestMappingID(name)
func testAccIdentityProviderProtocolCreate(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider1" {
  name     = "%s"
  protocol = "saml"
}

resource "huaweicloud_identity_provider_conversion" "conversion" {
  provider_id = "%s"

  conversion_rules {
    local {
      username = "Tom"
    }
    remote {
      attribute = "Tom"
    }
  }
}

resource "huaweicloud_identity_provider_protocol" "protocol" {
  provider_id = huaweicloud_identity_provider.provider1.id
  protocol_id = "saml"
  mapping_id  = huaweicloud_identity_provider_conversion.conversion.id
}
`, name, generateTestMappingID(name))
}

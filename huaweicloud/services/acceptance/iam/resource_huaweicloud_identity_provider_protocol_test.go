package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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
		return nil, err
	}
	return response.Protocol, nil
}

func TestAccProviderProtocol_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_provider_protocol.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getProviderProtocolFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderProtocol_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "provider_id", name),
					resource.TestCheckResourceAttr(resourceName, "protocol_id", "saml"),
					resource.TestCheckResourceAttrPair(resourceName, "mapping_id", "huaweicloud_identity_provider_conversion.test", "id"),
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

func testAccProviderProtocol_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "test" {
  name = "%[1]s"

  lifecycle {
    ignore_changes = [conversion_rules]
  }
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

resource "huaweicloud_identity_provider_protocol" "test" {
  provider_id = huaweicloud_identity_provider.test.id
  protocol_id = "saml"
  mapping_id  = huaweicloud_identity_provider_conversion.test.id
}
`, name)
}

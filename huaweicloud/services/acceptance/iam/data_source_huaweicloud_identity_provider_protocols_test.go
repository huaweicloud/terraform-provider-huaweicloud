package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIamIdentityProviderProtocols_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_provider_protocols.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceIamIdentityProviderProtocols(acceptance.RandomAccResourceName()),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.#", "1"),
				),
			},
			{
				Config: testDataSourceIamIdentityProviderProtocolsWithProtocolId(acceptance.RandomAccResourceName()),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "protocols.0.id", "saml"),
				),
			},
		},
	})
}

func testDataSourceIamIdentityProviderProtocols(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider1" {
  name     = "%s"
  protocol = "oidc"
}

data "huaweicloud_identity_provider_protocols" "test" {
  provider_id = huaweicloud_identity_provider.provider1.id
}
`, name)
}

func testDataSourceIamIdentityProviderProtocolsWithProtocolId(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_provider" "provider1" {
  name     = "%s"
  protocol = "saml"
}

data "huaweicloud_identity_provider_protocols" "test" {
  provider_id = huaweicloud_identity_provider.provider1.id
  protocol_id = "saml"
}
`, name)
}

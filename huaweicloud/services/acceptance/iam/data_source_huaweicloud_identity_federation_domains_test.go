package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityFederationDomains_basic(t *testing.T) {
	idpId := "YourIdpId"
	protocolId := "YourProtocolId"
	idToken := "YourIdToken"
	dataSourceName := "data.huaweicloud_identity_federation_domains.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestDataSourceIdentityFederationDomains(idpId, protocolId, idToken),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domains.0.enabled"),
					resource.TestCheckResourceAttr(dataSourceName, "domains.0.description", ""),
				),
			},
		},
	})
}

func testTestDataSourceIdentityFederationDomains(idpId, protocolId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
 idp_id      = "%[1]s"
 protocol_id = "%[2]s"
 id_token    = "%[3]s"
}

data "huaweicloud_identity_federation_domains" "test" {
 federation_token = huaweicloud_identity_unscoped_token_with_id_token.test.token
}
`, idpId, protocolId, idToken)
}

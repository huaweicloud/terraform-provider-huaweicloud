package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFederationDomains_basic(t *testing.T) {
	var (
		idpId      = "YourIdpId"
		protocolId = "YourProtocolId"
		idToken    = "YourIdToken"

		all = "data.huaweicloud_identity_federation_domains.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFederationDomains_basic(idpId, protocolId, idToken),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "domains.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "domains.0.id"),
					resource.TestCheckResourceAttrSet(all, "domains.0.name"),
					resource.TestCheckResourceAttrSet(all, "domains.0.enabled"),
					resource.TestCheckResourceAttr(all, "domains.0.description", ""),
				),
			},
		},
	})
}

func testAccDataFederationDomains_basic(idpId, protocolId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
  idp_id      = "%[1]s"
  protocol_id = "%[2]s"
  id_token    = "%[3]s"
}

data "huaweicloud_identity_federation_domains" "all" {
  federation_token = huaweicloud_identity_unscoped_token_with_id_token.test.token
}
`, idpId, protocolId, idToken)
}

package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFederationProjects_basic(t *testing.T) {
	var (
		idpId      = "YourIdpId"
		protocolId = "YourProtocolId"
		idToken    = "YourIdToken"

		all = "data.huaweicloud_identity_federation_projects.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFederationProjects_basic(idpId, protocolId, idToken),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "projects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "projects.0.id"),
					resource.TestCheckResourceAttrSet(all, "projects.0.name"),
					resource.TestCheckResourceAttrSet(all, "projects.0.enabled"),
				),
			},
		},
	})
}

func testAccDataFederationProjects_basic(idpId, protocolId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
  idp_id      = "%[1]s"
  protocol_id = "%[2]s"
  id_token    = "%[3]s"
}

data "huaweicloud_identity_federation_projects" "all" {
  federation_token = huaweicloud_identity_unscoped_token_with_id_token.test.token
}
`, idpId, protocolId, idToken)
}

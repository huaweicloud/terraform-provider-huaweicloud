package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityFederationProjects_basic(t *testing.T) {
	idpId := "YourIdpId"
	protocolId := "YourProtocolId"
	idToken := "YourIdToken"
	dataSourceName := "data.huaweicloud_identity_federation_projects.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestDataSourceIdentityFederationProjects(idpId, protocolId, idToken),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "projects.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "projects.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "projects.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "projects.0.enabled"),
				),
			},
		},
	})
}

func testTestDataSourceIdentityFederationProjects(idpId, protocolId, idToken string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
idp_id      = "%[1]s"
protocol_id = "%[2]s"
id_token    = "%[3]s"
}

data "huaweicloud_identity_federation_projects" "test" {
federation_token = huaweicloud_identity_unscoped_token_with_id_token.test.token
}
`, idpId, protocolId, idToken)
}

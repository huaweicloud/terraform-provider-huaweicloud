package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceApplications_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_applications.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	uuid, _ := uuid.GenerateUUID()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckIdentityCenterAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceApplications_basic(uuid),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.application_urn"),
					resource.TestCheckResourceAttr(dataSource,
						"applications.0.application_provider_urn", "IdentityCenter:::applicationProvider:custom-saml"),
					resource.TestCheckResourceAttr(dataSource, "applications.0.assignment_config.0.assignment_required", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.created_date"),
					resource.TestCheckResourceAttr(dataSource, "applications.0.description", "create"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.instance_urn"),
					resource.TestCheckResourceAttr(dataSource, "applications.0.name", "create"),
					resource.TestCheckResourceAttrSet(dataSource, "applications.0.application_account"),
					resource.TestCheckResourceAttr(dataSource, "applications.0.portal_options.0.visible", "true"),
					resource.TestCheckResourceAttr(dataSource, "applications.0.portal_options.0.visibility", "ENABLED"),
					resource.TestCheckResourceAttr(dataSource, "applications.0.portal_options.0.sign_in_options.0.origin", "IDENTITY_CENTER"),
					resource.TestCheckResourceAttr(dataSource, "applications.0.portal_options.0.sign_in_options.0.application_url", ""),
				),
			},
		},
	})
}

func testDataSourceApplications_basic(uuid string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_applications" "test"{
  depends_on  = [huaweicloud_identitycenter_application_instance.test]
  instance_id = data.huaweicloud_identitycenter_instance.test.id
}
`, testApplicationInstance_basic(uuid))
}

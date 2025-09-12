package identitycenter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterCatalogApplications_basic(t *testing.T) {
	rName := "data.huaweicloud_identitycenter_catalog_applications.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterCatalogApplications_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "id"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.application_id"),
					resource.TestCheckResourceAttr(rName, "applications.0.application_type", ""),
					resource.TestCheckResourceAttrSet(rName, "applications.0.display.0.display_name"),
					resource.TestCheckResourceAttrSet(rName, "applications.0.display.0.description"),
					resource.TestCheckResourceAttr(rName, "applications.0.display.0.icon", ""),
				),
			},
		},
	})
}

const testAccDatasourceIdentityCenterCatalogApplications_basic = `
data "huaweicloud_identitycenter_catalog_applications" "test" {}
`

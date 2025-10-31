package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityProjectQuotaDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_project_quota.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityProjectQuotaDataSource,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "resources.0.type", "project"),
				),
			},
		},
	})
}

const testAccIdentityProjectQuotaDataSource = `
data "huaweicloud_identity_projects" "test" {
  name = "MOS"
}

data "huaweicloud_identity_project_quota" "test" {
  project_id = data.huaweicloud_identity_projects.test.projects[0].id
}
`

package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityCatalogDataSource_basic(t *testing.T) {
	userName := acceptance.RandomAccResourceName()
	initPassword := acceptance.RandomPassword()
	resourceName := "data.huaweicloud_identity_catalog.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCatalogDataSource_basic(userName, initPassword),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "catalog.#"),
				),
			},
		},
	})
}

func testAccIdentityCatalogDataSource_basic(userName, initPassword string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_catalog" "test" {
  project_token = huaweicloud_identity_user_token.test.token
}
`, testAccIdentityUserToken_project(userName, initPassword))
}

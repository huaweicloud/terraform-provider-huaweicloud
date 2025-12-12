package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityEnterpriseProjectUsers_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_enterprise_project_users.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityEnterpriseProjectUsers(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "users.#"),
				),
			},
		},
	})
}

func testAccIdentityEnterpriseProjectUsers() string {
	return fmt.Sprintf(`
data "huaweicloud_identity_enterprise_project_users" "test" {
 	enterprise_project_id = "%s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID)
}

package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectUsers_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_enterprise_project_users.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectUsers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "users.0.id"),
					resource.TestCheckResourceAttrSet(all, "users.0.name"),
				),
			},
		},
	})
}

func testAccDataEnterpriseProjectUsers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_identity_enterprise_project_users" "all" {
  enterprise_project_id = "%[1]s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

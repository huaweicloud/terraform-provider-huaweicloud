package organizations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTagPolicyServices_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_organizations_tag_policy_services.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTagPolicyServices_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "services.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "services.0.service_name"),
					resource.TestMatchResourceAttr(all, "services.0.resource_types.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "services.0.resource_types.0"),
					resource.TestCheckResourceAttrSet(all, "services.0.support_all"),
				),
			},
		},
	})
}

const testAccDataTagPolicyServices_basic = `
data "huaweicloud_organizations_tag_policy_services" "test" {}
`

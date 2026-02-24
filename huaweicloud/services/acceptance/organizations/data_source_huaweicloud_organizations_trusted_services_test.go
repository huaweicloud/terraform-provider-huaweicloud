package organizations

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTrustedServices_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_organizations_trusted_services.test"
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
				Config: testAccDataTrustedServices_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "trusted_services.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "trusted_services.0.service_principal"),
					resource.TestCheckResourceAttrSet(all, "trusted_services.0.enabled_at"),
				),
			},
		},
	})
}

const testAccDataTrustedServices_basic = `
resource "huaweicloud_organizations_trusted_service" "test" {
  service = "service.SecMaster"
}

data "huaweicloud_organizations_trusted_services" "test" {
  depends_on = [huaweicloud_organizations_trusted_service.test]
}
`

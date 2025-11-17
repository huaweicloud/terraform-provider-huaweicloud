package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDomainOwnerVerification_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cdn_domain_owner_verification.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDomainOwnerVerification_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttrSet(rName, "dns_verify_type"),
					resource.TestCheckResourceAttrSet(rName, "dns_verify_name"),
					resource.TestCheckResourceAttrSet(rName, "file_verify_url"),
					resource.TestCheckResourceAttrSet(rName, "verify_domain_name"),
					resource.TestCheckResourceAttrSet(rName, "file_verify_filename"),
					resource.TestCheckResourceAttrSet(rName, "verify_content"),
					resource.TestCheckResourceAttrSet(rName, "file_verify_domains.#"),
				),
			},
		},
	})
}

func testAccDataDomainOwnerVerification_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cdn_domain_owner_verification" "test" {
  domain_name = "%[1]s"
}
`, acceptance.HW_CDN_DOMAIN_NAME)
}

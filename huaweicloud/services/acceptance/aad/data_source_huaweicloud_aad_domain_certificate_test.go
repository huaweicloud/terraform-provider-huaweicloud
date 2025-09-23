package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAadDomainCertificate_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_domain_certificate.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an AAD protected domain ID before running this test case.
			acceptance.TestAccPreCheckAadDomainID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAadDomainCertificate_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "domain_name"),
				),
			},
		},
	})
}

func testDataSourceAadDomainCertificate_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_aad_domain_certificate" "test" {
  domain_id = "%s"
}
`, acceptance.HW_AAD_DOMAIN_ID)
}

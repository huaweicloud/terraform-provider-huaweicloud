package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCertificates_basic(t *testing.T) {
	rName := "data.huaweicloud_ccm_certificates.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMCertificateName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCertificates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.name"),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.domain"),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.expire_time"),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.status"),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.domain_count"),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.wildcard_count"),
					resource.TestCheckResourceAttrSet(rName, "certificates.0.enterprise_project_id"),
				),
			},
		},
	})
}

func testAccDatasourceCertificates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ccm_certificates" "test" {
  name = "%s"
}
`, acceptance.HW_CCM_CERTIFICATE_NAME)
}

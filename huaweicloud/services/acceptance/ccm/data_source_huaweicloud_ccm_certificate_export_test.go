package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCertificateExport_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ccm_certificate_export.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the certificate ID of the completed certificate application
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCertificateExport_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "entire_certificate"),
					resource.TestCheckResourceAttrSet(dataSource, "certificate"),
					resource.TestCheckResourceAttrSet(dataSource, "certificate_chain"),
					resource.TestCheckResourceAttrSet(dataSource, "private_key"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCertificateExport_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ccm_certificate_export" "test" {
  certificate_id = "%s"
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID)
}

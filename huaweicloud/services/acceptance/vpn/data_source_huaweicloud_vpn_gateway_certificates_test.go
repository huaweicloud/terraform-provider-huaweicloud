package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnGatewayCertificates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_gateway_certificates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnGatewayCertificates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.vgw_id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.issuer"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.certificate_serial_number"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.certificate_subject"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.certificate_expire_time"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.certificate_chain_serial_number"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.certificate_chain_subject"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.certificate_chain_expire_time"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.enc_certificate_serial_number"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.enc_certificate_subject"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.enc_certificate_expire_time"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceVpnGatewayCertificates_basic() string {
	return `
data "huaweicloud_vpn_gateway_certificates" "test" {}
`
}

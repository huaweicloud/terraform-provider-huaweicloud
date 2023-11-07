package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccDataSourceWafCertificateV1_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_certificate.cert_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCertificateListV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafCertDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "expire_status", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "expiration"),
				),
			},
		},
	})
}

func TestAccDataSourceWafCertificateV1_withEpsID(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_certificate.cert_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafCertificateListV1_conf_withEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafCertDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "expire_status", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "expiration"),
				),
			},
		},
	})
}

func testAccCheckWafCertDataSourceID(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmtp.Errorf("Can't find waf data source: %s ", r)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("The Waf Certificate data source ID not set ")
		}
		return nil
	}
}

func testAccWafCertificateListV1_conf(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_certificate" "cert_1" {
  name          = huaweicloud_waf_certificate.certificate_1.name
  expire_status = 1

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name))
}

func testAccWafCertificateListV1_conf_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_certificate" "cert_1" {
  name                  = huaweicloud_waf_certificate.certificate_1.name
  enterprise_project_id = "%s"
  expire_status         = 1

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), epsID)
}

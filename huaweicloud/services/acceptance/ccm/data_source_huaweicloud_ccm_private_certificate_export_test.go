package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCcmPrivateCertificateExport_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_ccm_private_certificate_export.basiccert"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCcmPrivateCertificateExport_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "private_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate_chain"),
				),
			},
		},
	})
}

func TestAccCcmPrivateCertificateExport_SM2(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSourceSM2Name := "data.huaweicloud_ccm_private_certificate_export.sm2cert"
	dc := acceptance.InitDataSourceCheck(dataSourceSM2Name)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCcmPrivateCertificateSM2Export_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceSM2Name, "private_key"),
					resource.TestCheckResourceAttrSet(dataSourceSM2Name, "certificate"),
					resource.TestCheckResourceAttrSet(dataSourceSM2Name, "certificate_chain"),
					resource.TestCheckResourceAttrSet(dataSourceSM2Name, "enc_certificate"),
					resource.TestCheckResourceAttrSet(dataSourceSM2Name, "enc_sm2_enveloped_key"),
					resource.TestCheckResourceAttrSet(dataSourceSM2Name, "signed_and_enveloped_data"),
				),
			},
		},
	})
}
func TestAccCcmPrivateCertificateExport_compressed(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	tomcatDataSourceName := "data.huaweicloud_ccm_private_certificate_export.tomcatcert"
	iisDataSourceName := "data.huaweicloud_ccm_private_certificate_export.iiscert"
	dcTomcat := acceptance.InitDataSourceCheck(tomcatDataSourceName)
	dcIIS := acceptance.InitDataSourceCheck(iisDataSourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCcmPrivateCertificateExport_TOMCAT(rName),
				Check: resource.ComposeTestCheckFunc(
					dcTomcat.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(tomcatDataSourceName, "keystore_pass"),
					resource.TestCheckResourceAttrSet(tomcatDataSourceName, "server_jks"),
				),
			},
			{
				Config: testAccCcmPrivateCertificateExport_IIS(rName),
				Check: resource.ComposeTestCheckFunc(
					dcIIS.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(iisDataSourceName, "keystore_pass"),
					resource.TestCheckResourceAttrSet(iisDataSourceName, "server_pfx"),
				),
			},
		},
	})
}

// lintignore:AT004
func testAccCcmPrivateCertificateExport_basic(commonName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_certificate_export" "basiccert" {
    type           = "other"
    certificate_id = huaweicloud_ccm_private_certificate.test.id
}`, tesCmdbCertificate_basic(commonName))
}

// lintignore:AT004
func testAccCcmPrivateCertificateExport_TOMCAT(commonName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_certificate_export" "tomcatcert" {
    type           = "TOMCAT"
    certificate_id = huaweicloud_ccm_private_certificate.test.id
}`, tesCmdbCertificate_basic(commonName))
}

// lintignore:AT004
func testAccCcmPrivateCertificateExport_IIS(commonName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_certificate_export" "iiscert" {
    type           = "IIS"
    certificate_id = huaweicloud_ccm_private_certificate.test.id
}`, tesCmdbCertificate_basic(commonName))
}

// lintignore:AT004
func testAccCcmPrivateCertificateSM2Export_basic(commonName string) string {
	return fmt.Sprintf(`
provider "huaweicloud" {
  endpoints = {
    ccm = "https://ccm.cn-north-4.myhuaweicloud.com/"
  }
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type   = "ROOT"
  distinguished_name {
    common_name         = "%s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "SM2"
  signature_algorithm = "SM3"
  pending_days        = "7"
  validity {
    type  = "DAY"
    value = 2
  }
}

resource "huaweicloud_ccm_private_certificate" "test2" {
  distinguished_name {
    common_name = "%s"
  }
  issuer_id           = huaweicloud_ccm_private_ca.test_root.id
  key_algorithm       = "SM2"
  signature_algorithm = "SM3"
  validity {
    type  = "DAY"
    value = "1"
  }
}
	  
data "huaweicloud_ccm_private_certificate_export" "sm2cert" {
  type           = "other"
  certificate_id = huaweicloud_ccm_private_certificate.test2.id
  sm_standard    = "true"
}`, commonName, commonName)
}

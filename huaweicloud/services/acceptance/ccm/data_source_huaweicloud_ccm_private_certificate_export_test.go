package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPrivateCertificateExport_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceNameWithDash()
		apacheName = "data.huaweicloud_ccm_private_certificate_export.apache"
		apacheDc   = acceptance.InitDataSourceCheck(apacheName)

		nginxName = "data.huaweicloud_ccm_private_certificate_export.nginx"
		nginxDc   = acceptance.InitDataSourceCheck(nginxName)

		otherName = "data.huaweicloud_ccm_private_certificate_export.other"
		otherDc   = acceptance.InitDataSourceCheck(otherName)

		iisName = "data.huaweicloud_ccm_private_certificate_export.iis"
		iisDc   = acceptance.InitDataSourceCheck(iisName)

		tomcatName = "data.huaweicloud_ccm_private_certificate_export.tomcat"
		tomcatDc   = acceptance.InitDataSourceCheck(tomcatName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateCertificateExport_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					apacheDc.CheckResourceExists(),
					resource.TestCheckResourceAttr(apacheName, "type", "APACHE"),
					resource.TestCheckResourceAttrSet(apacheName, "certificate"),
					resource.TestCheckResourceAttrSet(apacheName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(apacheName, "password"),
					resource.TestCheckResourceAttrSet(apacheName, "private_key"),

					nginxDc.CheckResourceExists(),
					resource.TestCheckResourceAttr(nginxName, "type", "NGINX"),
					resource.TestCheckResourceAttrSet(nginxName, "certificate"),
					resource.TestCheckResourceAttrSet(nginxName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(nginxName, "private_key"),

					otherDc.CheckResourceExists(),
					resource.TestCheckResourceAttr(otherName, "type", "OTHER"),
					resource.TestCheckResourceAttrSet(otherName, "certificate"),
					resource.TestCheckResourceAttrSet(otherName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(otherName, "password"),
					resource.TestCheckResourceAttrSet(otherName, "private_key"),

					iisDc.CheckResourceExists(),
					resource.TestCheckResourceAttr(iisName, "type", "IIS"),
					resource.TestCheckResourceAttrSet(iisName, "keystore_pass"),
					resource.TestCheckResourceAttrSet(iisName, "server_pfx"),

					tomcatDc.CheckResourceExists(),
					resource.TestCheckResourceAttr(tomcatName, "type", "TOMCAT"),
					resource.TestCheckResourceAttrSet(tomcatName, "password"),
					resource.TestCheckResourceAttrSet(tomcatName, "server_jks"),
				),
			},
		},
	})
}

func testAccPrivateCertificateExport_basic(commonName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_certificate_export" "apache" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  type           = "APACHE"
  password       = "encryption_password"
}

data "huaweicloud_ccm_private_certificate_export" "nginx" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  type           = "NGINX"
}

data "huaweicloud_ccm_private_certificate_export" "other" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  type           = "OTHER"
  password       = "encryption_password"
}

data "huaweicloud_ccm_private_certificate_export" "iis" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  type           = "IIS"
}

data "huaweicloud_ccm_private_certificate_export" "tomcat" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  type           = "TOMCAT"
  password       = "encryption_password"
}
`, testPrivateCertificate_basic(commonName))
}

func TestAccPrivateCertificateExport_sm2cert(t *testing.T) {
	var (
		rName     = acceptance.RandomAccResourceNameWithDash()
		otherName = "data.huaweicloud_ccm_private_certificate_export.other"
		otherDc   = acceptance.InitDataSourceCheck(otherName)

		sm2certName = "data.huaweicloud_ccm_private_certificate_export.sm2cert"
		sm2certDc   = acceptance.InitDataSourceCheck(sm2certName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateCertificateExport_sm2cert(rName),
				Check: resource.ComposeTestCheckFunc(
					otherDc.CheckResourceExists(),
					resource.TestCheckResourceAttr(otherName, "type", "OTHER"),
					resource.TestCheckResourceAttrSet(otherName, "certificate"),
					resource.TestCheckResourceAttrSet(otherName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(otherName, "certificate_id"),
					resource.TestCheckResourceAttrSet(otherName, "enc_certificate"),
					resource.TestCheckResourceAttrSet(otherName, "enc_private_key"),
					resource.TestCheckResourceAttrSet(otherName, "password"),
					resource.TestCheckResourceAttrSet(otherName, "private_key"),

					sm2certDc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sm2certName, "type", "OTHER"),
					resource.TestCheckResourceAttrSet(sm2certName, "certificate"),
					resource.TestCheckResourceAttrSet(sm2certName, "certificate_chain"),
					resource.TestCheckResourceAttrSet(sm2certName, "enc_certificate"),
					resource.TestCheckResourceAttrSet(sm2certName, "enc_sm2_enveloped_key"),
					resource.TestCheckResourceAttrSet(sm2certName, "private_key"),
					resource.TestCheckResourceAttrSet(sm2certName, "signed_and_enveloped_data"),
					resource.TestCheckResourceAttrSet(sm2certName, "sm_standard"),
				),
			},
		},
	})
}

func testAccPrivateCertificateExport_sm2cert(commonName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_certificate_export" "other" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  type           = "OTHER"
  password       = "encryption_password"
}

data "huaweicloud_ccm_private_certificate_export" "sm2cert" {
  certificate_id = huaweicloud_ccm_private_certificate.test.id
  type           = "OTHER"
  sm_standard    = "true"
}
`, testAccPrivateCertificateExport_sm2Base(commonName))
}

func testAccPrivateCertificateExport_sm2Base(commonName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca" "test" {
  type   = "ROOT"
  distinguished_name {
    common_name         = "%[1]s-root"
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

resource "huaweicloud_ccm_private_certificate" "test" {
  distinguished_name {
    common_name = "%[1]s"
  }
  issuer_id           = huaweicloud_ccm_private_ca.test.id
  key_algorithm       = "SM2"
  signature_algorithm = "SM3"
  validity {
    type  = "DAY"
    value = "1"
  }
}
`, commonName)
}

package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCcmPrivateCertificateExport_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	expEesourceName := "data.huaweicloud_ccm_private_certificate_export.test3"
	resourceName := "huaweicloud_ccm_private_certificate.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcmPrivateCertificateExport_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(expEesourceName, "private_key"),
					resource.TestCheckResourceAttrSet(expEesourceName, "certificate"),
					resource.TestCheckResourceAttrSet(expEesourceName, "certificate_chain"),
				),
			},
		},
	})
}

func TestAccCcmPrivateCertificateExport_SM2(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	expEesourceSM2Name := "data.huaweicloud_ccm_private_certificate_export.test4"
	resourceName := "huaweicloud_ccm_private_certificate.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCcmPrivateCertificateSM2Export_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(expEesourceSM2Name, "private_key"),
					resource.TestCheckResourceAttrSet(expEesourceSM2Name, "certificate"),
					resource.TestCheckResourceAttrSet(expEesourceSM2Name, "certificate_chain"),
					resource.TestCheckResourceAttrSet(expEesourceSM2Name, "enc_certificate"),
					resource.TestCheckResourceAttrSet(expEesourceSM2Name, "enc_sm2_enveloped_key"),
					resource.TestCheckResourceAttrSet(expEesourceSM2Name, "signed_and_enveloped_data"),
				),
			},
		},
	})
}

// lintignore:AT004
func testAccCcmPrivateCertificateExport_basic(commonName string) string {
	return fmt.Sprintf(`
  %s

data "huaweicloud_ccm_private_certificate_export" "test3" {
    type           = "other"
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
	  
data "huaweicloud_ccm_private_certificate_export" "test4" {
  type           = "other"
  certificate_id = huaweicloud_ccm_private_certificate.test2.id
  sm_standard    = "true"
}`, commonName, commonName)
}

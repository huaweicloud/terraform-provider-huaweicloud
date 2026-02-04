package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ccm"
)

func getCCMCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("scm", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCM client: %s", err)
	}

	return ccm.ReadCCMCertificate(client, state.Primary.ID)
}

func TestAccCCMCertificate_basic(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ccm_certificate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCCMCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCCMCertificate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cert_brand", "GEOTRUST"),
					resource.TestCheckResourceAttr(rName, "cert_type", "OV_SSL_CERT"),
					resource.TestCheckResourceAttr(rName, "domain_type", "SINGLE_DOMAIN"),
					resource.TestCheckResourceAttr(rName, "effective_time", "1"),
					resource.TestCheckResourceAttr(rName, "domain_numbers", "1"),
					resource.TestCheckResourceAttr(rName, "tags.%", "2"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(rName, "validity_period"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "order_id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "push_support"),
				),
			},
			{
				Config: testCCMCertificate_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "tags.%", "1"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar_update"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"effective_time", "single_domain_number", "tags"},
			},
		},
	})
}

func testCCMCertificate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand            = "GEOTRUST"
  cert_type             = "OV_SSL_CERT"
  domain_type           = "SINGLE_DOMAIN"
  effective_time        = 1
  domain_numbers        = 1
  enterprise_project_id = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testCCMCertificate_basic_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand            = "GEOTRUST"
  cert_type             = "OV_SSL_CERT"
  domain_type           = "SINGLE_DOMAIN"
  effective_time        = 1
  domain_numbers        = 1
  enterprise_project_id = "%s"

  tags = {
    foo = "bar_update"
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func TestAccCCMCertificate_multiDomain(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ccm_certificate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCCMCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCCMCertificate_multiDomain,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cert_brand", "GEOTRUST"),
					resource.TestCheckResourceAttr(rName, "cert_type", "OV_SSL_CERT"),
					resource.TestCheckResourceAttr(rName, "domain_type", "MULTI_DOMAIN"),
					resource.TestCheckResourceAttr(rName, "effective_time", "1"),
					resource.TestCheckResourceAttr(rName, "domain_numbers", "4"),
					resource.TestCheckResourceAttr(rName, "primary_domain_type", "SINGLE_DOMAIN"),
					resource.TestCheckResourceAttr(rName, "single_domain_number", "1"),
					resource.TestCheckResourceAttr(rName, "wildcard_domain_number", "2"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "order_id"),
					resource.TestCheckResourceAttrSet(rName, "push_support"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "validity_period"),
					resource.TestCheckResourceAttrSet(rName, "wildcard_domain_number"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"effective_time", "single_domain_number"},
			},
		},
	})
}

const testCCMCertificate_multiDomain = `
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand             = "GEOTRUST"
  cert_type              = "OV_SSL_CERT"
  domain_type            = "MULTI_DOMAIN"
  effective_time         = 1
  domain_numbers         = 4
  primary_domain_type    = "SINGLE_DOMAIN"
  single_domain_number   = 1
  wildcard_domain_number = 2
}
`

func TestAccCCMCertificate_wildcardDomain(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ccm_certificate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCCMCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCCMCertificate_wildcardDomain,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "cert_brand", "GEOTRUST"),
					resource.TestCheckResourceAttr(rName, "cert_type", "OV_SSL_CERT"),
					resource.TestCheckResourceAttr(rName, "domain_type", "WILDCARD"),
					resource.TestCheckResourceAttr(rName, "effective_time", "1"),
					resource.TestCheckResourceAttr(rName, "domain_numbers", "1"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "order_id"),
					resource.TestCheckResourceAttrSet(rName, "push_support"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "validity_period"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"effective_time", "single_domain_number"},
			},
		},
	})
}

const testCCMCertificate_wildcardDomain = `
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand     = "GEOTRUST"
  cert_type      = "OV_SSL_CERT"
  domain_type    = "WILDCARD"
  effective_time = 1
  domain_numbers = 1
}
`

func TestAccCCMCertificate_epsID(t *testing.T) {
	var (
		obj   interface{}
		rName = "huaweicloud_ccm_certificate.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCCMCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCCMCertificate_epsID,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testCCMCertificate_epsID_update1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testCCMCertificate_epsID_update2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"effective_time", "single_domain_number"},
			},
		},
	})
}

var testCCMCertificate_epsID = fmt.Sprintf(`
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand            = "GEOTRUST"
  cert_type             = "OV_SSL_CERT"
  domain_type           = "SINGLE_DOMAIN"
  effective_time        = 1
  domain_numbers        = 1
  enterprise_project_id = "%s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)

const testCCMCertificate_epsID_update1 = `
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand            = "GEOTRUST"
  cert_type             = "OV_SSL_CERT"
  domain_type           = "SINGLE_DOMAIN"
  effective_time        = 1
  domain_numbers        = 1
  enterprise_project_id = "0"
}
`

var testCCMCertificate_epsID_update2 = fmt.Sprintf(`
resource "huaweicloud_ccm_certificate" "test" {
  cert_brand            = "GEOTRUST"
  cert_type             = "OV_SSL_CERT"
  domain_type           = "SINGLE_DOMAIN"
  effective_time        = 1
  domain_numbers        = 1
  enterprise_project_id = "%s"
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)

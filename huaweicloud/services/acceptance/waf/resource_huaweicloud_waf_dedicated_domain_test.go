package waf

import (
	"fmt"
	"testing"

	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWafDedicateDomainV1_basic(t *testing.T) {
	var domain domains.PremiumHost
	resourceName := "huaweicloud_waf_dedicated_domain.domain_1"
	randName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_name"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "protect_status"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "tls"),
					resource.TestCheckResourceAttrSet(resourceName, "cipher"),
					resource.TestCheckResourceAttrSet(resourceName, "alarm_page.template_name"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_3ds"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_dss"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.15"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_policy(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
			},
		},
	})
}

func TestAccWafDedicateDomainV1_withEpsID(t *testing.T) {
	var domain domains.PremiumHost
	resourceName := "huaweicloud_waf_dedicated_domain.domain_1"
	randName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafDedicatedDomainV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basic_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_name"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "protect_status"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "tls"),
					resource.TestCheckResourceAttrSet(resourceName, "cipher"),
					resource.TestCheckResourceAttrSet(resourceName, "alarm_page.template_name"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_3ds"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_dss"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_update_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.15"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_policy_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafDedicatedDomainV1Exists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
		},
	})
}

func testAccCheckWafDedicatedDomainV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	c, err := config.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud WAF dedicated client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_dedicated_domain" {
			continue
		}

		_, err := domains.GetWithEpsID(c, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"])
		if err == nil {
			return fmt.Errorf("WAF dedicated mode domain still exists")
		}
	}

	return nil
}

func testAccCheckWafDedicatedDomainV1Exists(n string, domain *domains.PremiumHost) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		c, err := config.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating HuaweiCloud WAF dedicated client: %s", err)
		}
		found, err := domains.GetWithEpsID(c, rs.Primary.ID, rs.Primary.Attributes["enterprise_project_id"])
		if err != nil {
			return err
		}
		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF dedicated domain not found")
		}
		*domain = *found
		return nil
	}
}

func testAccWafDedicatedDomainV1_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = huaweicloud_waf_certificate.certificate_1.id
  keep_policy    = false
  proxy          = false
  tls            = "TLS v1.1"
  cipher         = "cipher_1"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDedicatedDomainV1_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = huaweicloud_waf_certificate.certificate_1.id
  keep_policy    = false
  proxy          = true
  tls            = "TLS v1.2"
  cipher         = "cipher_2"
  pci_3ds        = true
  pci_dss        = true

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name)
}

func testAccWafDedicatedDomainV1_policy(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name = "%s"
}

resource "huaweicloud_waf_dedicated_domain" "domain_1" {
  domain         = "www.%s.com"
  certificate_id = huaweicloud_waf_certificate.certificate_1.id
  policy_id      = huaweicloud_waf_policy.policy_1.id
  keep_policy    = true
  proxy          = true
  tls            = "TLS v1.2"
  protect_status = 0

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf(name), name, name)
}

func testAccWafDedicatedDomainV1_basic_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_domain" "domain_1" {
  domain                = "www.%s.com"
  certificate_id        = huaweicloud_waf_certificate.certificate_1.id
  keep_policy           = false
  proxy                 = false
  tls                   = "TLS v1.1"
  cipher                = "cipher_1"
  enterprise_project_id = "%s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}

func testAccWafDedicatedDomainV1_update_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_dedicated_domain" "domain_1" {
  domain                = "www.%s.com"
  certificate_id        = huaweicloud_waf_certificate.certificate_1.id
  keep_policy           = false
  proxy                 = true
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = true
  pci_dss               = true
  enterprise_project_id = "%s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}

func testAccWafDedicatedDomainV1_policy_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "policy_1" {
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
}

resource "huaweicloud_waf_dedicated_domain" "domain_1" {
  domain                = "www.%[2]s.com"
  certificate_id        = huaweicloud_waf_certificate.certificate_1.id
  policy_id             = huaweicloud_waf_policy.policy_1.id
  keep_policy           = true
  proxy                 = true
  tls                   = "TLS v1.2"
  protect_status        = 0
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}

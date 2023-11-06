package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getWafDedicateDomainResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF dedicated client: %s", err)
	}

	epsID := state.Primary.Attributes["enterprise_project_id"]
	return domains.GetWithEpsID(client, state.Primary.ID, epsID)
}

func TestAccWafDedicateDomainV1_basic(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_dedicated_domain.domain_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWafDedicateDomainResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "website_name", "websiteName"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.http_return_code", "404"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.block_page_type", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key1", "$time_local"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$tenant_id"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_page.0.page_content"),
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
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "website_name", "websiteName_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "${http_host}/error.html"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.15"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$request_length"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key3", "$remote_addr"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_policy(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_dedicated_domain.domain_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWafDedicateDomainResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomainV1_basic_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "${http_host}/error.html"),
					resource.TestCheckResourceAttr(resourceName, "website_name", ""),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
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
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", ""),
					resource.TestCheckResourceAttr(resourceName, "website_name", "websiteName"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.15"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$request_length"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key3", "$remote_addr"),
				),
			},
			{
				Config: testAccWafDedicatedDomainV1_policy_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
				ImportStateIdFunc:       testWAFResourceImportState(resourceName),
			},
		},
	})
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
  protect_status = 1
  website_name   = "websiteName"
  description    = "test description"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }

  custom_page {
    http_return_code = "404"
    block_page_type  = "application/json"
    page_content     = <<EOF
{
  "event_id": "$${waf_event_id}",
  "error_msg": "error message"
}
EOF
  }

  forward_header_map = {
    "key1" = "$time_local"
    "key2" = "$tenant_id"
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
  protect_status = 0
  website_name   = "websiteName_update"
  description    = "test description update"
  redirect_url   = "$${http_host}/error.html"

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

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
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

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
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
  redirect_url          = "$${http_host}/error.html"
  enterprise_project_id = "%s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = huaweicloud_vpc.test.id
  }
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
  website_name          = "websiteName"
  description           = "test description"
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

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}

func testAccWafDedicatedDomainV1_policy_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "policy_1" {
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"

  depends_on = [
    huaweicloud_waf_certificate.certificate_1
  ]
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
}
`, testAccWafCertificateV1_conf_withEpsID(name, epsID), name, epsID)
}

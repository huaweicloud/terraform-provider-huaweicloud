package waf

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourceObj(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}
	return domains.GetWithEpsID(c, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"]).Extract()
}

// Before running the test case, please ensure that there is at least one WAF cloud instance in the current region.
func TestAccDomain_basic(t *testing.T) {
	var (
		domain domains.Domain

		resourceName    = "huaweicloud_waf_domain.test"
		randName        = acceptance.RandomAccResourceName()
		domainName      = fmt.Sprintf("%s.huawei.com", randName)
		certificateBody = testAccWafCertificate_basic(randName, generateCertificateBody())
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceObj,
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
				Config: testAccWafDomain_basic(certificateBody, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "http2_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.http_return_code", "400"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.block_page_type", "application/json"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_page.0.page_content"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "200"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "200"),
					resource.TestCheckResourceAttr(resourceName, "description", "web_description_1"),
					resource.TestCheckResourceAttr(resourceName, "lb_algorithm", "ip_hash"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key1", "$time_local"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "website_name", "websiteName"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "access_status"),
					resource.TestCheckResourceAttrSet(resourceName, "access_code"),
				),
			},
			{
				Config:      testAccWafDomain_withIpv6Enable(certificateBody, randName, domainName),
				ExpectError: regexp.MustCompile(`when type in server contains IPv6 address`),
			},
			{
				Config: testAccWafDomain_update1(certificateBody, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv6"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "2"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.0"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "$remote_addr"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag"),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "http2_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "${http_host}/error.html"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "description", "web_description_2"),
					resource.TestCheckResourceAttr(resourceName, "lb_algorithm", "round_robin"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$request_length"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key3", "$remote_addr"),
					resource.TestCheckResourceAttr(resourceName, "website_name", "websiteNameUpdate"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "1"),
				),
			},
			{
				Config: testAccWafDomain_update2(certificateBody, randName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "3"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "ip_tag_another"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag_update"),
					resource.TestCheckResourceAttrPair(resourceName, "policy_id", "huaweicloud_waf_policy.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", ""),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "180"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "3600"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "3600"),
					resource.TestCheckResourceAttr(resourceName, "lb_algorithm", "session_hash"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy", "charging_mode", "ipv6_enable"},
				ImportStateIdFunc:       testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccWafDomain_basic(certificateBody, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_domain" "test" {
  domain                = "%[2]s"
  certificate_id        = huaweicloud_waf_certificate.test.id
  certificate_name      = huaweicloud_waf_certificate.test.name
  proxy                 = false
  description           = "web_description_1"
  website_name          = "websiteName"
  lb_algorithm          = "ip_hash"
  protect_status        = 0
  enterprise_project_id = "%[3]s"

  custom_page {
    http_return_code = "400"
    block_page_type  = "application/json"
    page_content     = <<EOF
{
  "event_id": "$${waf_event_id}",
  "error_msg": "error message"
}
EOF
  }

  timeout_settings {
    connection_timeout = 50
    read_timeout       = 200
    write_timeout      = 200
  }

  forward_header_map = {
    "key1" = "$time_local"
    "key2" = "$tenant_id"
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
  }
}
`, certificateBody, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDomain_withIpv6Enable(certificateBody, randName, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "%[4]s"
}

resource "huaweicloud_waf_domain" "test" {
  domain                = "%[3]s"
  certificate_id        = huaweicloud_waf_certificate.test.id
  certificate_name      = huaweicloud_waf_certificate.test.name
  policy_id             = huaweicloud_waf_policy.test.id
  proxy                 = true
  ipv6_enable           = false
  lb_algorithm          = "session_hash"
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = "true"
  pci_dss               = "true"
  enterprise_project_id = "%[4]s"

  timeout_settings {
    connection_timeout = 180
    read_timeout       = 3600
    write_timeout      = 3600
  }

  traffic_mark {
    ip_tags     = ["ip_tag_update", "ip_tag_another"]
    session_tag = "session_tag_update"
    user_tag    = "user_tag_update"
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "3ffe:1900:fe21:4545::0"
    port            = 8443
    type            = "ipv6"
    weight          = 3
  }
}
`, certificateBody, randName, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDomain_update1(certificateBody, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "test" {
  domain                = "%s"
  certificate_id        = huaweicloud_waf_certificate.test.id
  certificate_name      = huaweicloud_waf_certificate.test.name
  proxy                 = true
  http2_enable          = true
  ipv6_enable           = true
  redirect_url          = "$${http_host}/error.html"
  description           = "web_description_2"
  lb_algorithm          = "round_robin"
  website_name          = "websiteNameUpdate"
  tls                   = "TLS v1.0"
  cipher                = "cipher_1"
  protect_status        = 1
  enterprise_project_id = "%[3]s"
  
  timeout_settings {
    connection_timeout = 100
    read_timeout       = 100
    write_timeout      = 100
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }

  traffic_mark {
    ip_tags     = ["ip_tag", "$remote_addr"]
    session_tag = "session_tag"
    user_tag    = "user_tag"
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "3ffe:1900:fe21:4545::0"
    port            = 8443
    type            = "ipv6"
    weight          = 2
  }
}
`, certificateBody, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDomain_update2(certificateBody, randName, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "%[4]s"
}

resource "huaweicloud_waf_domain" "test" {
  domain                = "%[3]s"
  certificate_id        = huaweicloud_waf_certificate.test.id
  certificate_name      = huaweicloud_waf_certificate.test.name
  policy_id             = huaweicloud_waf_policy.test.id
  proxy                 = true
  lb_algorithm          = "session_hash"
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = "true"
  pci_dss               = "true"
  protect_status        = 0
  enterprise_project_id = "%[4]s"

  timeout_settings {
    connection_timeout = 180
    read_timeout       = 3600
    write_timeout      = 3600
  }

  traffic_mark {
    ip_tags     = ["ip_tag_update", "ip_tag_another"]
    session_tag = "session_tag_update"
    user_tag    = "user_tag_update"
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    weight          = 3
  }
}
`, certificateBody, randName, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// This test case can only be run on the international site. China site does not support postPaid.
// Before running the test case, please ensure that there is at least one WAF cloud instance in the current region.
func TestAccDomain_postPaid(t *testing.T) {
	var (
		domain domains.Domain

		resourceName    = "huaweicloud_waf_domain.test"
		randName        = acceptance.RandomAccResourceName()
		domainName      = fmt.Sprintf("%s.huawei.com", randName)
		certificateBody = testAccWafCertificate_basic(randName, generateCertificateBody())
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceObj,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case can only be tested at international state, so a separate switch is configured.
			acceptance.TestAccPreCheckWafInternationalInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomain_postPaid(certificateBody, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
				),
			},
		},
	})
}

func testAccWafDomain_postPaid(certificateBody, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_domain" "test" {
  domain                = "%[2]s"
  certificate_id        = huaweicloud_waf_certificate.test.id
  certificate_name      = huaweicloud_waf_certificate.test.name
  keep_policy           = false
  proxy                 = false
  charging_mode         = "postPaid"
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
  }
}
`, certificateBody, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

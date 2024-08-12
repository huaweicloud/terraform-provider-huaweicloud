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
		return nil, fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}
	return domains.GetWithEpsID(c, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"]).Extract()
}

func TestAccWafDomainV1_basic(t *testing.T) {
	var domain domains.Domain
	resourceName := "huaweicloud_waf_domain.domain_1"
	randName := acceptance.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.huawei.com", randName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceObj,
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
				Config: testAccWafDomainV1_basic(randName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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
					resource.TestCheckResourceAttr(resourceName, "protect_status", "-1"),
					resource.TestCheckResourceAttrSet(resourceName, "access_status"),
					resource.TestCheckResourceAttrSet(resourceName, "access_code"),
				),
			},
			{
				Config:      testAccWafDomainV1_withIpv6Enable(randName, domainName),
				ExpectError: regexp.MustCompile(`when type in server contains IPv6 address`),
			},
			{
				Config: testAccWafDomainV1_update1(randName, domainName),
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
				Config: testAccWafDomainV1_update2(randName, domainName),
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
					resource.TestCheckResourceAttrPair(resourceName, "policy_id", "huaweicloud_waf_policy.policy_1", "id"),
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
			},
		},
	})
}

func TestAccWafDomainV1_withEpsID(t *testing.T) {
	var domain domains.Domain
	resourceName := "huaweicloud_waf_domain.domain_1"
	randName := acceptance.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.huawei.com", randName)

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
				Config: testAccWafDomainV1_basic_withEpsID(randName, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv6"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "1"),
					resource.TestCheckResourceAttr(resourceName, "website_name", ""),
				),
			},
			{
				Config: testAccWafDomainV1_update_withEpsID(randName, domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "3"),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
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

func TestAccWafDomainV1_withPolicy(t *testing.T) {
	var domain domains.Domain
	resourceName := "huaweicloud_waf_domain.domain_1"
	randName := acceptance.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.huawei.com", randName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceObj,
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
				Config: testAccWafDomainV1_policy(randName, domainName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.weight", "3"),
					resource.TestCheckResourceAttrPair(resourceName, "policy_id", "huaweicloud_waf_policy.policy_1", "id"),
				),
			},
		},
	})
}

// This test case can only be run on the international site. China site does not supprot postPaid
func TestAccWafDomainV1_postPaid(t *testing.T) {
	var domain domains.Domain
	resourceName := "huaweicloud_waf_domain.domain_1"
	randName := acceptance.RandomAccResourceName()
	domainName := fmt.Sprintf("%s.huawei.com", randName)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceObj,
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
				Config: testAccWafDomainV1_postPaid(randName, domainName),
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

func testAccWafDomainV1_base(randName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code    = "enterprise"
  enterprise_project_id = "0"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}

resource "huaweicloud_waf_certificate" "certificate_1" {
  name = "%s"

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUehx07qc7un7IB7/X9lHCLkt/jPowDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMTA1MzEwOTI1NTJaFw0yMjA1
MzEwOTI1NTJaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCvmuH5ViGtGOlevJ8vOoN3Ak4pp3SescdAfQa/r4cO
z/bmBqBcZJTX9HODhiQzdemyLLs9aOkQXYIc8OrcaIsjns92XITVDpFW0ThGyjhT
ZdELj9LsbIcVzNPPclTcebZBlzAyX0oLqpHK73OUYQY2E6l44U9G8Id763Bnws9N
Rn3cg0qufrlUgdim/pYZ8ubjvlDJ9eEIhcsu9zu8c8i2+8qLjEsonx5PrwzNlYP3
JqAmZ2dcbQeSPfv5U6ZceKEZfegK+Cxv4rFd5F4Rdxl+SAIY+6mr7qu1dAlcVMLS
QcLlJLRWQ5NmqL9xju7Fbj2VZt+L6nb512iKaedPo2GfAgMBAAGjUzBRMB0GA1Ud
DgQWBBR5yzB/GujpSlLrn0l2p+BslakGzjAfBgNVHSMEGDAWgBR5yzB/GujpSlLr
n0l2p+BslakGzjAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCj
TqvcIk0Au/yOxOfIGUZzVkTiORxbwATAfRN6n/+mrgWnIbHG4XFqqjmFr7gGvHeH
+BuyU06VXJgKYaPUqbYl7eBd4Spm5v3Wq7C7i96dOHmG8fVcjQnTWleyEmUsEarv
A6/lhTqXV1+AuNUaH+9EbBUBsrCHGLkECBMKl0+cJN8lo5XncAtp7z1+O/Mn0Zi6
XyNOyvqcmmn8HUkSIS4RlJ2ohuZN6oFC3sYX9g9Vo++IkjGl3dRbf/7JutqBGHNE
RVKoPyaivymDDIIL/qSy/Pi2s0hzUhwc1M8td0K/AMxyeigwNG7mTH0RzX32bUkf
ZoURg5WiRskhtHEvBsLF
-----END CERTIFICATE-----
EOT

  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCvmuH5ViGtGOle
vJ8vOoN3Ak4pp3SescdAfQa/r4cOz/bmBqBcZJTX9HODhiQzdemyLLs9aOkQXYIc
8OrcaIsjns92XITVDpFW0ThGyjhTZdELj9LsbIcVzNPPclTcebZBlzAyX0oLqpHK
73OUYQY2E6l44U9G8Id763Bnws9NRn3cg0qufrlUgdim/pYZ8ubjvlDJ9eEIhcsu
9zu8c8i2+8qLjEsonx5PrwzNlYP3JqAmZ2dcbQeSPfv5U6ZceKEZfegK+Cxv4rFd
5F4Rdxl+SAIY+6mr7qu1dAlcVMLSQcLlJLRWQ5NmqL9xju7Fbj2VZt+L6nb512iK
aedPo2GfAgMBAAECggEAeMAvDS3uAEI2Dx/y8h3xUn9yUfBFH+6tTanrXxoK6+OT
Kj96O64qL4l3eQRflkdJiGx74FFomglCtDXxudflfXvxurkJ2hunUySQ5xScwLQt
mB6w8kP6a8IqD+bVdbn32ohk6u5dU0JZ+ErJlklVZRAGJAoCYox5DXwrEh6CP+bJ
pItgjv71tEEnX5sScQwV7FMRbjsPzXoJp8vCQjlUdetM1fk9rs3R2WSeFbPgLLtC
xY0+8Hexy0q6BLmyPZvFCaVIAzAHCYeCyzPK3xcm4odbrBmRL/amOg24CCny065N
MU9RFhEjQsY1RaK7dgkvjsntUZvU+aDcL8o6djOTuQKBgQDlDN/j2ntpGCtbTWH0
cVTW13Ze7U7iE3BfDO3m4VYP3Xi/v5FI8nHlmLrcl30H1dPKvMTec0dCBOqD1wzF
KiqHy8ELowO2CbXMYJpjuPzXH40/AE3eOJVTJM8mOeuFdeFgYCd/9cB7o5jfTA5Y
4zj8EmcRzsH1rNSnvo7/O9q6+wKBgQDERDSvP8RScEbzDKuN6uhzj1K2CAEnY6//
rDA1so18UhAie9NcAvlKa46jQTOcYD77g5h0WSlNt9ZbK9Plq9CY9psI0KNqN3Fl
YVKOKdD5m6Rifmg+lt8KLc/WocQ10DXpPTXzzuRlN/TaMDdN2pedEre/0AAMs8Ia
MIUnu4oyrQKBgQC6b6BNdqi9Ak9IIdR5g0XrGbXfzolGu0vcEkoSg5fpkfuXF/bJ
yY2rtIVkyGmc1w9tFfmol2yI8Ddy2LgsRAYaQl7/edCre3vev0LrqMck0ynE/hpj
purkojF6i+qI10p7h8ie/wmNmbv1BZMoBst7Yf9DH2gA8IynfRQn7DA9wQKBgGaU
M2kJDgX8UsjDbYKuLTIAzb0AMAIzUxBxIX1fRh2dEnvDdjOYBk1EK/fdoyjvENwJ
6ouc8j6BgBKEtKpMg6j+8wbHbTGdqrHPDQPqjSN4mpEz+i4EUqySRxep0tBBc3vl
FybHko3okhvbqXwSbL2Ww90HzI7XAPMJOv8KQO+9AoGBAJxxftNWvypBXGkPCdH2
f3ikvT2Vef9QZjqkvtipCecAkjM6ReLshVsdqFSv/ZmsVUeNKoTHvX2GnhweJM44
x7N2mFK4skBzVtMVbjAHVjG78UitVu+FrzqGreaJXHaduhgUH2iFWfw09joOotAM
X7ioLbTeWGBqFM+C80PkdBNp
-----END PRIVATE KEY-----
EOT

  depends_on = [
    huaweicloud_waf_cloud_instance.test
  ]
}
`, randName)
}

func testAccWafDomainV1_basic(randName, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  proxy            = false
  description      = "web_description_1"
  website_name     = "websiteName"
  lb_algorithm     = "ip_hash"
  protect_status   = -1

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
`, testAccWafDomainV1_base(randName), domainName)
}

func testAccWafDomainV1_update1(randName, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  proxy            = true
  http2_enable     = true
  ipv6_enable      = true
  redirect_url     = "$${http_host}/error.html"
  description      = "web_description_2"
  lb_algorithm     = "round_robin"
  website_name     = "websiteNameUpdate"
  tls              = "TLS v1.0"
  cipher           = "cipher_1"
  protect_status   = 1
  
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
`, testAccWafDomainV1_base(randName), domainName)
}

func testAccWafDomainV1_update2(randName, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "policy_1" {
  name = "%[2]s"
  
  depends_on = [
    huaweicloud_waf_cloud_instance.test
  ]
}

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%[2]s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  policy_id        = huaweicloud_waf_policy.policy_1.id
  proxy            = true
  lb_algorithm     = "session_hash"
  tls              = "TLS v1.2"
  cipher           = "cipher_2"
  pci_3ds          = "true"
  pci_dss          = "true"
  protect_status   = 0

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
`, testAccWafDomainV1_base(randName), domainName)
}

func testAccWafDomainV1_withIpv6Enable(randName, domainName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "policy_1" {
  name = "%[2]s"

  depends_on = [
    huaweicloud_waf_cloud_instance.test
  ]
}

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%[2]s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  policy_id        = huaweicloud_waf_policy.policy_1.id
  proxy            = true
  ipv6_enable      = false
  lb_algorithm     = "session_hash"
  tls              = "TLS v1.2"
  cipher           = "cipher_2"
  pci_3ds          = "true"
  pci_dss          = "true"

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
`, testAccWafDomainV1_base(randName), domainName)
}

func testAccWafDomainV1_policy(randName, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_policy" "policy_1" {
  name = "%s"

  depends_on = [
    huaweicloud_waf_cloud_instance.test
  ]
}

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  policy_id        = huaweicloud_waf_policy.policy_1.id
  keep_policy      = true
  proxy            = true

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    weight          = 3
  }
}
`, testAccWafDomainV1_base(randName), randName, domainName)
}

func testAccWafDomainV1_postPaid(randName, domainName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain           = "%s"
  certificate_id   = huaweicloud_waf_certificate.certificate_1.id
  certificate_name = huaweicloud_waf_certificate.certificate_1.name
  keep_policy      = false
  proxy            = false
  charging_mode    = "postPaid"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
  }
}
`, testAccWafDomainV1_base(randName), domainName)
}

func testAccWafDomainV1_base_withEpsID(randName, epsID string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cloud_instance" "test" {
  resource_spec_code    = "enterprise"
  enterprise_project_id = "%[2]s"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}

resource "huaweicloud_waf_certificate" "certificate_1" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIDazCCAlOgAwIBAgIUehx07qc7un7IB7/X9lHCLkt/jPowDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMTA1MzEwOTI1NTJaFw0yMjA1
MzEwOTI1NTJaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCvmuH5ViGtGOlevJ8vOoN3Ak4pp3SescdAfQa/r4cO
z/bmBqBcZJTX9HODhiQzdemyLLs9aOkQXYIc8OrcaIsjns92XITVDpFW0ThGyjhT
ZdELj9LsbIcVzNPPclTcebZBlzAyX0oLqpHK73OUYQY2E6l44U9G8Id763Bnws9N
Rn3cg0qufrlUgdim/pYZ8ubjvlDJ9eEIhcsu9zu8c8i2+8qLjEsonx5PrwzNlYP3
JqAmZ2dcbQeSPfv5U6ZceKEZfegK+Cxv4rFd5F4Rdxl+SAIY+6mr7qu1dAlcVMLS
QcLlJLRWQ5NmqL9xju7Fbj2VZt+L6nb512iKaedPo2GfAgMBAAGjUzBRMB0GA1Ud
DgQWBBR5yzB/GujpSlLrn0l2p+BslakGzjAfBgNVHSMEGDAWgBR5yzB/GujpSlLr
n0l2p+BslakGzjAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCj
TqvcIk0Au/yOxOfIGUZzVkTiORxbwATAfRN6n/+mrgWnIbHG4XFqqjmFr7gGvHeH
+BuyU06VXJgKYaPUqbYl7eBd4Spm5v3Wq7C7i96dOHmG8fVcjQnTWleyEmUsEarv
A6/lhTqXV1+AuNUaH+9EbBUBsrCHGLkECBMKl0+cJN8lo5XncAtp7z1+O/Mn0Zi6
XyNOyvqcmmn8HUkSIS4RlJ2ohuZN6oFC3sYX9g9Vo++IkjGl3dRbf/7JutqBGHNE
RVKoPyaivymDDIIL/qSy/Pi2s0hzUhwc1M8td0K/AMxyeigwNG7mTH0RzX32bUkf
ZoURg5WiRskhtHEvBsLF
-----END CERTIFICATE-----
EOT

  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCvmuH5ViGtGOle
vJ8vOoN3Ak4pp3SescdAfQa/r4cOz/bmBqBcZJTX9HODhiQzdemyLLs9aOkQXYIc
8OrcaIsjns92XITVDpFW0ThGyjhTZdELj9LsbIcVzNPPclTcebZBlzAyX0oLqpHK
73OUYQY2E6l44U9G8Id763Bnws9NRn3cg0qufrlUgdim/pYZ8ubjvlDJ9eEIhcsu
9zu8c8i2+8qLjEsonx5PrwzNlYP3JqAmZ2dcbQeSPfv5U6ZceKEZfegK+Cxv4rFd
5F4Rdxl+SAIY+6mr7qu1dAlcVMLSQcLlJLRWQ5NmqL9xju7Fbj2VZt+L6nb512iK
aedPo2GfAgMBAAECggEAeMAvDS3uAEI2Dx/y8h3xUn9yUfBFH+6tTanrXxoK6+OT
Kj96O64qL4l3eQRflkdJiGx74FFomglCtDXxudflfXvxurkJ2hunUySQ5xScwLQt
mB6w8kP6a8IqD+bVdbn32ohk6u5dU0JZ+ErJlklVZRAGJAoCYox5DXwrEh6CP+bJ
pItgjv71tEEnX5sScQwV7FMRbjsPzXoJp8vCQjlUdetM1fk9rs3R2WSeFbPgLLtC
xY0+8Hexy0q6BLmyPZvFCaVIAzAHCYeCyzPK3xcm4odbrBmRL/amOg24CCny065N
MU9RFhEjQsY1RaK7dgkvjsntUZvU+aDcL8o6djOTuQKBgQDlDN/j2ntpGCtbTWH0
cVTW13Ze7U7iE3BfDO3m4VYP3Xi/v5FI8nHlmLrcl30H1dPKvMTec0dCBOqD1wzF
KiqHy8ELowO2CbXMYJpjuPzXH40/AE3eOJVTJM8mOeuFdeFgYCd/9cB7o5jfTA5Y
4zj8EmcRzsH1rNSnvo7/O9q6+wKBgQDERDSvP8RScEbzDKuN6uhzj1K2CAEnY6//
rDA1so18UhAie9NcAvlKa46jQTOcYD77g5h0WSlNt9ZbK9Plq9CY9psI0KNqN3Fl
YVKOKdD5m6Rifmg+lt8KLc/WocQ10DXpPTXzzuRlN/TaMDdN2pedEre/0AAMs8Ia
MIUnu4oyrQKBgQC6b6BNdqi9Ak9IIdR5g0XrGbXfzolGu0vcEkoSg5fpkfuXF/bJ
yY2rtIVkyGmc1w9tFfmol2yI8Ddy2LgsRAYaQl7/edCre3vev0LrqMck0ynE/hpj
purkojF6i+qI10p7h8ie/wmNmbv1BZMoBst7Yf9DH2gA8IynfRQn7DA9wQKBgGaU
M2kJDgX8UsjDbYKuLTIAzb0AMAIzUxBxIX1fRh2dEnvDdjOYBk1EK/fdoyjvENwJ
6ouc8j6BgBKEtKpMg6j+8wbHbTGdqrHPDQPqjSN4mpEz+i4EUqySRxep0tBBc3vl
FybHko3okhvbqXwSbL2Ww90HzI7XAPMJOv8KQO+9AoGBAJxxftNWvypBXGkPCdH2
f3ikvT2Vef9QZjqkvtipCecAkjM6ReLshVsdqFSv/ZmsVUeNKoTHvX2GnhweJM44
x7N2mFK4skBzVtMVbjAHVjG78UitVu+FrzqGreaJXHaduhgUH2iFWfw09joOotAM
X7ioLbTeWGBqFM+C80PkdBNp
-----END PRIVATE KEY-----
EOT

  depends_on = [
    huaweicloud_waf_cloud_instance.test
  ]

}
`, randName, epsID)
}

func testAccWafDomainV1_basic_withEpsID(randName, domainName, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain                = "%s"
  certificate_id        = huaweicloud_waf_certificate.certificate_1.id
  certificate_name      = huaweicloud_waf_certificate.certificate_1.name
  proxy                 = false
  keep_policy           = false
  enterprise_project_id = "%s"
  ipv6_enable           = true

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "3ffe:1900:fe21:4545::0"
    port            = 8080
    type            = "ipv6"
    weight          = 1
  }
}
`, testAccWafDomainV1_base_withEpsID(randName, epsID), domainName, epsID)
}

func testAccWafDomainV1_update_withEpsID(randName, domainName, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_domain" "domain_1" {
  domain                = "%s"
  certificate_id        = huaweicloud_waf_certificate.certificate_1.id
  certificate_name      = huaweicloud_waf_certificate.certificate_1.name
  proxy                 = true
  keep_policy           = false
  enterprise_project_id = "%s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    weight          = 3
  }
}
`, testAccWafDomainV1_base_withEpsID(randName, epsID), domainName, epsID)
}
